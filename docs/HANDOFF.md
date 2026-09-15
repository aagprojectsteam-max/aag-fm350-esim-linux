# Complete Handoff: Fibocom FM350 / MediaTek T700 eSIM on Linux

## Purpose

This document hands off the complete recovered engineering path that led to practical eSIM access on Linux with a Fibocom FM350 / MediaTek T700 WWAN modem. It is intended for another Linux developer or owner of similar hardware who wants to understand what was discovered, what was actually demonstrated, what remains hardware/firmware dependent, and how the source in this repository fits together.

This is a reconstruction of the successful development work. It deliberately omits subscriber-private data, raw logs, activation codes, EID/ICCID/IMSI values, and carrier-specific credentials.

## 1. Original objective

The goal was not merely to use a profile that had already been provisioned elsewhere. The target was native Linux eSIM management:

1. access the embedded UICC directly from Linux;
2. obtain the EID;
3. identify the correct physical/logical UICC slot;
4. communicate with the eUICC over the modem's Linux transport;
5. accept an SGP.22 activation code containing an SM-DP+ address and matching ID;
6. communicate with the SM-DP+ service;
7. retrieve profile metadata;
8. download a new profile;
9. enumerate installed profiles;
10. enable/switch profiles;
11. do all of this without requiring Windows as part of the normal workflow.

The development machine was an HP EliteBook 840 G11 with the Fibocom FM350 / MediaTek T700 family modem running Ubuntu Linux.

## 2. The key architectural discovery

The useful Linux path was MBIM, not a proprietary Windows-only eSIM application.

On the tested system the relevant control device was:

```text
/dev/wwan0mbim0
```

The embedded eUICC was reachable as **UICC slot 2**. The Go library used by the recovered tools numbers slots from one, so the source uses:

```go
const esimSlot = 2
```

This is a property of the tested machine/firmware combination, not a universal rule. A different laptop, firmware build, kernel driver, or modem configuration may expose a different device name or slot arrangement.

The userspace implementation is built on:

- `github.com/damonto/euicc-go`
- `github.com/damonto/wwan-go`

The recovered revisions are pinned in `go.mod` / `go.sum`.

## 3. Why MBIM matters

The eUICC communication path ultimately uses MBIM UICC low-level access/APDU functionality. This means Linux userspace can reach the eUICC through the modem rather than needing a vendor GUI.

A parallel investigation of `lpac` also showed that current lpac code contains an MBIM APDU backend using Microsoft UICC Low Level Access. Therefore this repository is **not** an lpac fork. It preserves the separate Go implementation and the FM350/T700 workflow that was built and exercised during the investigation.

## 4. What was demonstrated during development

The work progressed beyond simple modem detection.

The recovered development state demonstrates that the implementation was built to perform the following real operations:

- communicate with the MBIM device;
- target eUICC slot 2;
- read the EID;
- query/initialize terminal capabilities needed by the tested eUICC path;
- inspect the slot and card/application state;
- list eSIM profiles;
- parse a normal SGP.22 activation code;
- contact the SM-DP+ endpoint using strict TLS validation;
- obtain proposed profile metadata;
- show the user the profile/provider information and a masked ICCID;
- require explicit confirmation before installation;
- execute the profile download flow;
- enable/switch a selected profile.

During the original work an existing mobile profile could be read and the implementation reached real SM-DP+ profile-download activity. Later modem/ModemManager/MBIM state problems were investigated separately. Those later transport-state problems are one reason this repository is published as experimental rather than claiming universal plug-and-play support.

## 5. Source recovery and privacy cleanup

The original development directory contained source trees, locally compiled Go binaries, logs, state files, and a local EID file. Before publication the material was classified and only source material was retained.

The public repository excludes:

- the original `eid.txt`;
- raw eSIM and MBIM logs;
- recovery state/lock files;
- compiled ELF binaries;
- real activation codes and matching IDs;
- real EID, ICCID, IMSI, phone-number, or subscriber-account values.

Several temporary `profile-manager-switch.*` source directories survived. Their relevant source was compared during recovery and represented the same implementation, so the public tree contains one canonical `cmd/profile-manager` rather than publishing temporary duplicate directories.

A carrier/profile-specific fingerprint that existed in the private development workflow was also removed. The public profile manager is intentionally generic and requires the operator/user to identify the desired profile rather than embedding the original subscriber's identity.

## 6. Repository map

The public tree is organized as follows:

```text
cmd/
  eid-reader/
  slot2-diagnostic/
  terminal-capability-init/
  passthrough-toggle/
  profile-downloader/
  profile-manager/
docs/
  ARCHITECTURE.md
  SAFETY.md
  TROUBLESHOOTING.md
  HANDOFF.md
scripts/
  build.sh
README.md
go.mod
go.sum
Makefile
SECURITY.md
CHANGELOG.md
LICENSE
```

## 7. Tool-by-tool handoff

### 7.1 `eid-reader`

Purpose: prove that the eUICC is actually reachable and obtain its EID.

