FROM golang as build

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /dist

COPY --from=build /app/main ./main

EXPOSE 8080

ENTRYPOINT [ "/dist/main" ]