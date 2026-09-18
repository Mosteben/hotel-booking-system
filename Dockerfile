FROM golang:1.26.3-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s" -o /app/server ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app/server /app/server

ENV PORT=10000
EXPOSE 10000

USER 1000:1000
CMD ["/app/server"]
