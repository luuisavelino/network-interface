FROM golang:alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM scratch

COPY --from=builder /app/main /main

EXPOSE 8080

CMD ["/main"]
