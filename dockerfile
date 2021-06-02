FROM golang:1.16.3-alpine
RUN apk add --no-cache git


WORKDIR /app/goapp
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .


RUN go build -o ./out/goapp .

EXPOSE 10000

CMD ["./out/goapp"]
