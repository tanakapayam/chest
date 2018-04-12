# building...
FROM golang:alpine as builder

ENV APP="chest" \
    CGO_ENABLED=0 \
    TERM="xterm"

WORKDIR /go/src/github.com/tanakapayam/"$APP"
COPY . .

RUN apk --update --no-cache add \
        bash \
        git \
        make \
        ncurses \
    && make install

# minimal image
FROM alpine:3.8

ENV APP="chest" \
    EJSON_KEYDIR="/ejson" \
    TERM="xterm"

RUN apk --update --no-cache add \
        ncurses \
    && addgroup -S "$APP" \
    && adduser -D -S -G "$APP" -H -h "/$APP" "$APP" \
    && rm -rf /var/cache/apk/*

COPY --from=builder \
    /go/bin/"$APP" \
    /usr/bin/"$APP"

USER "$APP"
WORKDIR "/$APP"
VOLUME ["/chest", "/ejson"]
ENTRYPOINT ["chest"]
