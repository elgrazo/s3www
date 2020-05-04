FROM golang:1.14 AS build-env

ADD go.mod /go/src/github.com/elgrazo/s3www/go.mod
ADD go.sum /go/src/github.com/elgrazo/s3www/go.sum
WORKDIR /go/src/github.com/elgrazo/s3www/
# Get dependencies - will also be cached if we won't change mod/sum
RUN go mod download

ADD . /go/src/github.com/elgrazo/s3www/
WORKDIR /go/src/github.com/elgrazo/s3www/

ENV CGO_ENABLED 0

RUN go build -ldflags '-w -s' -a -o s3www .

FROM alpine
EXPOSE 8080

COPY --from=build-env /go/src/github.com/elgrazo/s3www/s3www /s3www

CMD ["/s3www"]
