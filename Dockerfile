FROM --platform=arm64 golang:1.16.3-alpine3.13 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /app/bin/generate-nginx-block ./main.go

# for testing purposes we download nginx and certbot

RUN apk add --no-cache nginx certbot

# we should expect the program to fail with certbot not being able to run since we don't want to link a domain to this container

FROM --platform=arm64 alpine:3.13

WORKDIR /app

COPY --from=builder /app/bin/generate-nginx-block /bin/generate-nginx-block

CMD ["generate-nginx-block"]