FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY .config /.config

ENV PORT=3000

# it needs to define in ci/cd pipeline
ENV APP_ENV=dev

EXPOSE $PORT

CMD ["./app"]