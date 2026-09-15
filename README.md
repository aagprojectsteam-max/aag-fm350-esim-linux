# AAG FM350 eSIM Linux

Experimental Linux tooling recovered from a working Fibocom FM350 / MediaTek T700 eSIM setup.

This repository documents and packages the userspace tools used to access the eUICC over MBIM on Linux, inspect slot 2, read the EID, initialize terminal capabilities, toggle passthrough, download an eSIM profile, and enable/switch profiles.

## Hardware and transport

The recovered setup used:

- Fibocom FM350 / MediaTek T700 WWAN hardware
- Ubuntu Linux
- MBIM device: `/dev/wwan0mbim0`
- eUICC: slot 2 (one-based numbering in `euicc-go`)
- `euicc-go` + `wwan-go`

Device names and slot numbering can differ on other systems. Do not assume these values are universal.

## What is included

The `cmd/` directory contains six recovered Go utilities:

- `eid-reader` — reads the eUICC EID and prints only a masked form; the original private EID is not included in this repository.
- `slot2-diagnostic` — diagnostics for the second UICC/eSIM slot.
- `terminal-capability-init` — initializes terminal capabilities required by the tested eUICC path.
- `passthrough-toggle` — controls the passthrough state used during the recovered workflow.
- `profile-downloader` — accepts an SGP.22 activation code, performs strict SM-DP+ TLS validation, shows profile metadata, asks for an explicit `INSTALL` confirmation, and downloads the profile.
- `profile-manager` — lists profiles and performs profile enable/switch operations.

See `docs/ARCHITECTURE.md` and `docs/SAFETY.md` before running anything that changes modem/eUICC state.

## Build

A recent Go toolchain is required. The recovered modules declare Go 1.26.3 and pin the exact `euicc-go` / `wwan-go` revisions used during development.

```bash
./scripts/build.sh
```

or:

```bash
make build
```

Binaries are written to `build/`.

## Suggested safe progression

Start with read-only inspection. Confirm that the MBIM device exists and that your hardware/slot layout matches the documented setup before attempting any state-changing operation.

A sensible progression is:

```text
hardware identification
        ↓
slot diagnostics
        ↓
EID read
        ↓
terminal-capability validation/initialization
        ↓
profile listing
        ↓
profile download / enable / switch
```

Do not run profile download, profile enable/switch, passthrough changes, or terminal-capability initialization merely as a diagnostic test. Those operations may change modem/eUICC state.

## Profile download safety

The downloader intentionally does not contain an activation code. It reads the activation code at runtime and supports normal SGP.22 `LPA:1$...` input.

Before installation it:

1. validates the SM-DP+ endpoint using the trusted roots supplied by `euicc-go`;
2. retrieves profile metadata;
3. displays the provider/profile information and a masked ICCID;
4. requires the user to type `INSTALL` explicitly.

Do not put real activation codes, EIDs, ICCIDs, IMSIs, phone numbers, or account credentials into issues, logs, screenshots, commits, or bug reports.

## Relationship to lpac

This project is not a fork of lpac. Current lpac already contains an MBIM APDU backend using Microsoft UICC Low Level Access. The code here records a separate Go-based integration and the FM350/T700-specific workflow recovered from the tested machine.

## GNSS

GNSS support for the same FM350/T700 family is maintained separately in the `fibocom-fm350-t700-gnss-linux` project. GNSS is not required by the eSIM userspace tools in this repository.

## Status

`v0.1.0` is a source-recovery/publication release. The source was recovered from the working development machine and sanitized before publication. It should be treated as experimental on other firmware, laptops, carriers, and eUICC implementations.

The publication environment could format and audit the recovered source but could not complete a fresh network-dependent Go dependency/toolchain download. Users should therefore build and validate on their own target system before relying on it.

## Privacy

The repository intentionally excludes:

- the original EID file;
- raw eSIM/MBIM logs;
- activation codes and matching IDs;
- real ICCID/IMSI values;
- recovery state files;
- locally built binaries.

If you discover private subscriber data in the repository, follow `SECURITY.md` and do not open a public issue containing that data.

## License

MIT. See `LICENSE`.
