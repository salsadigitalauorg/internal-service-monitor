FROM golang:1.22 AS builder

ARG VERSION
ARG COMMIT

ADD . $GOPATH/src/github.com/salsadigitalauorg/internal-services-monitor/

WORKDIR $GOPATH/src/github.com/salsadigitalauorg/internal-services-monitor

ENV CGO_ENABLED 0

RUN apt-get install ca-certificates

RUN go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" -o build/internal-services-monitor

FROM scratch

ARG PORT=3000
ARG CONFIG=cfg.yml

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/salsadigitalauorg/internal-services-monitor/build/internal-services-monitor /usr/local/bin/internal-services-monitor

EXPOSE $PORT

CMD [ "internal-services-monitor", "-port", $PORT, "-config", $CONFIG ]
