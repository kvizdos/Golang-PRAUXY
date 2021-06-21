FROM golang:alpine
WORKDIR /src
COPY . .
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get github.com/pquerna/otp/totp
RUN go get github.com/go-redis/redis/v8
RUN go mod download golang.org/x/crypto
RUN go mod download golang.org/x/sync
RuN apk add --update alpine-sdk
CMD ["go", "test", "-v", "e2e_test.go"]