FROM golang:1.20.5

WORKDIR /carchi

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

ENTRYPOINT ["./main", "-mode=server"]
