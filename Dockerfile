# syntax = docker/dockerfile:1.2
FROM registry.access.redhat.com/ubi8/ubi as js-builder

RUN dnf makecache

WORKDIR /usr/src/nikt-link-proxy/loader

RUN dnf install -yq nodejs npm

COPY loader/package-lock.json .
COPY loader/package.json .

RUN npm ci

ARG NODE_ENV=production

COPY loader/tsconfig.json .
COPY loader/src src

RUN npm run build-loader

FROM registry.access.redhat.com/ubi8/ubi as go-builder

RUN dnf makecache

WORKDIR /usr/src/nikt-link-proxy

RUN dnf install -yq golang

ENV GOPATH /root/go

COPY loader/go.mod loader/go.mod
COPY go.mod .
COPY go.sum .

RUN --mount=type=cache,id=gopath,target=${GOPATH} \
    go mod \
      download

COPY loader loader
COPY pkg pkg
COPY --from=js-builder \
    /usr/src/nikt-link-proxy/loader/binary loader/binary

ARG APP_VERSION=containerized

RUN --mount=type=cache,id=gopath,target=${GOPATH} \
    go build \
      -o nikt-link-proxy -ldflags "-X main.Version=${APP_VERSION}" \
      pkg/cmd/main.go

FROM registry.access.redhat.com/ubi8/ubi as binary

WORKDIR /etc/nikt-link-proxy

COPY config.yaml .
COPY --from=go-builder \
    /usr/src/nikt-link-proxy/nikt-link-proxy \
    /usr/local/bin/nikt-link-proxy

CMD ["/usr/local/bin/nikt-link-proxy"]