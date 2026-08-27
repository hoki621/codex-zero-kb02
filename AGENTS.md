# Agent instructions

- Read the assigned parent GitHub issue before editing.
- One agent owns one issue. Do not edit the same child repository concurrently.
- Host work stays in `host/`; Firmware work stays in `firmware/`.
- Commit and push child changes before changing the parent gitlink.
- Use an issue branch such as `issue-5-herdr-status`; never force-push.
- Keep v1 limited to the issue scope. Approve/Deny is allowed only as the exact
  fixed K9/K10 operations defined by Issue #17. Do not add K11 push-to-talk,
  joystick-push actions, persistent/session/prefix approval, arbitrary input or
  shell execution, Zed ACP, Vial control, or a settings GUI.
- Never flash firmware or open/control a real device unless the user explicitly
  approves that operation.
- Preserve third-party attribution. The TinyGo workshop is reference-only until
  it has an explicit reusable license; do not copy its code.
- Run the smallest focused check that proves the acceptance criteria and report
  exactly what was and was not tested.
