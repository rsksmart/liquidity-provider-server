FROM --platform=linux/amd64 node:20-bookworm@sha256:2c3f34d2d28e4c13b26f7244c653527d15544626e85b1a21fb67a95ba4df9a01

RUN apt-get update -y && apt-get install -y -qq --no-install-recommends jq && apt-get clean

RUN npm install --ignore-scripts -g truffle

USER node

WORKDIR /home/node

ARG LBC_GIT_BRANCH="master"

RUN gitBranch=${LBC_GIT_BRANCH} && \
    git init && \
    git remote add origin https://github.com/rsksmart/liquidity-bridge-contract.git && \
    git fetch --depth 1 origin "$gitBranch" && \
    git checkout "$gitBranch"

COPY --chown=root:root --chmod=444 truffle-config.patch ./

RUN git apply truffle-config.patch
RUN --mount=type=secret,id=gh_token,required=true,uid=1000 echo "//npm.pkg.github.com/:_authToken=$(cat /run/secrets/gh_token)" >> .npmrc

RUN npm ci --ignore-scripts

RUN npx truffle compile

COPY --chown=root:root --chmod=555 deploy-lbc.sh ./
