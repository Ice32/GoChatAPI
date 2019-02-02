FROM golang:alpine

WORKDIR /go/src/bitbucket.org/KenanSelimovic/GoChatServer
COPY . .
RUN apk add --no-cache bash
RUN chmod +x wait-for-it.sh

RUN apk add --no-cache git mercurial
RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 3183
CMD ["GoChatServer"]