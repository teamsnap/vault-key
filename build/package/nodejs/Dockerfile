FROM node:8.16.1-stretch

ENV \
 GOARCH="amd64" \
 PATH="$PATH:/usr/local/go/bin" \
 GO111MODULE=on

RUN curl -Lo /tmp/go1.13.1.linux-amd64.tar.gz https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz \
 && tar -C /usr/local -xzf /tmp/go1.13.1.linux-amd64.tar.gz \
 && rm /tmp/go1.13.1.linux-amd64.tar.gz \
 && npm install -g node-gyp@latest

WORKDIR /app

COPY . .

CMD make build
