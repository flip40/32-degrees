FROM golang:1.14-alpine

LABEL description="Image for 32 degrees server"

RUN mkdir /go/src/32-degrees

WORKDIR /go/src/32-degrees
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["32-degrees"]
