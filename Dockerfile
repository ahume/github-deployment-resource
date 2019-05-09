FROM golang:1.11-alpine
RUN apk update && \
    apk upgrade && \
    apk add git

ADD . /go/src/github.com/ahume/github-deployment-resource
RUN go install github.com/ahume/github-deployment-resource

WORKDIR /go/src/github.com/ahume/github-deployment-resource
RUN /bin/ash ./scripts/build
