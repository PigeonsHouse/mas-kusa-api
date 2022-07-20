FROM golang:1.18.3-alpine

RUN mkdir -p /go/src/OnakaAPI/
COPY . /go/src/OnakaAPI/
WORKDIR /go/src/OnakaAPI/

RUN go install

RUN go get -u github.com/cosmtrek/air && \
    go build -o /go/bin/air github.com/cosmtrek/air

CMD ["air"]
