# Copilot CLI Error Behavior (No JSONL Output)

CLI argument validation errors (e.g. invalid --model) produce NO JSONL output.
Error goes to stderr as plain text with exit code 1.

Example stderr:
```
error: option '--model <model>' argument 'nonexistent-model' is invalid.
```

For runtime errors during execution, copilot may emit a result event with non-zero exitCode.
The session may also end abruptly with fewer events than expected.

Rate limit errors surface as text in assistant.message content, not as a separate error event type.

Note: This file replaces `error_exit.jsonl` which contained `_meta` commentary objects
that were not valid copilot JSONL events. Since CLI errors produce no JSONL, there is
no fixture to capture for this scenario.
