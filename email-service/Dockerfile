FROM golang:1.12.0-alpine3.9 AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . /app

# If vendor isn't there run the command.
RUN if [ ! -d vendor ]; \
    then echo "Need to build"; go mod vendor; \
    fi

RUN go build -ldflags "-s -w" -mod=vendor -o email .

FROM alpine:3.11
WORKDIR /app
COPY --from=builder /app/email /app

CMD ["/app/email"]