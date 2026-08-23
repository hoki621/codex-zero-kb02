# codex-zero-kb02

An unofficial zero-kb02 controller for Codex agents running in Herdr.

This repository pins a compatible pair of independently versioned components:

- [`host/`](https://github.com/hoki621/codex-zero-kb02-host): Herdr and USB CDC bridge
- [`firmware/`](https://github.com/hoki621/codex-zero-kb02-firmware): TinyGo firmware for zero-kb02

This project is not affiliated with or endorsed by OpenAI, Work Louder, Herdr,
Waveshare, or the upstream projects listed in [UPSTREAMS.md](UPSTREAMS.md).

## Clone

```sh
git clone --recurse-submodules https://github.com/hoki621/codex-zero-kb02.git
cd codex-zero-kb02
mise install
```

## Development workflow

Work is coordinated in the [parent issue tracker](https://github.com/hoki621/codex-zero-kb02/issues).
The child repositories intentionally have Issues disabled.

1. Pick one issue and work only in the repository named by that issue.
2. Commit and push the child repository first.
3. Run the issue's focused verification.
4. Update the parent submodule pointer only during integration.

Do not use `git submodule update --remote`; the parent commit is the compatibility
record for the exact Host and Firmware commits.

## Safety

Do not flash firmware, open the device serial port, or send input to a Herdr pane
without explicit user approval for that operation. Build and mock tests are safe
defaults.
