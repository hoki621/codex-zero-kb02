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
cd host
npm ci
```

The parent commit pins the tested component pair:

| Component | Commit |
| --- | --- |
| Host | `ac0192aec142836fbc7fdcf6c31fa518acdead80` |
| Firmware | `9c39c5fcc5e77454b97c6ad54ba18cc6d43fac9a` |

## Run

Flash the pinned firmware only after verifying the target device and obtaining
explicit permission for that hardware operation. Build instructions are in
[`firmware/README.md`](firmware/README.md).

Start the bridge from a normal terminal outside Herdr-managed panes. Set
`HERDR_SOCKET_PATH` explicitly so the daemon can reconnect after a Herdr server
restart; use the exact serial path when more than one USB modem is connected:

```sh
cd host
npm run build
codex app-server daemon start
herdr plugin link --enabled "$(pwd)"
HERDR_SOCKET_PATH="$HOME/.config/herdr/herdr.sock" ZERO_KB02_PORT=/dev/cu.usbmodemzero_kb02_v11 npm start
```

The link command is an explicit setup step that mutates the local Herdr plugin
registry; the daemon never runs it automatically. It installs only the fixed
local `hoki621.zero-kb02` manifest and its read-only `status` popup.
Start managed Codex panes with `codex --remote unix://`; Encoder control reuses
that Codex CLI 0.149.1 App Server's local Unix endpoint.

Stop it with Ctrl-C. The daemon maps the six agent keys to safe `agent.focus`
requests only after resolving the current pane for the assigned terminal. A
daemon running inside a Herdr-managed pane cannot survive a server restart.

## Update

Update only to another parent commit so the tested child pair stays intact:

```sh
git pull --ff-only
git submodule update --init --recursive
mise install
cd host
npm ci
```

Do not use `git submodule update --remote`.

## Recovery

- USB disconnect: leave the daemon running; it retries the exact configured
  port and sends a fresh handshake and complete state after reconnect.
- Herdr restart: an externally running daemon reconnects to the configured
  Herdr socket, rebuilds the six slots, and retransmits a complete state. A
  five-second reconcile repairs missed Herdr events.
- Host restart: run the same explicit `HERDR_SOCKET_PATH` and `ZERO_KB02_PORT`
  command; generations are not reused across the new USB session. A stale,
  unreachable `status.sock` is removed only if its identity is unchanged;
  a live, replaced, or ambiguous socket remains untouched and startup fails.
- Firmware recovery: stop the daemon first and follow the verified procedure in
  [`firmware/docs/hardware-diagnostics.md`](firmware/docs/hardware-diagnostics.md).
  Flashing and BOOTSEL/RST always require separate explicit permission.

## Known constraints

- v1 supports up to six detected Codex agents and Herdr protocol 20 as shipped
  by Herdr 0.8.2.
- K1 sends scoped Escape to the focused mapped Codex pane. K2, K3, and K5-K8
  focus agent slots 0-5. K4 toggles the Herdr session's active popup globally:
  it closes any active popup, including one from another plugin, or opens the
  fixed status popup when none is open. Encoder CW/CCW changes the focused
  managed Codex CLI thread's reasoning effort by one supported level and clamps
  at the endpoints. K9-K12 have no Host action. The joystick moves the USB HID
  relative pointer; its push has no action.
- macOS USB discovery is limited to `/dev/cu.usbmodem*`; set
  `ZERO_KB02_PORT` when discovery is ambiguous.
- There is no launchd service, settings GUI, Vial control, arbitrary shell
  execution, Approve/Deny, push-to-talk, model switching, desktop App control,
  or Zed ACP support.

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
