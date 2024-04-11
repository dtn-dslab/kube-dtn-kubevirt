#-------------------Build Stage (sidecar)-------------------

FROM golang:1.21.6 AS build_base

WORKDIR /kubedtn

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubedtnhack hack/hack.go

#-------------------Run Stage-------------------

FROM alpine:3.19

WORKDIR /kubedtn

COPY --from=build_base /kubedtn/kubedtnhack .
COPY ./hack/entrypoint.sh .

RUN chmod +x ./entrypoint.sh
RUN chmod +x ./kubedtnhack

ENTRYPOINT ./entrypoint.sh