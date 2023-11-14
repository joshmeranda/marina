FROM golang:1.21 as builder

COPY . /src

RUN cd /src && \
    make marina-server \
    && cp bin/marina-server /marina-server

FROM golang:1.21

COPY --from=builder /marina-server /marina-server

ENTRYPOINT ["/marina-server"]