#!/bin/bash

# Usage: ./ship.sh "Your commit message"
# If no message provided, use timestamp

if [ -z "$1" ]
then
    MSG="Update: $(date)"
else
    MSG="$1"
fi

echo "🚀 Shipping: $MSG"

# 1. Stage changes
git add .

# 2. Commit
git commit -m "$MSG"

# 3. Deploy
# We are now relying on GitHub Actions to build and deploy with the correct API URL.
# npm run deploy <--- Removed to prevent local build (localhost) from overriding.

# 4. Push source code to main (triggers GitHub Action)
git push origin main

echo "✅ Deployed and Pushed!"
