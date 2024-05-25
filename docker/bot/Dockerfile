ARG GO_VERSION=1.22.2
ARG ALPINE_VERSION=3.17.2

FROM golang:${GO_VERSION}-alpine AS build
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

WORKDIR /go/src
#COPY --exclude=*.key ca/ ca/
COPY ca/ ca/
COPY cmd/ cmd/
COPY lib/ lib/
COPY res/ res/
COPY go.* ./

RUN apk add build-base && \
    go build -v -race -o ./bin/bot ./cmd/bot

FROM alpine:${ALPINE_VERSION}

ENV http_proxy http://80.78.243.67:3128
ENV https_proxy http://80.78.243.67:3128

WORKDIR /usr/local
COPY --from=build /go/src/ca/ share/ca-certificates/
COPY --from=build /go/src/bin/bot bin/
COPY --from=build /go/src/cmd/bot/config.json bin/
COPY --from=build /go/src/res/ res/

WORKDIR /usr/local/bin
EXPOSE 65080

ENTRYPOINT ["./bot", "config.json"]