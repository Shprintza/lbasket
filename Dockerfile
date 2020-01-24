###########################
# Global build args
###########################
# change service and workdir with your 
ARG service=lbasket
ARG workdir=/${service}
ARG netrcFile=./.netrc
ARG GOPRIVATE=github.com/orov-io/*
############################
# STEP 1 build executable binary
############################

FROM golang:1.13-alpine AS build-env

ARG workdir
ARG GOPRIVATE
ARG netrcFile

WORKDIR ${workdir}

# Install git.
## Git is required for fetching dependencies.
RUN apk update && apk add --no-cache git openssh-client build-base

## This line adds gomod access to private repositories
COPY ${netrcFile} /root/.netrc

COPY . ./

# Building the app
RUN go build -o app

# ENTRYPOINT [ "./app" ]

############################
# STEP 2 build tiny executable container
############################

FROM alpine

ARG workdir

WORKDIR ${workdir}
COPY --from=build-env ${workdir}/app ./
COPY --from=build-env ${workdir}/migrations/* ./migrations/
RUN apk --update add ca-certificates

EXPOSE 8080
ENTRYPOINT [ "./app" ]