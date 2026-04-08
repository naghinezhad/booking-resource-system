# ---------- Build Stage ----------
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# use Go module mirror
ENV GOPROXY=https://package-mirror.liara.ir/repository/go/
ENV GOSUMDB=off

# download dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server


# ---------- Runtime Stage ----------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]