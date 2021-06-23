FROM golang:alpine

RUN mkdir /urlshortener
ADD . /urlshortener
WORKDIR /urlshortener/server

RUN go mod download
RUN go build -o server .

CMD ["./server"]