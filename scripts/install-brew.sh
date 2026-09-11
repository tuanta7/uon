#!/usr/bin/env bash
set -euo pipefail

if command -v brew &>/dev/null; then
  echo "brew is already installed: $(brew --version | head -1)"
  exit 0
fi

/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

test -d /opt/homebrew && eval "$(/opt/homebrew/bin/brew shellenv)"
test -d /usr/local/Homebrew && eval "$(/usr/local/bin/brew shellenv)"

if ! grep -q 'brew shellenv' ~/.zprofile 2>/dev/null; then
  echo "eval \"\$($(brew --prefix)/bin/brew shellenv)\"" >> ~/.zprofile
fi

echo "brew installed: $(brew --version | head -1)"