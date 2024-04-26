FROM golang:1.22-alpine

# necessary dependencies for SQLite3
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# build with CGO enabled
RUN CGO_ENABLED=1 go build -o main .

EXPOSE 8080

CMD ["./main"]