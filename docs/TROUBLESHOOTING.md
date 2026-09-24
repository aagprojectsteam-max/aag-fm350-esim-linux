# Troubleshooting

## `/dev/wwan0mbim0` is missing

Confirm the kernel driver and WWAN device are present before using these tools. This repository does not install or patch the T700 kernel driver.

## MBIM open fails

The recovered tools intentionally use the MBIM proxy. This avoids directly taking the device away from ModemManager. Verify `mbim-proxy`/ModemManager are healthy rather than changing the tools to direct-open as a first response.

## Slot 2 does not work

Slot 2 is the mapping validated on the tested FM350/T700 machine. Other firmware/platform combinations may differ. Do not blindly switch slots on production hardware; collect read-only evidence first.

## Subscriber-ready state is `no-esim-profile`

That can be a valid state for an empty eUICC and was accepted by the diagnostic path.

## Profile download fails at TLS preflight

Do not bypass TLS verification. Check whether the network performs TLS interception/filtering and retry on a direct connection.

## Profile manager refuses multiple profiles

That is intentional. The recovered helper was designed around a controlled single-profile state. Use a full LPA for multi-profile selection rather than guessing which ICCID should be changed.


## FM350 loops `disabled -> enabling -> disabled` with `Unexpected data value`

Enable ModemManager DEBUG logging and capture the first failing command. On the validated FM350 firmware, the failure was:

```text
ATZ
+CME ERROR: 59
```

while `AT`, `CFUN=1`, SIM READY and FCC-unlocked checks were healthy.

A tested downstream ModemManager 1.25.95 workaround is documented in `docs/INCIDENT-2026-09-24.md` and stored under `patches/`. The tested patch modifies a generic ModemManager path and is not upstream-ready; do not deploy it blindly to unrelated modems.

## Profile enable reports APDU status `910B`

Do not immediately retry or delete the profile. Re-list the eUICC profiles and inspect the actual profile state first. During the validated recovery, an enable request returned `910B` but the following readback showed the profile was already enabled.

## Registration fails with Cause 112

`apn-restriction-value-incompatible-with-active-pdp-context (112)` is a network registration rejection, not proof that the Linux MBIM transport is broken.

Keep modem transport recovery and carrier-registration debugging separate. Record the exact PLMN and rejection cause, stop automatic retry loops, and compare against a clean subscriber/provisioning state if one is available.

In the 2026-09-24 validation, a newly created subscription registered `home` and `attached` on the same FM350/T700 setup after the previous subscription repeatedly returned Cause 112. The exact carrier-side backend fault was not observable from the client.

## NetworkManager says no suitable WWAN device after replacing the eSIM

Check `gsm.sim-id` in the NetworkManager profile. A profile bound to the previous SIM identifier will no longer match the replacement eSIM.

Update the local profile to the current SIM identifier or create a clean GSM profile. Never publish the real SIM identifier.

## NetworkManager bearer fails with `MBIM status error: InvalidParameters`

First prove registration is already healthy:

```text
registration-state: home (or roaming)
packet-service-state: attached
network-rejection: none
```

Then compare NetworkManager with an explicit ModemManager connection:

```bash
mmcli -m <modem> --simple-connect="apn=<APN>,ip-type=ipv4,allowed-auth=none"
```

If that succeeds while NetworkManager fails, inspect the failed bearer with `mmcli -b <id>`.

On the tested NetworkManager 1.54.3 system, failed bearers advertised every authentication method. Constraining the profile to no authentication fixed the MBIM bearer:

```bash
nmcli connection modify "<WWAN profile>" \
  ppp.noauth yes \
  ppp.refuse-eap yes \
  ppp.refuse-pap yes \
  ppp.refuse-chap yes \
  ppp.refuse-mschap yes \
  ppp.refuse-mschapv2 yes
```

Disconnect and reconnect once, then verify cellular-source ping and HTTPS.

## MBIM tests time out or return permission errors

Interpret the failure in context:

- `Permission denied` means the diagnostic itself did not have access to the device.
- direct-open tests can conflict with ModemManager.
- proxy tests can still be affected by a wedged proxy/client.

Prefer bounded, read-only queries and verify `mbim-proxy`, ModemManager and device permissions before concluding that MBIM is broken.

A successful health query on the validated system was:

```bash
sudo mbimcli -d /dev/wwan0mbim0 --device-open-proxy --query-device-caps
```

## Avoid overlapping automatic and manual connection attempts

Disable NetworkManager WWAN autoconnect during controlled registration/bearer experiments. Parallel `nmcli`, `mmcli` and automatic retries can create misleading timeouts and multiple stale bearers.

Restore normal autoconnect only after the controlled test is complete.
