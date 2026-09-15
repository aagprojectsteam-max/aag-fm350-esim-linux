# Safety and privacy

## Secrets and identifiers

Treat these as private: EID, ICCID, IMSI, IMEI when linked to a person/device, activation code, matching ID, confirmation code, and raw provisioning logs.

Do not commit generated `eid.txt` or activation-code files. `.gitignore` blocks common names, but it is not a substitute for review.

## State-changing operations

The diagnostic tool is designed for read-only queries. The following operations can alter modem/eUICC state:

- downloading/installing a profile;
- enabling a profile;
- handling/removing eUICC notifications;
- changing UICC pass-through state;
- writing terminal capability state.

Do not run those merely to test whether the modem exists.

## TLS

Provisioning is intentionally strict. Do not weaken certificate validation or add an interception CA to make an SM-DP+ download pass. Use a direct connection that can validate the provisioning endpoint normally.

## Recovery principle

If MBIM begins timing out, stop state-changing experiments. Restore ordinary modem operation first, then return to read-only diagnostics.
