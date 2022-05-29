# =========================================
# build stage
# =========================================
FROM golang:1.18-alpine as build-stage

WORKDIR /opt/vpn_client

ARG TZ
RUN apk --no-cache add tzdata && \
  cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
  echo "${TZ}" > /etc/timezone && \
  apk del tzdata

COPY go.mod go.sum ./
RUN apk --no-cache add openvpn make

COPY ./ ./
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -gcflags="all=-N -l" -o client

CMD ["/bin/sh"]

# =========================================
# production stage
# =========================================
FROM alpine as production-stage

WORKDIR /opt/vpn_client

ARG TZ
RUN apk --no-cache add tzdata && \
  cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
  echo "${TZ}" > /etc/timezone && \
  apk del tzdata

RUN apk --no-cache add openvpn
COPY --from=build-stage /opt/vpn_client/client /usr/bin/client

ENTRYPOINT ["client"]
