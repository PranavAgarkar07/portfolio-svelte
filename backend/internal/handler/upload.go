package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/image/draw"

	"portfolio-backend/internal/model"
)

var allowedMIME = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/gif":  true,
	"image/webp": true,
}

func (h *Handlers) handleUpload(c *fiber.Ctx) (string, string, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return "", "", fiber.NewError(400, "no image file provided")
	}

	mimeType := file.Header.Get("Content-Type")
	if !allowedMIME[mimeType] {
		return "", "", fiber.NewError(400, "invalid MIME type, use PNG/JPG/GIF/WEBP")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExt := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true}
	if !allowedExt[ext] {
		return "", "", fiber.NewError(400, "unsupported format, use PNG/JPG/JPEG/GIF/WEBP")
	}

	if file.Size > 10<<20 {
		return "", "", fiber.NewError(400, "file too large, max 10MB")
	}

	src, err := file.Open()
	if err != nil {
		return "", "", fiber.NewError(500, "failed to read file")
	}
	defer src.Close()

	sniff := make([]byte, 512)
	if _, err := io.ReadFull(src, sniff); err != nil && err != io.ErrUnexpectedEOF {
		return "", "", fiber.NewError(500, "failed to validate file")
	}
	sniffMIME := http.DetectContentType(sniff)
	if !allowedMIME[sniffMIME] {
		return "", "", fiber.NewError(400, "file content does not match image type")
	}

	src.Seek(0, io.SeekStart)
	data, err := io.ReadAll(src)
	if err != nil {
		return "", "", fiber.NewError(500, "failed to read file data")
	}

	hash := sha256.Sum256(data)
	filename := hex.EncodeToString(hash[:8]) + ext

	var thumbData []byte
	var thumbFilename string
	img, _, decodeErr := image.Decode(bytes.NewReader(data))
	if decodeErr == nil {
		const maxThumbWidth = 640
		bounds := img.Bounds()
		w := bounds.Dx()
		hh := bounds.Dy()

		thumbW, thumbH := w, hh
		if w > maxThumbWidth {
			thumbW = maxThumbWidth
			thumbH = hh * maxThumbWidth / w
		}

		thumb := image.NewRGBA(image.Rect(0, 0, thumbW, thumbH))
		draw.BiLinear.Scale(thumb, thumb.Bounds(), img, bounds, draw.Over, nil)

		var thumbBuf bytes.Buffer
		if err := jpeg.Encode(&thumbBuf, thumb, &jpeg.Options{Quality: 70}); err == nil {
			thumbData = thumbBuf.Bytes()
			thumbFilename = strings.TrimSuffix(filename, ext) + "_thumb.jpg"
		}
	}

	if h.S3Client != nil {
		_, err = h.S3Client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket:      aws.String(model.S3Bucket),
			Key:         aws.String("static/uploads/" + filename),
			Body:        bytes.NewReader(data),
			ContentType: aws.String(mimeType),
		})
		if err != nil {
			slog.Error("Failed to upload to S3", "error", err)
			return "", "", fiber.NewError(500, "failed to save image")
		}
		if thumbData != nil {
			h.S3Client.PutObject(context.Background(), &s3.PutObjectInput{
				Bucket:      aws.String(model.S3Bucket),
				Key:         aws.String("static/thumbs/" + thumbFilename),
				Body:        bytes.NewReader(thumbData),
				ContentType: aws.String("image/jpeg"),
			})
		}
		url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/static/uploads/%s", model.S3Bucket, model.S3Region, filename)
		return url, filename, nil
	}

	savePath := filepath.Join("static", "uploads", filename)
	if err := os.WriteFile(savePath, data, 0644); err != nil {
		slog.Error("Failed to save uploaded image", "error", err)
		return "", "", fiber.NewError(500, "failed to save image")
	}
	if thumbData != nil {
		thumbSavePath := filepath.Join("static", "thumbs", thumbFilename)
		os.MkdirAll(filepath.Dir(thumbSavePath), 0755)
		if err := os.WriteFile(thumbSavePath, thumbData, 0644); err != nil {
			slog.Error("Failed to save thumbnail", "error", err)
		}
	}

	scheme := "http"
	if c.Protocol() == "https" || strings.HasPrefix(c.Hostname(), "sentinel") {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s/static/uploads/%s", scheme, c.Hostname(), filename)
	return url, filename, nil
}

func (h *Handlers) HandleAdminUploadImage(c *fiber.Ctx) error {
	if !h.checkAdminKey(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	url, _, err := h.handleUpload(c)
	if err != nil {
		var fe *fiber.Error
		if errors.As(err, &fe) {
			return c.Status(fe.Code).JSON(fiber.Map{"error": fe.Message})
		}
		return ErrInternal(c)
	}
	return c.JSON(fiber.Map{"url": url})
}

func (h *Handlers) HandleBlogUploadImage(c *fiber.Ctx) error {
	url, _, err := h.handleUpload(c)
	if err != nil {
		var fe *fiber.Error
		if errors.As(err, &fe) {
			return c.Status(fe.Code).JSON(fiber.Map{"error": fe.Message})
		}
		return ErrInternal(c)
	}
	return c.JSON(fiber.Map{"url": url})
}
