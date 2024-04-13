ARG TAGS="analyzer_angular output_teamcity"
ARG VERSION="0.0.0"

FROM golang:1.21-alpine AS builder
ARG TAGS
ARG VERSION
WORKDIR /app
COPY . .
RUN go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -tags "${TAGS}" -o ./bin/semantic-release .

FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/bin/semantic-release /usr/local/bin/semantic-release
RUN apk add --no-cache git

VOLUME /app
ENTRYPOINT ["/usr/local/bin/semantic-release"]