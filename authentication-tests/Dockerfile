FROM golang:alpine
WORKDIR /src
COPY . .
RuN apk add --update alpine-sdk
CMD ["go", "test", "-v", "e2e_test.go"]