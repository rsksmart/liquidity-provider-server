FROM localstack/localstack:3.1.0@sha256:9d278167f2b7fda866b8bbe2a077a3fcf2e2f704b024fa76211c7c468428f977

ARG BTC_WALLET_PASSWORD=test
ARG RSK_ENCRYPTED_JSON_PASSWORD=test

RUN apt-get install --no-install-recommends -y jq && apt-get clean

COPY local-key.json /tmp/local-key.json
COPY setup/ /etc/localstack/init/ready.d/

USER localstack