FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /go/bin/api cmd/main.go

FROM gcr.io/distroless/static-debian12
WORKDIR /root/
COPY --from=builder /go/bin/api .

EXPOSE ${PORT}
CMD ["./api"]
