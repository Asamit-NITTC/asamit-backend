FROM golang:1.20 AS build-env

WORKDIR /api

COPY ./api/go.mod .
COPY ./api/go.sum .

RUN go mod download

COPY /api .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM gcr.io/distroless/base-debian10

COPY --from=build-env /api/main /api/main

EXPOSE 8080

CMD ["/api/main"]