It opens the MBIM path, selects the eUICC slot, and calls the eUICC client EID operation. Console output is masked. The historical private workflow also saved the full EID locally; the public source was changed so that a user's private identifier is not part of this repository.

Treat the EID as subscriber/device-sensitive information. Do not paste it into public issues.

### 7.2 `slot2-diagnostic`

Purpose: inspect the UICC/eUICC path before performing profile-changing operations.

Use this when porting the workflow to another FM350/T700 machine. The important question is not merely whether `/dev/wwan0mbim0` exists, but whether the expected UICC slot and application/card state can be reached through it.

### 7.3 `terminal-capability-init`

Purpose: query/set the terminal capability data required by the tested eUICC path.

The recovered machine recorded the slot-2 capability value:

```text
830107
```

This is included as engineering evidence, not as a claim that every firmware requires the identical sequence.

**This tool changes state.** Do not use it casually just to see what happens.

### 7.4 `passthrough-toggle`

Purpose: control the passthrough state used during the development workflow.

This was part of making the eUICC path usable on the tested modem configuration. It is a state-changing operation and should be used only when the developer understands the current modem state and has a recovery path.

### 7.5 `profile-downloader`

This is the core SGP.22 download utility.

Its workflow is deliberately defensive:

1. obtain the activation code at runtime rather than embedding it in source;
2. normalize common activation-code input forms into `LPA:1$...` form;
3. parse the SM-DP+ address and matching ID;
4. perform a strict TLS preflight;
5. use the trusted Root CAs supplied by `euicc-go` rather than accepting a TLS-interception CA;
6. communicate with the SM-DP+ service;
7. capture bounded ES9+/server diagnostic status without dumping secrets indiscriminately;
8. obtain the offered profile metadata;
9. display profile name, service provider, class, owner MCC/MNC when present, and only a **masked** ICCID;
10. require the operator to type `INSTALL` before committing the installation;
11. perform the download;
12. report the resulting profile using masked identifiers.

The TLS policy is important. eSIM provisioning is a security-sensitive protocol. The recovered implementation intentionally disables HTTP proxy inheritance for this client and does not add a local interception CA to the trusted set. If a filtering/proxy network breaks provisioning TLS, use a direct network path rather than weakening certificate verification.

### 7.6 `profile-manager`

Purpose: enumerate installed profiles and perform profile enable/switch operations.

The private development copy knew which test profile was expected. That subscriber-specific assumption was removed for publication. A reusable implementation must not ship somebody else's ICCID or carrier-profile fingerprint.

Enabling a profile is a state-changing operation. Depending on modem/firmware behavior, network registration may temporarily disappear while the modem/UICC state settles.

## 8. Recommended bring-up procedure on another machine

Do not begin with profile download. Reproduce the lower layers first.

### Stage A — identify hardware

Confirm that the machine really contains the expected FM350/T700-family modem and that the Linux kernel exposes its WWAN/MBIM interfaces.

Useful read-only inspection typically includes:

```bash
lspci -nn
ls -l /dev/wwan*
ls -l /sys/class/wwan
```

If the MBIM device is not present, stop here. The eSIM userspace tools cannot compensate for a missing kernel/modem transport.

### Stage B — identify ownership/conflicts

ModemManager or another process may already own the MBIM control channel. Do not randomly kill services. Determine whether the chosen tool/library can coexist with the MBIM proxy arrangement on your system and whether exclusive access is required.

The original investigation encountered later MBIM timeouts and ModemManager state failures. A timeout does not prove the eUICC disappeared; it can indicate that the modem/control channel is in a bad state or being contended.

### Stage C — run read-only diagnostics

Build the repository and begin with the slot diagnostic and EID reader. Your first success criterion is simple:

```text
Linux userspace -> MBIM -> expected UICC slot -> eUICC
```

Do not continue to provisioning until this is reliable.

### Stage D — validate terminal capability requirements

Inspect the terminal-capability behavior and compare it with the tested FM350 setup. Only perform initialization if the diagnostics show it is required and you understand the operation.

### Stage E — enumerate profiles

Confirm that installed profile metadata can be read. Never publish the full identifiers from your own eUICC while asking for help.

### Stage F — prepare the network path

Before a real profile download, make sure the provisioning endpoint can be reached without TLS interception. The downloader's strict verification is intentional.

### Stage G — perform a real profile download

Use a valid, unused activation code supplied by your carrier/provider. Do not put it on the command line if that would expose it in shell history. The recovered downloader supports secret input from the terminal/stdin path.

Carefully inspect the profile metadata displayed before typing `INSTALL`.

### Stage H — enable/switch

After successful installation, list profiles again and select the intended profile. Expect network registration to need time to recover. Do not interpret a transient registration failure as proof that the profile installation failed.

## 9. Building

The recovered source modules declare Go 1.26.3 and pin the library revisions used during development.

From the repository root:

```bash
./scripts/build.sh
```

or:

```bash
make build
```

The binaries are written under `build/` and are intentionally ignored by Git.

At publication time the source tree was audited and formatted, but the publication environment could not complete a fresh network-dependent Go toolchain/dependency download because access to the Go module proxy was unavailable. Consequently the public release does **not** claim a fresh independent build validation in that environment. Build on the target Linux system and report reproducible compiler errors without private identifiers if anything has drifted.

