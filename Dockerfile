FROM golang:1.18-alpine AS go-builder
COPY ./ /homebrew-exporter
WORKDIR /homebrew-exporter
RUN go mod download
RUN GOOS=linux go build -o exporter .

FROM alpine
EXPOSE 9888
COPY --from=go-builder /homebrew-exporter/exporter /exporter
ENTRYPOINT ["/exporter"]
