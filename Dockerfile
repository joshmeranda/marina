FROM golang:1.21 as builder

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /src

RUN cd /src && \
    make marina-server \
    && cp bin/marina-server /marina-server

FROM golang:1.21

RUN go install github.com/grpc-ecosystem/grpc-health-probe@latest

COPY --from=builder /marina-server /marina-server

ENTRYPOINT ["/marina-server"]
