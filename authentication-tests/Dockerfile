FROM golang:alpine
WORKDIR /src
COPY . .
RUN go get go.mongodb.org/mongo-driver/mongo

RuN apk add --update alpine-sdk
CMD ["go", "test", "-v", "e2e_test.go"]