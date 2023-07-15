FROM golang:1-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN GOOS=linux GOARCH=amd64 go build -o main

FROM public.ecr.aws/lambda/go:1

COPY --from=build /app/main ${LAMBDA_TASK_ROOT}

CMD [ "main" ]
