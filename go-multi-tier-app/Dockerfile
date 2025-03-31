FROM golang:alpine AS builder
WORKDIR /application
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-multi-tier-app

FROM scratch
COPY --from=builder /go-multi-tier-app /go-multi-tier-app
ENTRYPOINT ["/go-multi-tier-app"]