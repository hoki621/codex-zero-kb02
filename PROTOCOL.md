# Host-Firmware CDC protocol v1

Status: v1. This file is the sole contract between `host/` and `firmware/`.
Protocol changes must update this document and focused parser tests in both
child repositories.

## Framing and grammar

- The transport is USB CDC carrying one ASCII message per line.
- Each line, including its terminating LF (`0x0a`), is at most 128 bytes.
  Senders use LF. Receivers may also accept CRLF and discard the CR.
- A message contains printable ASCII (`0x20` through `0x7e`) only. Tokens are
  separated by one space, with no leading or trailing spaces.
- Command names and enum values are uppercase and case-sensitive.
- Decimal integers have no sign or leading zeroes, except the value `0`.
- `generation` is a decimal integer from 1 through 18446744073709551615.
- `sequence` is a decimal integer from 0 through 4294967295.
- A receiver processes a message only after receiving its complete LF-terminated
  line. Partial lines have no effect.

The fixed v1 messages are:

```text
# Host to firmware
HELLO HOST 1
STATE <generation> <selected> <states>
PING <sequence>
OFFLINE <generation>

# Firmware to host
HELLO ZERO-KB02 1
PONG <sequence>
ESC <generation> <DOWN|UP>
KEY <generation> <slot> <DOWN|UP>
ENC <generation> <CW|CCW|DOWN|UP>
JOY <generation> <UP|DOWN|LEFT|RIGHT>
```

`slot` and a numeric `selected` are decimal integers from `0` through `5`.
`selected` may instead be `-` when no slot is selected. `states` is exactly six
characters, one per slot from 0 through 5:

| Code | Meaning |
| --- | --- |
| `W` | working |
| `I` | idle |
| `B` | blocked |
| `D` | done |
| `U` | unknown |
| `E` | empty; no agent is assigned to the slot |

`STATE` is always a complete snapshot. Differential updates are not part of v1.
The host sends it when the snapshot changes and at least once every five
seconds while Herdr is available.

## Handshake and versioning

1. After opening a CDC candidate, the host sends `HELLO HOST 1` and does not
   send any other command until the handshake succeeds.
2. Firmware receiving that exact line replies `HELLO ZERO-KB02 1`, enters the
   connected-but-offline state, and waits for `STATE` or `OFFLINE`.
3. The host accepts the port only after receiving that exact reply. It then
   sends a complete `STATE`, or `OFFLINE` if Herdr is unavailable.
4. Firmware emits input events only after receiving `STATE` for the current
   connection.

The final integer in each `HELLO` is the protocol major. A receiver that does
not support the received major must not process any non-`HELLO` message on that
connection. Firmware remains offline and replies with its supported
`HELLO ZERO-KB02 1`; the host treats any reply other than that exact v1 line as
an incompatible or unrelated device and closes the port. Neither side falls
back to another major.

Receiving a valid host `HELLO` at any time starts a new protocol session:
firmware discards its prior generation, replies with its `HELLO`, enters the
connected-but-offline state, and waits for a new `STATE` or `OFFLINE`.

## State and generation

The host owns `generation`. It chooses a fresh nonzero value for every protocol
session and a new value whenever any slot's agent binding changes. Status-only
changes, selection changes, and periodic retransmission may keep the same
generation. Firmware treats the value as an opaque token and copies the most
recent online `STATE` generation into every input event.

The host accepts an input event only when all of these are true:

1. the handshake is complete;
2. Herdr is online and the last state sent was `STATE`;
3. the event generation exactly equals the generation in the host's current
   complete slot mapping; and
4. every other field is valid for that event.

Otherwise the host discards the event without invoking a Herdr operation.
In particular, an event is stale after a slot binding change, `OFFLINE`, USB
reconnect, firmware re-handshake, or host restart. After rejecting a well-formed
stale event, the host retransmits the current complete `STATE` when online or
the current `OFFLINE` when offline. It does not reuse a generation while events
from an earlier mapping or connection may still be buffered.

