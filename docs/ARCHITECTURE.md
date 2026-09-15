# Architecture

The working path is:

`tool -> wwan-go/mbim -> mbim-proxy -> /dev/wwan0mbim0 -> UIM slot 2 -> euicc-go MBIM channel -> eUICC ISD-R`

`profile-downloader` adds the SGP.22 LPA and HTTPS path to the SM-DP+ server. It performs a strict TLS preflight using the trusted root set supplied by the eUICC library, disables HTTP proxy inheritance for that client, captures bounded ES9+ status information, shows profile metadata, and requires the literal confirmation `INSTALL` before loading the bound profile package.

The original lpac tree was inspected separately. Its working tree had no source diff; only an untracked install/build directory was present. Therefore this repository publishes the recovered integration code rather than pretending to contain a custom lpac fork.

## Tools

- `slot2-diagnostic`: read-only queries for MBIM version, subscriber-ready state, pass-through status, ATR, and terminal capability count.
- `eid-reader`: reads the EID and writes it to a user-selected mode-0600 file; console output is masked.
- `profile-downloader`: profile download/install workflow. State-changing.
- `profile-manager`: conservative single-profile status/enable/notification helper. Enabling and notification handling are state-changing.
- `passthrough-toggle`: changes UICC pass-through state. State-changing.
- `terminal-capability-init`: disables pass-through when necessary and initializes terminal capability. State-changing.
