FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

WORKDIR /app/cmd/app

RUN CGO_ENABLED=0 GOOS=linux go build -o /vigilant-carnival-backend

EXPOSE 5000

CMD ["/vigilant-carnival-backend"]
