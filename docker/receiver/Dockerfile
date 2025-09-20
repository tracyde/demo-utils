FROM golang:1.25-alpine
WORKDIR /demo

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/receiver ./cmd/receiver

CMD ["/demo/bin/receiver"]
EXPOSE 8081