## 10. Important failure modes learned during the project

### MBIM timeouts

Later in the investigation ModemManager experienced timeouts on the MBIM control device and could temporarily reject the modem. This can be caused by modem/control-channel state and should be debugged separately from the SGP.22 server flow.

### Wrong slot assumptions

Slot 2 was correct on the tested system. Blindly forcing slot 2 on another firmware is unsafe engineering. Diagnose first.

### TLS interception

A network that replaces remote TLS certificates can break provisioning. Do not solve this by teaching the provisioning tool to trust the interception CA. Use an appropriate direct connection.

### Confusing profile success with network registration

Profile installation/enable and cellular registration are related but distinct layers. Diagnose them separately.

### Publishing identifiers while debugging

EID, ICCID, IMSI, activation matching IDs and raw provisioning logs can be sensitive. Mask them before sharing.

## 11. What is FM350/T700-specific and what is generic

Generic pieces:

- SGP.22 activation-code parsing;
- SM-DP+ HTTPS/TLS flow;
- profile metadata confirmation;
- eUICC profile operations exposed by `euicc-go`.

Tested-machine assumptions:

- `/dev/wwan0mbim0`;
- eUICC on one-based slot 2;
- terminal-capability behavior/value;
- passthrough behavior;
- interaction with the particular FM350/T700 firmware and Linux kernel stack.

A good future contribution would make device and slot selection configurable after automatic discovery rather than hard-coding the tested values.

## 12. Relationship to the T700 kernel/GNSS work

The same modem was also used in a separate Linux GNSS investigation involving the MediaTek `t7xx` driver and additional WWAN/AT channel mapping. That work is intentionally kept in the separate `fibocom-fm350-t700-gnss-linux` repository.

Do not assume the GNSS patches are required for this eSIM userspace implementation. Keeping the projects separate makes it easier to distinguish eSIM requirements from unrelated modem improvements.

## 13. Relationship to the original lpac experiments

The development machine also retained an `lpac-fm350-current` tree. Inspection during recovery showed no local source patch in that lpac working tree; the meaningful local untracked material was build/install output. This was an important conclusion: the public eSIM project should not pretend that a large private lpac fork is required.

The reusable contribution from this investigation is the FM350/T700 bring-up knowledge and the recovered Go tooling, while upstream lpac remains an independent alternative implementation.

## 14. Current project status

The public repository should be understood as:

```text
SOURCE_RECOVERY=COMPLETE
PRIVATE_IDENTIFIERS_REMOVED=YES
RAW_LOGS_PUBLISHED=NO
COMPILED_PRIVATE_BINARIES_PUBLISHED=NO
MBIM_TRANSPORT_TESTED_ON_ORIGINAL_MACHINE=YES
EID_ACCESS_DEMONSTRATED=YES
PROFILE_ENUMERATION_IMPLEMENTED=YES
REAL_SMDP_PLUS_FLOW_REACHED=YES
PROFILE_DOWNLOAD_IMPLEMENTED=YES
PROFILE_ENABLE_SWITCH_IMPLEMENTED=YES
UNIVERSAL_FM350_FIRMWARE_COMPATIBILITY=NOT_CLAIMED
FRESH_PUBLICATION_ENVIRONMENT_BUILD=NOT_COMPLETED_NETWORK_BLOCKED
```

## 15. Next engineering steps

For someone continuing the project, the highest-value improvements are:

1. make MBIM device and UICC slot configurable;
2. add a read-only discovery command that finds candidate MBIM devices and slots safely;
3. add automated tests for activation-code normalization and identifier masking;
4. add mock transport tests for profile operations without touching a real eUICC;
5. test multiple FM350 firmware revisions and laptop vendors;
6. document ModemManager coexistence/recovery more rigorously;
7. add CI once the required Go toolchain is generally available in the runner environment;
8. only after cross-machine evidence, promote individual steps from experimental to supported.

## 16. Rules for contributors

Never commit:

- a real EID;
- a real ICCID or IMSI;
- a usable activation code/matching ID;
- subscriber phone/account data;
- raw provisioning traces that have not been manually reviewed;
- carrier credentials or API secrets.

When reporting a problem, provide hardware/firmware/kernel versions, sanitized error text, the operation being attempted, and whether the failure is in MBIM transport, eUICC/APDU access, SM-DP+ communication, profile management, or later cellular registration.

## 17. Bottom line

The project established a practical route from Linux userspace to the FM350/T700 embedded eUICC through MBIM. The important result was not a Windows workaround: it was direct Linux access to the eUICC and an implementation of the profile-management/provisioning workflow around that transport.

The safest way for another owner of this modem to reproduce the work is to validate each layer in order — hardware, MBIM, UICC slot, EID, terminal capability, profile enumeration, SM-DP+ connectivity, download, then enable/switch — rather than jumping directly to a live activation code.

This handoff, together with the source and the shorter architecture/safety/troubleshooting documents, is intended to make that reproduction possible without requiring access to the original private development logs.