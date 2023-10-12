FROM golang:1.20-alpine as build

WORKDIR /app

RUN apk add --no-cache tzdata

COPY . ./

RUN go mod tidy

RUN go build -o tax-aggregator-service app/main.go

FROM alpine AS main

WORKDIR /app

RUN apk --no-cache add ca-certificates curl

COPY --from=build /app/tax-aggregator-service /app
COPY --from=build /app/config/config.json /app/config/config.json

ENTRYPOINT [ "/app/tax-aggregator-service" ]