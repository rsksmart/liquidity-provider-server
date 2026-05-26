FROM node:25.7.0@sha256:18e02657e2a2cc3a87210ee421e9769ff28a1ac824865d64f74d6d2d59f74b6b

WORKDIR /app

COPY package.json ./
RUN npm install

# Make node_modules available to scripts mounted at /scripts at runtime
ENV NODE_PATH=/app/node_modules
