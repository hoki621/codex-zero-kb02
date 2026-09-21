#!/bin/sh
# Build only in a fresh temporary directory. Never flash or open a device.
set -eu
check_versions() {
case "$(tinygo version)" in
  'tinygo version 0.40.1 '*'go1.25.13 '*) ;;
  *) echo 'Run mise install and use TinyGo 0.40.1 / Go 1.25.13 from the parent repository.' >&2; exit 1 ;;
esac
}
check_versions
tinygo version
go version
source_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
probe_dir=$(mktemp -d /tmp/kb02-library-probe.XXXXXX)
printf 'Build directory: %s\n' "$probe_dir"
cp "$source_dir/main.go" "$source_dir/go.mod" "$source_dir/go.sum" "$probe_dir/"
git clone --quiet https://github.com/sago35/tinygo-keyboard.git "$probe_dir/tinygo-keyboard"
git -C "$probe_dir/tinygo-keyboard" checkout --quiet cf173e98f60329b7f7feba941461bb95c065c418
git -C "$probe_dir/tinygo-keyboard" apply --check "$source_dir/../patches/tinygo-keyboard-input-only.patch"
git -C "$probe_dir/tinygo-keyboard" apply "$source_dir/../patches/tinygo-keyboard-input-only.patch"
cd "$probe_dir"
check_versions
tinygo list -target waveshare-rp2040-zero -tags kb02_inputonly -deps . > deps.txt
if grep -Fxq machine/usb/hid/keyboard deps.txt; then
  echo 'Unexpected HID keyboard dependency' >&2
  exit 1
fi
tinygo build -target waveshare-rp2040-zero -tags kb02_inputonly -stack-size 8kb -size short -o probe.uf2 .
shasum -a 256 probe.uf2
printf 'Compile-only probe retained at %s; do not flash it.\n' "$probe_dir"
