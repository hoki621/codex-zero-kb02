# Agent instructions

- Read the assigned parent GitHub issue before editing.
- One agent owns one issue. Do not edit the same child repository concurrently.
- Host work stays in `host/`; Firmware work stays in `firmware/`.
- Commit and push child changes before changing the parent gitlink.
- Use an issue branch such as `issue-5-herdr-status`; never force-push.
- Keep v1 limited to the issue scope. Do not add Approve/Deny, push-to-talk,
  reasoning-level control, Zed ACP, HID/Vial, a settings GUI, or arbitrary shell
  execution.
- Never flash firmware or open/control a real device unless the user explicitly
  approves that operation.
- Preserve third-party attribution. The TinyGo workshop is reference-only until
  it has an explicit reusable license; do not copy its code.
- Run the smallest focused check that proves the acceptance criteria and report
  exactly what was and was not tested.
