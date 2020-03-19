FROM golang:1.12.0-alpine3.9 AS builder

RUN apk update && apk add --no-cache git nodejs npm

WORKDIR /app

COPY . /app

# If vendor isn't there run the command.
RUN if [ ! -d vendor ]; \
    then echo "Need to build"; go mod vendor; \
    fi



RUN go build -ldflags "-s -w" -mod=vendor -o web .

WORKDIR /app/webfiles/chorelist

RUN if [ ! -d "dist" ]; then \
# TODO: Look into why node-sass requires a rebuild...
    npm install \
    npm rebuild node-sass \
    npm run build; \
    fi

FROM alpine:3.11
WORKDIR /app
COPY --from=builder /app/web /app
COPY --from=builder /app/webfiles/chorelist/dist/ /app/webfiles

CMD ["/app/web"]