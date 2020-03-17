FROM golang:1.12.0-alpine3.9 AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

ADD . /app

# Pass in prebuilt modules.  If not using them, uncomment next line to have go build them.
# RUN go mod vendor
RUN go build -ldflags "-s -w" -mod=vendor -o users .

FROM alpine:3.11
WORKDIR /app
COPY --from=builder /app/users /app

CMD ["/app/users"]