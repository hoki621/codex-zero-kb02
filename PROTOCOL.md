# Host-Firmware protocol

Status: draft. The v1 contract is defined and reviewed in
[parent issue #3](https://github.com/hoki621/codex-zero-kb02/issues/3).

This file is the sole contract between `host/` and `firmware/`. Protocol changes
must update this document and focused parser tests in both child repositories.

The intended v1 is a bounded ASCII line protocol over USB CDC. It must use an
explicit major version and device handshake, send complete six-slot states, and
reject stale physical input by generation. It must not expose arbitrary Herdr
methods, shell commands, or a generic RPC mechanism to the device.
