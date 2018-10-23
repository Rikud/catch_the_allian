FROM golang:latest
MAINTAINER Ivan Nemshilov
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]