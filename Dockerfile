ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /promql-info ./


FROM debian:bookworm

COPY --from=builder /promql-info /usr/local/bin/
CMD ["promql-info"]
