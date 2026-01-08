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

# 3. Deploy (Build + GH Pages Push)
# Note: You have both a GitHub Action AND a local deploy script.
# This script runs the local one as requested.
npm run deploy

# 4. Push source code to main (triggers Action backup or redundancy)
git push origin main

echo "✅ Deployed and Pushed!"
