# Builder
FROM golang:1.12.8-alpine3.10 as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app

COPY . .

RUN make engine

# Distribution#
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --no-cache --update add ca-certificates tzdata && \
    mkdir /app

WORKDIR /app 

EXPOSE 9090

COPY --from=builder /app/engine /app

CMD /app/engine