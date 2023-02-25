FROM golang:1.19.4-alpine3.17 as build
ARG OS
ARG ARCH
WORKDIR /build
COPY . .
RUN go mod download && go build -o cleanmyarr

FROM alpine:3.17
ARG VERSION
ARG user=cma
ARG group=cma
ARG uid=1000
ARG gid=1000
USER root
WORKDIR /app
COPY --from=build /build/cleanmyarr /app/cleanmyarr
COPY container-entrypoint.sh /app/container-entrypoint.sh
COPY sample-config.yaml /app/sample-config.yaml
RUN apk update && apk --no-cache add bash vim && addgroup -g ${gid} ${group} && adduser -h /app -u ${uid} -G ${group} -s /bin/bash -D ${user}
RUN mkdir /config
RUN chown cma:cma /app/cleanmyarr && chmod +x /app/cleanmyarr && \
    chown cma:cma /app/container-entrypoint.sh && chmod +x /app/container-entrypoint.sh && \
    chown cma:cma /config && chmod u+rw /config
USER cma
VOLUME [ "/config" ]
ENTRYPOINT [ "/app/container-entrypoint.sh"]