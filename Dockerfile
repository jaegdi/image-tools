# Multistage Stage 0 as base
FROM golang:alpine as builder
ENV CGO_ENABLED=1
ENV https_proxy=http://webproxy.sf-bk.de:8181/
ENV HTTPS_PROXY=$https_proxy
WORKDIR /usr/local/go/src/image-tool
COPY build ./build/
COPY docs ./docs/
COPY ocrequest ./ocrequest/
COPY *.go .
COPY go.mod .
COPY clusterconfig.json .
RUN apk update -U \
    && apk --no-cache add ca-certificates tzdata libc6-compat libgcc libstdc++ \
    && ls -l build/certs

RUN apk add --no-cache --update git build-base
RUN pwd \
    && ls -l \
    && go mod tidy \
	&& go build

# ----------------------------------------------------------------------------------------------------------------------------------

FROM alpine:latest as runner
LABEL maintainer="Dirk Jäger <dirk.jaeger@schufa.de>"

ENV LANG=en_US.UTF-8
ENV TZ=Europe/Berlin
ENV https_proxy=http://webproxy.sf-bk.de:8181/
ENV HTTPS_PROXY=$https_proxy

RUN apk update -U \
    && apk --no-cache add ca-certificates tzdata libc6-compat libgcc libstdc++
COPY --from=builder /usr/local/go/src/image-tool/build/certs/*.crt /usr/local/share/ca-certificates/
COPY --from=builder /usr/local/go/src/image-tool/image-tool /usr/bin/image-tool
COPY --from=builder /usr/local/go/src/image-tool/clusterconfig.json /usr/bin/

RUN chgrp root /usr/bin/image-tool \
    && chgrp root /usr/bin/clusterconfig.json \
    && chmod a+rx /usr/bin/image-tool \
    && chmod a+rw /usr/bin/clusterconfig.json \
    && update-ca-certificates

CMD ["image-tool", "-socks5=no", "-server", "-statcfg"]
