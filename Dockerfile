FROM golang:1.21.0 as builder
WORKDIR /sidecar
COPY . .

RUN apt-get update && apt-get install gettext-base
RUN go install github.com/client9/misspell/cmd/misspell@v0.3.4 \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.1 \
    && go install github.com/google/addlicense@v1.0.0
RUN apt update && apt install -y make
RUN make build

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /sidecar/bin/rungmpcol /rungmpcol
COPY --from=builder /sidecar/bin/run-gmp-entrypoint /run-gmp-entrypoint

ENTRYPOINT ["/run-gmp-entrypoint"]