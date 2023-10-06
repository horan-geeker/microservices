FROM golang:1.20 AS build

ARG PROJECT_ROOT=/src
ENV PROJECT_ROOT ${PROJECT_ROOT}
COPY . $PROJECT_ROOT/
# compile and dynamic linking
RUN set -ex && \
cd $PROJECT_ROOT && \
go mod tidy && \
GO111MODULE=on CGO_ENABLED=1 GOOS=linux go build -o /execute && \
([ -x "/execute" ] || exit -1)

FROM alpine:latest AS final

RUN apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++

RUN cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    mkdir -p /binary

COPY --from=build /execute /binary/execute

EXPOSE 80

WORKDIR /binary

ENTRYPOINT ["./execute"]