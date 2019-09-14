FROM golang:1.13.0

WORKDIR /go/src/work
COPY . .

RUN go build -o opt/resource/check ./cmd/check && \
    go build -o opt/resource/in ./cmd/in && \
    go build -o opt/resource/out ./cmd/out

FROM gcr.io/distroless/base
COPY --from=0 /go/src/work/opt/resource /opt/resource
