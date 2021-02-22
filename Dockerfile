FROM golang:1.15-alpine AS builder
WORKDIR /go/src/github.com/smvfal/faas-monitor
COPY . .
RUN go build

FROM alpine
WORKDIR /root/
COPY --from=builder /go/src/github.com/smvfal/resmonitor/faas-monitor .
CMD ["./faas-monitor"]