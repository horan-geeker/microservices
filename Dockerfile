FROM golang:alpine AS build
RUN apk update && apk add ca-certificates && apk add tzdata
WORKDIR /app
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM scratch AS final
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/app /
COPY --from=build /app/.env /

ENV TZ Asia/Shanghai
ENTRYPOINT ["/app"]