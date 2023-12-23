FROM golang:1.20

WORKDIR /app

COPY . .
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /strix-server

EXPOSE 7777

CMD ["/strix-server"]