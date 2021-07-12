FROM node:lts-alpine AS ui-builder
WORKDIR /node
COPY . .
RUN npm install
RUN echo $(ls)
RUN npm run build
RUN echo $(pwd)
RUN echo $(ls)

FROM golang:alpine AS server-builder
WORKDIR /go/delivery
COPY ./server .
RUN echo "$GOPATH"
RUN ls
RUN pwd
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o ./out/prauxy-authentication-web .

FROM alpine:latest AS host
WORKDIR /app/
COPY --from=server-builder /go/delivery/out/prauxy-authentication-web .
COPY --from=ui-builder /node/dist/ ./dist
RUN echo $(ls)
EXPOSE 8080
CMD ["./prauxy-authentication-web"]