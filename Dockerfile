#-------------------Build Stage (sidecar)-------------------

FROM golang:1.21.6 AS build_base

WORKDIR /go/src/github.com/kubevirt/sidecar-shim

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -o sidecar_shim sidecar_shim.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -o onUndefineDomain ./sidecar

#-------------------Run Stage-------------------
FROM alpine:3.19

RUN apk add --no-cache iproute2 sudo bash libvirt-client

# Install Open vSwitch and Supervisor
RUN apk add --no-cache openvswitch supervisor

# Install OvS
RUN mkdir -p /var/run/openvswitch
RUN ovsdb-tool create /etc/openvswitch/conf.db /usr/share/openvswitch/vswitch.ovsschema
ADD supervisord.conf /etc/supervisord.conf

COPY --from=build_base /go/src/github.com/kubevirt/sidecar-shim/sidecar_shim /sidecar-shim
COPY --from=build_base /go/src/github.com/kubevirt/sidecar-shim/onUndefineDomain /usr/local/bin/onDefineDomain

ENV REDIS_ADDR service-kubedtn-redis.kubedtn.svc.dev.dslab:6379
ENV REDIS_PASSWORD reins5401
ENV REDIS_DB 0

ADD entrypoint.sh /entrypoint.sh

ENTRYPOINT ["bash", "/entrypoint.sh"]