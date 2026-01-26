FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /proxy ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /proxy /app/proxy
COPY configs/ /app/configs/
COPY log/ /app/log/

EXPOSE 80

CMD ["./proxy"]
