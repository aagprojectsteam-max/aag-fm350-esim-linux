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
