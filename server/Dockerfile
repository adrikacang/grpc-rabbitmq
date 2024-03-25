# Dockerfile
FROM golang:alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod tidy

COPY . .

RUN go build -o server .

EXPOSE 50051

CMD ["./server"]