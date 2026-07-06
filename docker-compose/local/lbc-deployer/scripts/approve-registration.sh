#!/bin/bash
set -euo pipefail

# Capture the orchestrator-injected decision BEFORE sourcing the env file:
# the mounted env file also defines LP_REGISTRATION_DECISION and would otherwise
# clobber the value this container was started with.
DECISION="${LP_REGISTRATION_DECISION:-approve}"

# shellcheck disable=SC1090
source "/${ENV_FILE}"   # fresh DISCOVERY_ADDRESS (deploy just rewrote it)

LP="$LIQUIDITY_PROVIDER_RSK_ADDR"
ADMIN_KEY=$(cast wallet derive-private-key "$DEPLOYER_MNEMONIC")
state() { cast call "$DISCOVERY_ADDRESS" "getRegistrationState(address)(uint8)" "$LP" \
  --rpc-url "$RSK_ENDPOINT" | tr -d '[:space:]'; }

# wait for a Pending request (LPS not up yet -> starts at 0 None); bounded
S=""
for _ in $(seq 1 300); do
  S=$(state)
  [[ "$S" == "1" ]] && break
  [[ "$S" == "2" ]] && { echo "Already approved; nothing to do."; exit 0; }
  sleep 2
done
[[ "$S" == "1" ]] || { echo "Timed out waiting for a Pending registration (state=$S)"; exit 1; }

FN=$([[ "$DECISION" == "reject" ]] && echo rejectRegistration || echo approveRegistration)
echo "Admin applying '$DECISION' to $LP"
cast send "$DISCOVERY_ADDRESS" "$FN(address)" "$LP" \
  --private-key "$ADMIN_KEY" --rpc-url "$RSK_ENDPOINT" --legacy

S=$(state)
echo "Final registration state: $S"
[[ "$DECISION" == "reject" && "$S" == "3" ]] && exit 0
[[ "$DECISION" != "reject" && "$S" == "2" ]] && exit 0
exit 1
