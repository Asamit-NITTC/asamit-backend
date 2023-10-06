FROM golang:1.20 AS build-env

WORKDIR /api

COPY ./api/go.mod .
COPY ./api/go.sum .

RUN go mod download

COPY /api .
RUN GOOS=linux GOARCH=amd64 go build -mod=readonly -v -o api .


EXPOSE 8080

CMD ["./api"]

