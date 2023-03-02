FROM golang:alpine AS base

WORKDIR /app

COPY *go* /app
COPY pkg/ /app/pkg

RUN go build -o app

RUN apk add --no-cache --update ca-certificates tzdata

RUN adduser -D -H -h / -s /sbin/nologin app

FROM scratch AS packed

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /usr/share/zoneinfo/ /usr/share/zoneinfo/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /app/app /usr/local/bin/app

USER app

ENV TZ="Europe/Budapest"

ENTRYPOINT ["app"]
