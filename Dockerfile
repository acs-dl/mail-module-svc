FROM golang:1.19-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/acs-dl/mail-module-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/mail-module /go/src/github.com/acs-dl/mail-module-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/mail-module /usr/local/bin/mail-module
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["mail-module"]
