FROM golang:1.19-alpine

WORKDIR $GOPATH/src/stark-tag-api/cmd/app

COPY . .

RUN go install -v ./...

EXPOSE 8080 50051

CMD ["server"]