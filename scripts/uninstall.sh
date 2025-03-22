#!/usr/bin/env bash

set -e

echo "Removing quill..."
sudo rm -f /usr/local/bin/quill
rm -rf ~/.quill

echo "quill has been successfully uninstalled."
