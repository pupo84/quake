FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . /app

RUN go build -o /app/quake

FROM scratch

COPY --from=builder /app/quake /app/quake
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/log.txt /app/log.txt

WORKDIR /app

EXPOSE 8000

CMD ["./quake"]