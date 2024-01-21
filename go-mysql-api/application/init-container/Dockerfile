FROM golang:alpine AS builder
WORKDIR /application
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-init

FROM scratch
COPY --from=builder /go-init /go-init
ENTRYPOINT ["/go-init"]