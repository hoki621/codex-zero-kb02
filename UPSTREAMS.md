# Upstream references and licenses

Firmware uses maintained input/display libraries and TinyGo standard USB APIs.
No workshop code is copied. Product firmware implementation remains user-owned.

| Project | Pinned reference | License / use |
| --- | --- | --- |
| [sago35/tinygo-keyboard](https://github.com/sago35/tinygo-keyboard/tree/cf173e98f60329b7f7feba941461bb95c065c418) | `cf173e98f60329b7f7feba941461bb95c065c418` | MIT, Copyright 2022 sago35; matrix/debounce unchanged; input-only build isolation patch |
| [TinyGo drivers](https://github.com/tinygo-org/drivers/tree/v0.34.0) | v0.34.0 | BSD-3-Clause; Encoder and SSD1306 |
| [tinyfont](https://github.com/tinygo-org/tinyfont/tree/v0.6.0) | v0.6.0 | BSD-3-Clause; bitmap fonts |
| [tinydraw](https://github.com/tinygo-org/tinydraw/tree/v0.4.0) | v0.4.0 | BSD-3-Clause; drawing |
| [pio](https://github.com/tinygo-org/pio/tree/v0.2.0) | v0.2.0 | BSD-3-Clause; RP2040 WS2812B |
| [TinyGo](https://github.com/tinygo-org/tinygo/tree/v0.40.1) | 0.40.1, Go 1.25.13 | BSD-3-Clause; ADC, standard CDC/HID mouse; bundled runtime retains its notices |
| House of Herdr packages/codex-micro | `50b24e3f334a38a84bfa356f154d49835dff2499` | MIT; historical Host reference; preserve child repository notices |
| sago35/keyboards zero-kb02/firmware | `4b18114b66637c5909229704c0a503bcdeccc057` | MIT; hardware mapping reference |
| [makiuchi-d/zero-micro](https://github.com/makiuchi-d/zero-micro) | Reference only; no code imported | Related project demonstrating library-based composition |
| tinygo-keeb/workshop | `6ca413538952d83aa2067183839c8977da564bb8` | No reusable license confirmed; reference only, no copying |

The only local library patch is [tinygo-keyboard-input-only.patch](docs/patches/tinygo-keyboard-input-only.patch),
with the [complete upstream MIT notice](docs/patches/tinygo-keyboard-LICENSE.txt).
It adds build constraints and minimal input-only types. `kbmatrix.go`, its scan
algorithm and debounce remain byte-for-byte upstream. Firmware must retain the
patched dependency's license and pin, and build with `-tags kb02_inputonly`.
No upstream PR or external publication was made.

[The compile-only probe](docs/library-probe/main.go) checks these public APIs
together; [check.sh](docs/library-probe/check.sh) clones the exact SHA, applies the
patch, rejects a HID keyboard dependency, and builds without opening a device.
It is a learning prerequisite, not flashable product firmware or USB acceptance.
The standard TinyGo CDCHID descriptor still includes unused keyboard report
items; the input-only dependency installs neither a keyboard handler nor a
Vial vendor interface. Real USB enumeration remains an H1 acceptance check.

This project does not import official product branding, icons or USB identifiers.

TomThumb has its own original font copyright/attribution in addition to the
tinyfont package license. Preserve [the complete TomThumb notice](docs/patches/tomthumb-LICENSE.txt)
with distributed firmware binaries/documentation, including the named original
authors and subsequent conversion/modification credits.
