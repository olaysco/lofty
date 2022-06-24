FROM golang:1.18.3-alpine AS builder
RUN mkdir /build
WORKDIR /build
COPY . .
RUN go build -o main application.go

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
WORKDIR /build
COPY --from=builder /build/main .
COPY db/migration ./db/migration
COPY .env .
CMD ["./main"]