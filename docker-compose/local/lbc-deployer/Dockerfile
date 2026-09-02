# This is a temporary solution until we update the LBC docker image in the registry
FROM node:26.8.1@sha256:9f94d34c787165dca03b74e5bf9c3bf90e8de79b19aa3d87fe1fa1694bf75c89

ARG DEPLOYER_MNEMONIC
ARG LBC_COMMIT
ARG FOUNDRY_VERSION=1.7.0

# Install Foundry and required tools
RUN apt-get update -y && \
    apt-get install -y -qq --no-install-recommends jq make curl && \
    apt-get clean && \
    curl -L https://foundry.paradigm.xyz | bash && \
    /root/.foundry/bin/foundryup --install "${FOUNDRY_VERSION}" && \
    cp -r /root/.foundry /home/node/.foundry && \
    chown -R node:node /home/node/.foundry

ENV PATH="/home/node/.foundry/bin:${PATH}"
ENV WORKDIR=/home/node

WORKDIR $WORKDIR

RUN git config --global http.version HTTP/1.1 && \
    git clone https://github.com/rsksmart/liquidity-bridge-contract.git lbc

WORKDIR "$WORKDIR/lbc"

RUN git checkout ${LBC_COMMIT} && \
    make install && \
    npm ci --ignore-scripts && \
    npm run compile
