# This is a temporary solution until we update the LBC docker image in the registry
FROM node:20.15.1@sha256:6326b52a508f0d99ffdbfaa29a69380321b215153db6f32974835bac71b38fa4

ARG DEPLOYER_MNEMONIC
ARG LBC_COMMIT

# Install Foundry and required tools
RUN apt-get update -y && \
    apt-get install -y -qq --no-install-recommends jq make curl && \
    apt-get clean && \
    curl -L https://foundry.paradigm.xyz | bash && \
    /root/.foundry/bin/foundryup && \
    cp -r /root/.foundry /home/node/.foundry && \
    chown -R node:node /home/node/.foundry

ENV PATH="/home/node/.foundry/bin:${PATH}"
ENV WORKDIR=/home/node

WORKDIR $WORKDIR

RUN git clone https://github.com/rsksmart/liquidity-bridge-contract.git lbc

WORKDIR "$WORKDIR/lbc"

RUN git checkout ${LBC_COMMIT} && \
    make install && \
    npm ci --ignore-scripts && \
    npm run compile
