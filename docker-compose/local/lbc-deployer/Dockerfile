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

# Branch, tag, or full commit SHA. A SHA keeps the deployed contracts in step with the LPS bindings.
ARG LBC_GIT_REF=87101fe4729eeda2b6fec9ff8cafccca0dcc6cc7
ARG LBC_GIT_URL=https://github.com/rsksmart/liquidity-bridge-contract.git
# Commit LBC_GIT_REF resolved to. lps-local.sh passes it, so a moved branch or tag rebuilds this
# layer instead of reusing a cached clone of an older commit.
ARG LBC_GIT_SHA=

# git fetch accepts a branch, a tag, or a commit SHA; git clone --branch does not take a SHA.
RUN git init -q /home/node/lbc && \
    git -C /home/node/lbc remote add origin "${LBC_GIT_URL}" && \
    git -C /home/node/lbc fetch -q --depth 1 origin "${LBC_GIT_REF}" && \
    git -C /home/node/lbc checkout -q FETCH_HEAD && \
    git -C /home/node/lbc submodule update -q --init --recursive --depth 1 && \
    HEAD_SHA="$(git -C /home/node/lbc rev-parse HEAD)" && \
    echo "LBC ${LBC_GIT_REF} at ${HEAD_SHA}" && \
    if [ -n "${LBC_GIT_SHA}" ] && [ "${HEAD_SHA}" != "${LBC_GIT_SHA}" ]; then \
      echo "ERROR: expected LBC ${LBC_GIT_SHA}, cloned ${HEAD_SHA}"; exit 1; \
    fi

WORKDIR /home/node/lbc

RUN npm ci --ignore-scripts && forge build
