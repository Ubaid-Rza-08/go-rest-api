# =============================
# BUILD STAGE
# =============================

FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/server/main.go

# =============================
# RUN STAGE
# =============================

FROM alpine:latest

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/main .

COPY .env .

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8000

CMD ["./main"]