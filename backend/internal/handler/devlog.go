package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"portfolio-backend/internal/model"
)

func (h *Handlers) generateDevLog() (string, error) {
	h.Breaker.Mu.Lock()
	if h.Breaker.Failures >= model.MaxFailures {
		if time.Since(h.Breaker.LastFailure) < model.BreakerTimeout {
			h.Breaker.Mu.Unlock()
			slog.Warn("Circuit breaker open — skipping OpenRouter call")
			return model.QuotaFallback, nil
		}
		h.Breaker.Failures = 0
	}
	h.Breaker.Mu.Unlock()

	apiKey := h.Config.OpenRouterAPIKey
	if apiKey == "" {
		return "System Error: Neural Link Disconnected (Missing API Key). Please configure the satellite uplink.", nil
	}

	events, err := h.fetchGitHubEvents()
	if err != nil {
		return "", err
	}

	h.ResponsesMu.Lock()
	diversityInstruction := ""
	if len(h.LastResponses) > 0 {
		diversityInstruction = fmt.Sprintf("Avoid repeating these recent responses: %s", strings.Join(h.LastResponses, "; "))
	}
	h.ResponsesMu.Unlock()

	prompt := fmt.Sprintf(`You are a senior software engineer named Pranav. A friend just asked "what have you been working on?" — answer naturally.

Today is %s. Here's what you've pushed to GitHub recently:

%s

Write ONE sentence (max 40 words) like you're telling a friend. Mention the project and what you actually did. Be casual, specific, and human.

%s`, time.Now().Weekday().String(), events, diversityInstruction)

	summary, err := h.callOpenRouter(apiKey, prompt)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "429") || strings.Contains(errStr, "RESOURCE_EXHAUSTED") || strings.Contains(errStr, "quota") {
			h.Breaker.Mu.Lock()
			h.Breaker.Failures++
			h.Breaker.LastFailure = time.Now()
			h.Breaker.Mu.Unlock()
			slog.Warn("OpenRouter quota exhausted", "failures", h.Breaker.Failures)
			return model.QuotaFallback, nil
		}
		return "", err
	}

	h.ResponsesMu.Lock()
	h.LastResponses = append(h.LastResponses, summary)
	if len(h.LastResponses) > 3 {
		h.LastResponses = h.LastResponses[len(h.LastResponses)-3:]
	}
	h.ResponsesMu.Unlock()

	return summary, nil
}

func (h *Handlers) fetchGitHubEvents() (string, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public", model.GitHubUsername)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if token := h.Config.GitHubToken; token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := h.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("github api error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var events []map[string]interface{}
	if err := json.Unmarshal(body, &events); err != nil {
		return "", err
	}

	relevantTypes := map[string]bool{
		"PushEvent":    true,
		"CreateEvent":  true,
		"ReleaseEvent": true,
		"IssuesEvent":  true,
	}

	var filtered []map[string]interface{}
	for _, event := range events {
		eventType, _ := event["type"].(string)
		if relevantTypes[eventType] {
			filtered = append(filtered, event)
		}
	}
	events = filtered

	var summaryBuilder strings.Builder
	count := 0
	for _, event := range events {
		if count >= 30 {
			break
		}
		eventType, _ := event["type"].(string)
		repo, _ := event["repo"].(map[string]interface{})
		repoName, _ := repo["name"].(string)

		summaryBuilder.WriteString(fmt.Sprintf("- %s on %s", eventType, repoName))

		if payload, ok := event["payload"].(map[string]interface{}); ok {
			if commits, ok := payload["commits"].([]interface{}); ok {
				for _, c := range commits {
					commit, _ := c.(map[string]interface{})
					msg, _ := commit["message"].(string)
					summaryBuilder.WriteString(fmt.Sprintf(": %s", msg))
				}
			}
		}
		summaryBuilder.WriteString("\n")
		count++
	}

	return summaryBuilder.String(), nil
}

func (h *Handlers) callOpenRouter(apiKey, text string) (string, error) {
	url := "https://openrouter.ai/api/v1/chat/completions"

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": model.OpenRouterModel,
		"messages": []interface{}{
			map[string]interface{}{
				"role":    "user",
				"content": text,
			},
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), model.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := h.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openrouter api error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	return "Analysis complete. Systems nominal (Default Response).", nil
}
