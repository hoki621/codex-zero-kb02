# Agent instructions

- Read the assigned parent GitHub issue before editing.
- One agent owns one issue. Do not edit the same child repository concurrently.
- Host work stays in `host/`; Firmware work stays in `firmware/`.
- Commit and push child changes before changing the parent gitlink.
- Use an issue branch such as `issue-5-herdr-status`; never force-push.
- Keep the product limited to Issue #34. Approve/Deny is allowed only as the
  exact fixed K9/K10 operations and verified approval conditions in Issue #25. Do not add K11 push-to-talk,
  joystick-push actions, persistent/session/prefix approval, arbitrary input or
  shell execution, Zed ACP, Vial control, or a settings GUI.
- Never flash firmware or open/control a real device unless the user explicitly
  approves that operation.
- Preserve third-party attribution. The TinyGo workshop is reference-only until
  it has an explicit reusable license; do not copy its code.
- Run the smallest focused check that proves the acceptance criteria and report
  exactly what was and was not tested.

- The product design is Issue #34; the wire contract is PROTOCOL.md major 2.
  Old major 1 remains reference-only and must never be auto-detected as major 2.
- Firmware implementation belongs to the user. Codex may prepare design docs,
  dependency patches and temporary compile-only probes; do not implement or
  modify `firmware/` without a separate request.
- TinyGo standard HID mouse is allowed for joystick movement. HID keyboard
  output, Vial controls and both push actions are outside the product scope.
- When the user starts firmware work, align firmware/AGENTS.md with the input-only
  library build tag and major 2 contract before changing its implementation.
