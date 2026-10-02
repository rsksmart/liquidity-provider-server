FROM node:20.15.1@sha256:6326b52a508f0d99ffdbfaa29a69380321b215153db6f32974835bac71b38fa4

RUN apt-get update -y && \
    apt-get install -y -qq --no-install-recommends jq make curl git ca-certificates && \
    apt-get clean

USER node

ENV PATH="/home/node/.foundry/bin:${PATH}"

# Same Foundry version as the liquidity-bridge-contract Dockerfile.
ARG FOUNDRY_VERSION=v1.7.1
RUN curl -L https://foundry.paradigm.xyz | bash && \
    /home/node/.foundry/bin/foundryup --install "${FOUNDRY_VERSION}"

ARG LBC_GIT_REF=fly-2699-wire-flyover-configurations
ARG LBC_GIT_URL=https://github.com/rsksmart/liquidity-bridge-contract.git
# Commit LBC_GIT_REF resolved to. lps-local.sh passes it, so a moved ref rebuilds this layer
# instead of reusing a cached clone of an older commit.
ARG LBC_GIT_SHA=

RUN git clone --depth 1 --recurse-submodules --shallow-submodules \
      --branch "${LBC_GIT_REF}" "${LBC_GIT_URL}" /home/node/lbc && \
    HEAD_SHA="$(git -C /home/node/lbc rev-parse HEAD)" && \
    echo "LBC ${LBC_GIT_REF} at ${HEAD_SHA}" && \
    if [ -n "${LBC_GIT_SHA}" ] && [ "${HEAD_SHA}" != "${LBC_GIT_SHA}" ]; then \
      echo "ERROR: expected LBC ${LBC_GIT_SHA}, cloned ${HEAD_SHA}"; exit 1; \
    fi

WORKDIR /home/node/lbc

RUN npm ci --ignore-scripts && forge build
