FROM golang:1.19-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/gitlab.com/distributed_lab/acs/mail-module
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/mail-module /go/src/gitlab.com/distributed_lab/acs/mail-module


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/mail-module /usr/local/bin/mail-module
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["mail-module"]
