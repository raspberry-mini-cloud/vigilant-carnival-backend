FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /vigilant-carnival-backend/cmd/app ./cmd/app

EXPOSE 5000

CMD ["/vigilant-carnival-backend/cmd/app"]
