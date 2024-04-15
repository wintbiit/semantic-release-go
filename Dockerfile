ARG TAGS="analyzer_angular output_teamcity"
ARG VERSION="0.0.0"
ARG TZ="Asia/Shanghai"

FROM golang:1.21-alpine AS builder
ARG TAGS
ARG VERSION
WORKDIR /app
COPY . .
RUN go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -tags "${TAGS}" -o ./bin/semantic-release .

FROM alpine:3.14
ARG TZ
WORKDIR /app
ENV TZ=${TZ}
COPY --from=builder /app/bin/semantic-release /usr/local/bin/semantic-release
RUN apk add --no-cache git alpine-conf && \
    setup-timezone -z ${TZ} && \
    apk del alpine-conf && \
    rm -rf /var/cache/apk/* && \
    chmod +x /usr/local/bin/semantic-release

VOLUME /app
ENTRYPOINT ["/usr/local/bin/semantic-release"]