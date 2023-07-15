FROM golang:1-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /app/out/api.go

FROM golang:1-alpine

COPY --from=build /app/out /app

EXPOSE 8080

CMD ["/app/api.go"]
