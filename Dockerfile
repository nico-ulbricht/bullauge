FROM golang:alpine as builder
RUN apk add git
WORKDIR /app

ENV GO1111MODULE=on
ENV GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./server ./cmd/server 




FROM alpine

WORKDIR /app
COPY --from=builder /app/server .

CMD ./server