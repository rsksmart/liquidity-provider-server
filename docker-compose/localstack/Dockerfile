FROM localstack/localstack:3.1.0

ARG BTC_WALLET_PASSWORD=test
ARG RSK_ENCRYPTED_JSON_PASSWORD=test

COPY local-key.json /tmp/local-key.json
COPY setup/ /etc/localstack/init/ready.d/

USER localstack