############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/plan-crypto/node
COPY . .

# Инсталлим
RUN sh ./scripts/install.sh

# Запускаем chto-to
ENTRYPOINT ["sh", "./scripts/start.sh"]