# zero-kb02 Codex controller

[日本語](README.md)

Turn a zero-kb02 into a controller for Codex CLI running in Herdr on macOS. Display the status of up to six conversations on the OLED and LEDs, switch conversations with the keys, and adjust reasoning effort with the encoder. The joystick moves the mouse pointer.

[Demo video on X](https://x.com/hoki621/status/2093605017819423047)

## System architecture

![Architecture](assets/architecture-en.drawio.svg)

Firmware on the zero-kb02 handles input and display. The Host bridge on your Mac connects the device to Herdr and changes conversation reasoning settings through Codex App Server.

## Requirements

- zero-kb02 and a USB data cable
- macOS and [Herdr](https://herdr.dev/)
- [Homebrew](https://brew.sh/), [mise](https://mise.jdx.dev/), and an account with access to Codex CLI

`mise.toml` specifies the Node.js, Go and TinyGo versions. The launcher uses Homebrew Codex CLI.

## Setup

1. Install Herdr, Homebrew and mise, then get the source and required tools:

   ```sh
   brew install --cask codex
   git clone --recurse-submodules https://github.com/hoki621/codex-zero-kb02.git
   cd codex-zero-kb02
   mise install
   mise exec -- sh -c 'cd host && npm ci && npm run build'
   ```

2. Authenticate Codex CLI and register the Herdr status plugin:

   ```sh
   codex login
   herdr plugin link --enabled "$PWD/host"
   ```

3. Follow the [Firmware instructions](firmware/README.md) to build and flash your device, then connect it over USB.

## Run

Start Herdr, then run the following from the repository root. Use a separate terminal or pane for each step.

1. Keep Codex App Server running in a normal Terminal. This service manages conversation execution and reasoning settings.

   ```sh
   mise exec -- node host/dist/src/codex-micro.js server
   ```

2. Start a device-connected Codex CLI in each Herdr pane, one pane per conversation, up to six.

   ```sh
   mise exec -- node host/dist/src/codex-micro.js
   ```

   This launcher links each Codex conversation to its Herdr pane. Use it instead of plain `codex` for reasoning and approval controls. Append `resume` to resume a conversation. To work in another project, change to that directory and specify the absolute path to `host/dist/src/codex-micro.js`.

3. Start the Host bridge in another normal Terminal. Replace the USB port with your device's exact path and close other serial monitors. List candidate ports with `ls /dev/cu.usbmodem*`.

   ```sh
   HERDR_SOCKET_PATH="$HOME/.config/herdr/herdr.sock" \
   ZERO_KB02_PORT=/dev/cu.usbmodemzero_kb02_v21 \
   mise exec -- node host/dist/src/main.js
   ```

   `HERDR_SOCKET_PATH` identifies the Herdr connection; `ZERO_KB02_PORT` identifies the zero-kb02 USB connection.

To stop, finish the Codex CLI sessions and stop the Host bridge, then press Ctrl-C in the App Server terminal.

## Controls

Keys are numbered from the top left: K1–K4, then K5–K8, then K9–K12.

| Input | Action |
| --- | --- |
| K1 | Escape in the focused Codex pane |
| K2 / K3 / K5 / K6 / K7 / K8 | Focus conversation panes 1–6 |
| K4 | Open/close the status popup |
| K9 / K10 | Approve/deny command execution on supported versions |
| K11 | Unassigned |
| K12 | New conversation when the focused Codex is idle/done |
| Encoder | Clockwise increases reasoning effort; counterclockwise decreases it |
| Joystick | Mouse pointer movement; push inputs are unassigned |

OLED slots are arranged as 1/2, 3/4, 5/6. W=working, I=idle, B=blocked, D=done, U=unknown and E=empty.

Approval keys support Codex CLI **0.155.1 and 0.160.0 only**, and require a verified operation target and approval prompt. **They are disabled on 0.162.0.** K4 controls Herdr's shared popup and may close another plugin's popup.

## Updates and development

Stop the CLI sessions, Host bridge and App Server before updating. From the repository root:

```sh
git pull --ff-only
git submodule update --init --recursive
mise install
mise exec -- sh -c 'cd host && npm ci && npm run build'
```

Restart App Server and CLI sessions after a Homebrew Codex upgrade as well. Source checks do not require a device:

```sh
mise exec -- sh -c 'cd host && npm run typecheck && npm test && npm run dry-run -- WIBDUE'
mise exec -- sh -c 'cd firmware && go test ./... && go vet ./...'
```

- [Host details and troubleshooting](host/README.md) / [Firmware](firmware/README.md)
- [Verification scope](docs/verification.md) / [Wire contract](PROTOCOL.md)
- [Libraries and licenses](UPSTREAMS.md)

Upstream library notices are retained in each child repository. This project's own code does not yet have a license.
