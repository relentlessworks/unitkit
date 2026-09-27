FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o unitkit ./cmd/unitkit

FROM scratch
COPY --from=builder /app/unitkit /unitkit
EXPOSE 8080
ENTRYPOINT ["/unitkit"]
