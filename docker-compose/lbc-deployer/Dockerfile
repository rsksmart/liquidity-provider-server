FROM --platform=linux/amd64 node:16

RUN apt-get update -y && apt-get install -qq --no-install-recommends jq

RUN npm install -g truffle

USER node

WORKDIR /home/node

ARG LBC_GIT_BRANCH="feature/peg_out"


# If you want to get liquidity-bridge-contract from github uncomment the first patch
# else clone your project liquidity-bridge-contract into /docker-compose/lbc-deployer/liquidity-bridge-contract/

#RUN gitBranch=${LBC_GIT_BRANCH} && \
#    git init && \
#    git remote add origin https://github.com/rsksmart/liquidity-bridge-contract.git && \
#    git fetch --depth 1 origin "$gitBranch" && \
#    git checkout "$gitBranch"


# this is to simplify the way to interact with liquidity-bridge-contract
COPY --chown=node liquidity-bridge-contract ./

COPY --chown=node truffle-config.patch ./

RUN git apply truffle-config.patch

RUN npm ci

RUN npx truffle compile

COPY --chown=node deploy-lbc.sh ./
