ARG GO_VERSION=1

FROM node:20-alpine as npm-builder

WORKDIR /usr/src/app
COPY *.mjs *.js package.json package-lock.json ./
COPY templates ./templates

RUN npm install
RUN npm run build

FROM golang:${GO_VERSION}-bookworm as go-builder

WORKDIR /usr/src/app

COPY go.mod go.sum main.go example.test ./
COPY templates ./templates
COPY vendor ./vendor
COPY --from=npm-builder /usr/src/app/static ./static

RUN go mod download && go mod verify
RUN go build -v -o /promql-info ./

FROM debian:bookworm

COPY --from=go-builder /promql-info /usr/local/bin/
CMD ["promql-info"]