`OFFLINE <generation>` invalidates the previous online state. The host uses a
new generation, firmware displays offline, clears its selected slot and slot
states, and stops emitting input events until a later `STATE` arrives.

## Heartbeat and recovery

After handshake, the host sends `PING <sequence>` at least once every five
seconds. Firmware immediately answers `PONG` with the identical sequence.
Sequence wraps from `4294967295` to `0`.

- Firmware enters offline state and stops emitting events if it receives no
  valid host message for 12 seconds.
- Host closes and reopens the port if it receives no matching `PONG` for 12
  seconds. Unsolicited or non-matching `PONG` messages do not acknowledge a
  heartbeat.
- A valid host message resets the firmware timeout. A matching `PONG` resets
  the host timeout. Invalid lines reset neither timeout.
- If Herdr disconnects while CDC remains connected, the host sends `OFFLINE`
  with a new generation. After Herdr recovers it sends a complete `STATE` with
  another new generation.

USB connection or either process restarting uses this recovery sequence:

```text
Host -> Firmware: HELLO HOST 1
Firmware -> Host: HELLO ZERO-KB02 1
Host -> Firmware: STATE 41827 2 WIBDUE
Host -> Firmware: PING 0
Firmware -> Host: PONG 0
```

No earlier state or event is replayed across the handshake.

Herdr becoming unavailable without a CDC disconnect is explicit:

```text
Host -> Firmware: OFFLINE 41828
```

A major mismatch never reaches normal traffic. For example, firmware receiving
`HELLO HOST 2` remains offline and replies `HELLO ZERO-KB02 1`; a v2-only host
then closes the port.

## Input events

`ESC` reports the debounced edge of physical K1. Firmware emits one `DOWN` and
one `UP` per physical press; the host sends Escape to the currently focused,
mapped Codex pane on `DOWN` only.

```text
Firmware -> Host: ESC 41827 DOWN
Firmware -> Host: ESC 41827 UP
```

`KEY` reports the debounced edge of one of the six agent keys. Firmware emits
one `DOWN` and one `UP` per physical press; the host performs focus on `DOWN`
only.

```text
Firmware -> Host: KEY 41827 2 DOWN
Firmware -> Host: KEY 41827 2 UP
```

`ENC` reports one encoder detent or a debounced encoder-button edge. Each
detent produces one `CW` or `CCW`; each button press produces one `DOWN` and one
`UP`.

```text
Firmware -> Host: ENC 41827 CW
Firmware -> Host: ENC 41827 DOWN
Firmware -> Host: ENC 41827 UP
```

`JOY` reports a single cardinal direction after the joystick crosses its
calibrated threshold from neutral. Firmware chooses one axis for a diagonal and
does not repeat a direction until the stick returns to neutral.

```text
Firmware -> Host: JOY 41827 LEFT
```

The mapping from these fixed events to allowlisted Herdr operations belongs to
the host. Firmware cannot name a Herdr method or shell command.

## Invalid input

A malformed line is any line with invalid ASCII, spacing, token count, integer,
enum, or field range. An unknown command is also malformed. Each receiver:

- discards the complete line without changing state or performing an action;
- continues reading later lines on the same connection; and
- sends no error response.

If LF has not appeared within the first 128 bytes, the receiver discards bytes
through the next LF using bounded memory. The oversized line has no partial
effect. A partial line is discarded when CDC disconnects.

Examples that must be ignored include:

```text
STATE 41827 2 WIBEU
ESC 0 DOWN
ESC 41827 DOWN extra
KEY 41826 2 DOWN extra
JOY 41827 DIAGONAL
RUN herdr focus 2
```

`ESC 0 DOWN` has an invalid generation, and the next two examples have extra
tokens. A syntactically valid `ESC 41826 DOWN` or `KEY 41826 2 DOWN` would
instead be rejected as stale when the current generation is `41827`, followed
by retransmission of the current `STATE`.
