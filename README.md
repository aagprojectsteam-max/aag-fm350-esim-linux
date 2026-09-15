# AAG FM350 eSIM Linux

Experimental Linux tooling recovered from a working Fibocom FM350 / MediaTek T700 eSIM setup.

This repository documents and packages the userspace tools used to access the eUICC over MBIM on Linux, inspect slot 2, read the EID, initialize terminal capabilities, toggle passthrough, download an eSIM profile, and enable/switch profiles.

## Hardware and transport

The recovered setup used Fibocom FM350 / MediaTek T700 hardware on Ubuntu, MBIM device `/dev/wwan0mbim0`, eUICC slot 2 (one-based numbering in `euicc-go`), and `euicc-go` + `wwan-go`. Device names and slot numbering can differ on other systems.

## Included tools

The `cmd/` directory contains six recovered utilities: `eid-reader`, `slot2-diagnostic`, `terminal-capability-init`, `passthrough-toggle`, `profile-downloader`, and `profile-manager`.

`slot2-diagnostic` is intended for read-only inspection. The downloader, profile enable/notification handling, passthrough changes, and terminal-capability initialization can change modem/eUICC state. Read `docs/SAFETY.md` first.

## Build

A recent Go toolchain is required. The recovered source declares Go 1.26.3 and pins the `euicc-go` / `wwan-go` revisions used during development.

```bash
./scripts/build.sh
```

or:

```bash
make build
```

Binaries are written to `bin/`.

## Suggested progression

Start with hardware identification and read-only slot diagnostics, then read the EID, validate the eUICC/terminal-capability path, list profiles, and only then perform profile download or enable/switch operations.

Do not run state-changing tools merely as diagnostics.

## Profile download safety

The downloader contains no activation code. It reads the SGP.22 activation code at runtime, performs strict SM-DP+ TLS validation, retrieves profile metadata, displays only a masked ICCID, and requires the literal confirmation `INSTALL` before installation.

Do not put real activation codes, EIDs, ICCIDs, IMSIs, phone numbers, or account credentials into issues, logs, screenshots, commits, or bug reports.

## Relationship to lpac

This project is not a fork of lpac. The final inspected lpac tree had no local source diff; only build/install artifacts were untracked. Current lpac also contains an MBIM APDU backend. This repository records the separate Go-based integration and the FM350/T700-specific workflow recovered from the tested machine.

## GNSS

GNSS support for the same hardware family is maintained separately in the `fibocom-fm350-t700-gnss-linux` repository and is not required by these eSIM userspace tools.

## Tested environment

The original validation was performed on an HP EliteBook 840 G11 with Fibocom FM350 / MediaTek T700 under Ubuntu. The device path and slot mapping are the values actually validated there, not universal assumptions.

## Status

`v0.1.0` is the initial public-source recovery. The recovered source and documentation were sanitized before publication. Treat it as experimental on other firmware, laptops, carriers, and eUICC implementations.

A fresh network-dependent Go build could not be completed inside the publication sandbox because that environment could not download the requested Go 1.26.3 toolchain. The recovered source was formatted and privacy-audited; users should build and validate it on their own target system before relying on state-changing operations.

## Privacy

The repository intentionally excludes the original EID file, raw eSIM/MBIM logs, activation codes and matching IDs, real ICCID/IMSI values, recovery state files, and locally built binaries.

If you discover private subscriber data, follow `SECURITY.md` and do not open a public issue containing it.

## License

MIT. See `LICENSE`.
