FROM --platform=linux/amd64 node:16@sha256:f77a1aef2da8d83e45ec990f45df50f1a286c5fe8bbfb8c6e4246c6389705c0b

RUN apt-get update -y && apt-get install -y -qq --no-install-recommends jq && apt-get clean

RUN npm install --ignore-scripts -g truffle@5.11.5@sha512-yCa2uWs5DmL0spuJUuIMtnVayRQrVuWLtcRXHMB0NLrtWDcRo7VM9RViveV4+oi9LdZ8VpFmmqHGm43LbzUxOA==

USER node

WORKDIR /home/node

ARG LBC_GIT_BRANCH="QA-Test"

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
