########################
# BASE
########################
FROM golang:1.24.4-alpine as base

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ARG APP_UID=1000
ARG APP_GID=1000
RUN addgroup -S app -g ${APP_GID} && adduser -u ${APP_UID} -S -D -G app app

RUN apk update \
    && apk add --no-cache bash ca-certificates tzdata curl \
    && update-ca-certificates

ENV TZ="Europe/Paris"

COPY ./app /app
WORKDIR /app
COPY ./docker/config/env.local /vault/secrets/.env

RUN go mod download && go mod verify

########################
# BUILD
########################
FROM base as build-env

RUN go build -ldflags="-w -s" -o /http-broadcaster

########################
# PROD ENV ###
########################
FROM alpine:3.22 as prod

ARG APP_UID=1000
ARG APP_GID=1000
RUN addgroup -S app -g ${APP_GID} && adduser -u ${APP_UID} -S -D -G app app

RUN apk update \
    && apk add --no-cache bash ca-certificates tzdata curl \
    && update-ca-certificates

ENV TZ="Europe/Paris"
COPY --from=build-env /http-broadcaster /usr/local/bin/http-broadcaster
RUN chmod +x /usr/local/bin/http-broadcaster
RUN mkdir /app && chown ${APP_UID}:${APP_GID} /app

USER app

########################
# DEV
########################
FROM base as dev

COPY --from=build-env /http-broadcaster /usr/local/bin/http-broadcaster
RUN chmod +x /usr/local/bin/http-broadcaster

ENTRYPOINT ["http-broadcaster"]
