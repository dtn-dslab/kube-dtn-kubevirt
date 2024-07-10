#-------------------Build Stage (sidecar)-------------------

FROM golang:1.22.3 AS build_base

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /go/src/github.com/kubevirt/sidecar-shim

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -o sidecar_shim sidecar_shim.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -o onUndefineDomain ./sidecar

#-------------------Run Stage-------------------
FROM alpine:3.19

RUN apk add --no-cache iproute2 sudo bash libvirt-client tcpdump

# Install Open vSwitch and Supervisor
RUN apk add --no-cache openvswitch supervisor

# Install OvS
RUN mkdir -p /var/run/openvswitch
RUN ovsdb-tool create /etc/openvswitch/conf.db /usr/share/openvswitch/vswitch.ovsschema
ADD supervisord.conf /etc/supervisord.conf

COPY --from=build_base /go/src/github.com/kubevirt/sidecar-shim/sidecar_shim /sidecar-shim
COPY --from=build_base /go/src/github.com/kubevirt/sidecar-shim/onUndefineDomain /usr/local/bin/onDefineDomain

ADD entrypoint.sh /entrypoint.sh

ENTRYPOINT ["bash", "/entrypoint.sh"]