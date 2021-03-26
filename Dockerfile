FROM golang

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /dist

RUN cp /app/main .

EXPOSE 8080

CMD ["/dist/main"]