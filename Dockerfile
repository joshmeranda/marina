FROM golang:1.22 as builder

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

RUN go install github.com/grpc-ecosystem/grpc-health-probe@v0.4.26

COPY . /src

RUN cd /src && \
    make gateway \
    && cp bin/gateway /gateway

# use smaller base image
FROM golang:1.22

COPY --from=builder /gateway /gateway
COPY --from=builder /go/bin/grpc-health-probe /grpc-health-probe

ENTRYPOINT ["/gateway"]
