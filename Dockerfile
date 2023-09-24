FROM golang:1.20 AS build-env

WORKDIR /api

COPY ./api/go.mod .
COPY ./api/go.sum .

RUN go mod download

COPY /api .

RUN GOOS=linux GOARCH=amd64 go build -mod=readonly -v -o main .

FROM gcr.io/distroless/base-debian10

COPY --from=build-env /api/main /api/main

EXPOSE 8080

CMD ["./main"]

