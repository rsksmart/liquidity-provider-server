FROM node:16

WORKDIR /code/lbc

RUN npm install -g truffle

RUN gitbranch=master && \
    git init && \
    git remote add origin https://github.com/rsksmart/liquidity-bridge-contract.git && \
    git fetch --depth 1 origin "$gitbranch" && \
    git checkout "$gitbranch"

RUN npm ci

RUN npx truffle compile

RUN sed -i "s/host: '127.0.0.1'/host: 'rskj'/" truffle-config.js

COPY deploy.sh ./

ENTRYPOINT [ "bash", "deploy.sh" ]
