#!/usr/bin/env bash
set -euo pipefail

REPO="Starrick2001/commit-craft"
ASSET="commit-craft"
INSTALL_DIR=""

usage() {
  cat <<'USAGE'
Usage: install.sh [options]

Options:
  -d, --dir DIR     Install directory (default: /usr/local/bin if writable, else ~/.local/bin)
  -a, --asset NAME  Release asset name to download (default: commit-craft)
  -h, --help        Show this help message

Examples:
  ./install.sh
  ./install.sh --dir ~/.local/bin
  ./install.sh --asset commit-craft
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -d|--dir)
      INSTALL_DIR="$2"
      shift 2
      ;;
    -a|--asset)
      ASSET="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$INSTALL_DIR" ]]; then
  if [[ -w "/usr/local/bin" ]]; then
    INSTALL_DIR="/usr/local/bin"
  else
    INSTALL_DIR="$HOME/.local/bin"
  fi
fi

if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
  echo "Error: curl or wget is required." >&2
  exit 1
fi

TMP_FILE="$(mktemp -t commit-craft.XXXXXX)"
cleanup() {
  rm -f "$TMP_FILE"
}
trap cleanup EXIT

URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"

if command -v curl >/dev/null 2>&1; then
  curl -fL -o "$TMP_FILE" "$URL"
else
  wget -O "$TMP_FILE" "$URL"
fi

chmod +x "$TMP_FILE"
mkdir -p "$INSTALL_DIR"

# Use sudo if needed for system-wide install directory.
if [[ ! -w "$INSTALL_DIR" ]]; then
  sudo mv "$TMP_FILE" "$INSTALL_DIR/commit-craft"
else
  mv "$TMP_FILE" "$INSTALL_DIR/commit-craft"
fi

printf "Installed commit-craft to %s/commit-craft\n" "$INSTALL_DIR"

if ! command -v commit-craft >/dev/null 2>&1; then
  echo "Note: ensure $INSTALL_DIR is on your PATH."
fi
