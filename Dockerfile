FROM golang:alpine AS base
RUN apk add --no-cache git ca-certificates \
  && update-ca-certificates

FROM base as development
ARG BIN_EXT
ARG BIN_NAME=proxy-gzip
ARG CGO_ENABLED=0
ARG GO111MODULE=on
ARG GOARCH=amd64
ARG GOOS=linux
WORKDIR /go/src
COPY . /go/src
RUN go build -ldflags "-extldflags -static" -a -o /go/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT}
RUN sha256sum /go/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT} | cut -d ' ' -f 1 > /go/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT}.sha256

FROM scratch AS production
ARG BIN_EXT
ARG BIN_NAME=proxy-gzip
ARG GOARCH=amd64
ARG GOOS=linux
COPY --from=development /etc/ssl/certs /etc/ssl/certs
COPY --from=development /go/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT} /bin/proxy-gzip
COPY --from=development /go/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT}.sha256 /bin/proxy-gzip.sha256
ENTRYPOINT ["/bin/proxy-gzip"]
