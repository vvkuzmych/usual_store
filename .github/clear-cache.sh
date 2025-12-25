#!/bin/bash

# Script to clear GitHub Actions cache for the UsualStore repository
# Usage: ./clear-cache.sh

set -e

echo "🧹 Clearing GitHub Actions cache..."
echo ""

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) is not installed."
    echo ""
    echo "Install it with:"
    echo "  macOS: brew install gh"
    echo "  Linux: https://github.com/cli/cli/blob/trunk/docs/install_linux.md"
    echo ""
    exit 1
fi

# Check if authenticated
if ! gh auth status &> /dev/null; then
    echo "❌ You are not authenticated with GitHub CLI."
    echo ""
    echo "Run: gh auth login"
    echo ""
    exit 1
fi

echo "📋 Listing all caches..."
gh cache list

echo ""
echo "🗑️  Deleting all Go-related caches..."

# Delete all caches with "go" in the key
gh cache list | grep -i "go" | awk '{print $1}' | while read -r cache_id; do
    echo "  Deleting cache: $cache_id"
    gh cache delete "$cache_id" || true
done

echo ""
echo "✅ Cache clearing complete!"
echo ""
echo "💡 Next steps:"
echo "   1. Go to: https://github.com/YOUR_USERNAME/usual_store/actions"
echo "   2. Re-run the failed workflow"
echo "   3. The cache will be rebuilt from scratch"
echo ""

