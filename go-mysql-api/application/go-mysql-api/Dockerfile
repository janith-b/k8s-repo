FROM golang:alpine AS builder
WORKDIR /application
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-mysql-api

FROM scratch
COPY --from=builder /go-mysql-api /go-mysql-api
ENTRYPOINT ["/go-mysql-api"]