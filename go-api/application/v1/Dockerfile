FROM golang:alpine AS builder
WORKDIR /application
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-api

FROM scratch
COPY --from=builder /go-api /go-api
ENTRYPOINT ["/go-api"]
