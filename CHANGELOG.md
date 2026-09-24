# Changelog

## 0.2.0 - 2026-09-24

- Documented a full end-to-end FM350/T700 eSIM recovery and validation on Ubuntu.
- Separated a carrier/subscription registration rejection from independent local modem/ModemManager failures.
- Recorded Cause 11 and Cause 112 registration evidence and the successful A/B result with a newly created subscription, without publishing subscriber identifiers.
- Isolated the ModemManager first-enable failure to `ATZ -> +CME ERROR: 59`.
- Added the exact tested downstream ModemManager 1.25.95 patch that allows enable to continue on the observed unexpected-data-value response.
- Documented FCC, CFUN, SIM-ready, MBIM and MBIMEx v3 validation.
- Documented APDU status `910B` where profile-state readback showed the enable operation had actually taken effect.
- Documented stale `gsm.sim-id` binding after eSIM replacement.
- Isolated NetworkManager bearer `InvalidParameters` to an overly broad authentication set.
- Documented the working NetworkManager fix: `ppp.noauth=yes` while refusing PAP/CHAP/MSCHAP/MSCHAPv2/EAP.
- Verified final LTE `home` registration, `attached` packet service, IPv4 assignment, carrier DNS, ICMP, HTTPS, and disconnect/reconnect.
- Added a sanitized incident report and updated troubleshooting guidance.

## 0.1.0 - 2026-09-15

- Initial public-source preparation from the validated FM350/T700 eSIM investigation.
- Preserved MBIM slot-2 diagnostics, EID reader, profile downloader, conservative profile manager, pass-through tool, and terminal-capability initializer.
- Removed carrier-specific profile fingerprints and all recovered private identifiers/logs.
- Added safety, architecture, troubleshooting, build, and privacy documentation.
