FROM node:18.20.3@sha256:70ae166dcb03837ebe1abbb78919a1a4bdc79cc970344d56f33b6c753f4b46cc AS builder

WORKDIR /code/flyover-ui

ARG VERSION="custom-lbc"
ARG LBC_ADDRESS


RUN --mount=type=secret,id=gh_token,required=true,uid=1000 gitrev="${VERSION}" && REPO_URL="https://$(cat /run/secrets/gh_token)@github.com/rsksmart/flyover-ui.git" &&\
    git init && \
    git remote add origin "$REPO_URL" && \
    git fetch --depth 1 origin tag "$gitrev" && \
    git checkout "$gitrev"

RUN --mount=type=secret,id=gh_token,required=true,uid=1000 npm config set //npm.pkg.github.com/:_authToken $(cat /run/secrets/gh_token)
RUN npm ci

RUN echo "REACT_APP_ENV=Regtest" > .env
RUN echo "REACT_APP_RPC_URL=http://localhost:4444" >> .env
RUN echo "REACT_APP_CUSTOM_LBC=$LBC_ADDRESS" >> .env


RUN npm run build


FROM nginx:1.27.1@sha256:135fbc7ed19c8f644ddf678e68292e678696908451dad7ee2fd4e0cf861f4b6f
COPY --from=builder /code/flyover-ui/build /usr/share/nginx/html

EXPOSE 80

STOPSIGNAL SIGTERM

CMD ["nginx", "-g", "daemon off;"]
