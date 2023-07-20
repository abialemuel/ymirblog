FROM telkomindonesia/alpine:go-1.20 AS builder
WORKDIR /usr/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main ./cmd/bin

FROM alpine:3.18
WORKDIR /usr/src/app
RUN apk add curl
RUN touch .env
COPY config.yaml .
COPY --from=builder /usr/src/app/main .
CMD ["./main"]