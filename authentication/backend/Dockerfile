FROM golang:alpine
WORKDIR /src
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o ./out/prauxy-authentication .

EXPOSE 8080
CMD ["./out/prauxy-authentication"]