#!/usr/bin/env python3
"""Generate low-quality WebP thumbnails + tiny base64 placeholders for cert images."""

import base64
import hashlib
import io
import json
import os
import sys
from PIL import Image, ImageFilter

STATIC_DIR = os.path.join(os.path.dirname(__file__), "..", "backend", "static")
THUMB_DIR = os.path.join(STATIC_DIR, "thumbs")
OUTPUT_JS = os.path.join(
    os.path.dirname(__file__), "..", "src", "lib", "cert-placeholders.ts"
)

SOURCES = [
    "cert-nptl.png",
    "cert-redhat-python.png",
    "cert-redhat-linux.png",
    "cert-joomla.jpeg",
]

THUMB_WIDTH = 640
THUMB_QUALITY = 40
PLACEHOLDER_WIDTH = 24
PLACEHOLDER_BLUR = 8

os.makedirs(THUMB_DIR, exist_ok=True)

placeholders: dict[str, str] = {}

for src_name in SOURCES:
    src_path = os.path.join(STATIC_DIR, src_name)
    if not os.path.exists(src_path):
        print(f"SKIP {src_name} — not found", file=sys.stderr)
        continue

    img = Image.open(src_path)
    img = img.convert("RGB")

    # --- Thumbnail (WebP, 640px, 40% quality) ---
    thumb = img.copy()
    thumb.thumbnail((THUMB_WIDTH, int(THUMB_WIDTH * thumb.height / thumb.width)))
    thumb_name = os.path.splitext(src_name)[0] + ".webp"
    thumb_path = os.path.join(THUMB_DIR, thumb_name)
    thumb.save(thumb_path, "WEBP", quality=THUMB_QUALITY, method=6)
    thumb_size_kb = os.path.getsize(thumb_path) / 1024
    orig_size_kb = os.path.getsize(src_path) / 1024
    print(
        f"{src_name}: {orig_size_kb:.0f}KB → {thumb_name}: {thumb_size_kb:.0f}KB"
        f" ({thumb_size_kb/orig_size_kb*100:.0f}%)"
    )

    # --- Tiny blurred placeholder (base64 data URI) ---
    tiny = img.copy()
    tiny.thumbnail((PLACEHOLDER_WIDTH, int(PLACEHOLDER_WIDTH * tiny.height / tiny.width)))
    tiny = tiny.filter(ImageFilter.GaussianBlur(PLACEHOLDER_BLUR))
    buf = io.BytesIO()
    tiny.save(buf, "JPEG", quality=20)
    b64 = base64.b64encode(buf.getvalue()).decode()
    data_uri = f"data:image/jpeg;base64,{b64}"
    placeholders[src_name] = data_uri

# Write TypeScript module with placeholders keyed by filename
with open(OUTPUT_JS, "w") as f:
    f.write("// Auto-generated — do not edit. Run scripts/optimize-cert-images.py\n\n")
    f.write("export const certPlaceholders: Record<string, string> = ")
    f.write(json.dumps(placeholders, indent=2))
    f.write(";\n")

print(f"\nDone. {len(placeholders)} placeholders written to {OUTPUT_JS}")
