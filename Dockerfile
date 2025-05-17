# Use official Go image
FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o chesso .

EXPOSE 8080

CMD ["./chesso"]

