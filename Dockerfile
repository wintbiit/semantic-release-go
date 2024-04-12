ARG TAGS="analyzer_angular output_teamcity"

FROM golang:1.21-alpine AS builder
ARG TAGS
WORKDIR /app
COPY . .
RUN go build -trimpath -ldflags "-s -w" -tags "$TAGS" -o ./bin/semantic-release .

FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/bin/semantic-release /usr/local/bin/semantic-release

VOLUME /app
ENTRYPOINT ["/usr/local/bin/semantic-release"]