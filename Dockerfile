FROM golang:1.19

WORKDIR /usr/app/src
COPY . .

RUN go mod download

EXPOSE 8080
CMD ["go", "run", "main.go"]