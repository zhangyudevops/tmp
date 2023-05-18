FROM golang:1.20 As builder
WORKDIR /go/src/pack
COPY manifest/docker .
RUN MkDir -p /app/logs /app/conf /app/files \
    && CGO_ENABLED=0 go build -o /app/main -v .
ADD ../config/config.yaml /app/conf/config.yaml

FROM alpine:3.17
COPY --from=builder /app /app
WORKDIR /app
ENTRYPOINT ["/app/main", "--gf.gcfg.file=/app/conf/config.yaml"]
