#TODO: rewrite Dockerfile to be minimal, what are build commands
FROM golang

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /c2-listener

EXPOSE 8080

CMD ["/c2-listener"]

