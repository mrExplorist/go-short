FROM golang:alpine as builder
# alpine is a lightweight linux distribution 

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go build -o main .

# second stage of the build 

FROM alpine

RUN adduser -S -D -H -h /app appuser

USER appuser

COPY . /app

COPY --from=builder /build/main /app/

WORKDIR /app

EXPOSE 3001

CMD ["./main"]
