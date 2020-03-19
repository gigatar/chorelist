FROM golang:1.12.0-alpine3.9 AS builder

RUN apk update && apk add --no-cache git nodejs npm

WORKDIR /app

ADD . /app

# Pass in prebuilt modules.  If not using them, uncomment next line to have go build them.
# RUN go mod vendor
RUN go build -ldflags "-s -w" -mod=vendor -o web .

# WORKDIR /app/webfiles/chorelist

# RUN npm install
# RUN npm run build

FROM alpine:3.11
WORKDIR /app
COPY --from=builder /app/web /app
# COPY --from=builder /app/webfiles/chorelist/dist/ /app/webfiles

CMD ["/app/web"]