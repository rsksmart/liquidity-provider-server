# TODO: stop building from source once a liquidity-bridge-contract image with the v3 contracts is published
FROM node:20.15.1@sha256:6326b52a508f0d99ffdbfaa29a69380321b215153db6f32974835bac71b38fa4

RUN apt-get update -y && \
    apt-get install -y -qq --no-install-recommends jq make curl git ca-certificates && \
    apt-get clean

USER node

ENV PATH="/home/node/.foundry/bin:${PATH}"

ARG FOUNDRY_VERSION=v1.7.1
RUN curl -L https://foundry.paradigm.xyz | bash && \
    /home/node/.foundry/bin/foundryup --install "${FOUNDRY_VERSION}"

ARG LBC_GIT_REF=25f16a9b07f8b86ecb9dd944cc81e01866fa161e
ARG LBC_GIT_URL=https://github.com/rsksmart/liquidity-bridge-contract.git
ARG LBC_GIT_SHA=

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
