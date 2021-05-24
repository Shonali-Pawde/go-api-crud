FROM golang:1.16.3-alpine

RUN mkdir /go-crud

WORKDIR /go-crud

COPY . /go-crud

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 10000

ENTRYPOINT [ "/go-crud" ]