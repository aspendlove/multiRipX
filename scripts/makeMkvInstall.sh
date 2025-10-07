#!/bin/bash
set -euo pipefail

if [[ -z "$1" ]]; then
    echo "Error: Missing version number." >&2
    echo "Usage: $0 <version>" >&2
    exit 1
fi

VERSION="$1"

function cleanup {
    if [[ -d "$TEMP_DIR" ]]; then
        echo "Cleaning up temporary directory: $TEMP_DIR"
        rm -rf "$TEMP_DIR" 2>/dev/null || true
    fi
}
trap cleanup EXIT INT TERM

TEMP_DIR=$(mktemp -d)
if [[ ! -d "$TEMP_DIR" ]]; then
    echo "Error: Could not create temporary directory." >&2
    exit 1
fi

cd "$TEMP_DIR" || exit 1

echo "Installing dependencies"
sudo dnf group install -y c-development development-tools
sudo dnf install -y zlib-devel openssl-devel expat-devel ffmpeg-devel qt5-qtbase-devel

echo "Downloading MakeMKV $VERSION"
wget --no-check-certificate "https://www.makemkv.com/download/makemkv-oss-${VERSION}.tar.gz"
wget --no-check-certificate "https://www.makemkv.com/download/makemkv-bin-${VERSION}.tar.gz"

tar -xzf makemkv-oss-${VERSION}.tar.gz
tar -xzf makemkv-bin-${VERSION}.tar.gz

echo "Building MakeMKV-OSS"
pushd "makemkv-oss-${VERSION}"
./configure
make
sudo make install
popd

echo "Building MakeMKV-BIN"
pushd "makemkv-bin-${VERSION}"

# Make the eula auto accept
echo "exit 0" > ./src/ask_eula.sh

make
sudo make install
popd

echo "Finished"