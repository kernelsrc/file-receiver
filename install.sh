#!/bin/bash
set -euo pipefail

echo "Installing file receiver dependencies..."

if command -v apt-get >/dev/null 2>&1; then
    sudo apt-get update -qq || true
    sudo apt-get install -y -qq gcc make qrencode || true
fi

echo "Dependencies installed: gcc, make, qrencode (optional)."
echo "To compile:  make -C c_file_receiver"
echo "To run:      ./c_file_receiver/server (or with sudo for port <1024)"
echo "QR Code:     after install, run 'qrencode -t ANSIUTF8 \"http://YOUR_IP:8080\"'"
