# build stage
FROM ubuntu:19.04 as builder

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y \
	&& apt-get install -y git-core golang-go

WORKDIR /root/dnsdist-autoconf
COPY *.go /root/dnsdist-autoconf/

ENV GOPATH=/tmp/go
ENV GOBIN=/tmp/go/bin
RUN /usr/bin/go get ./... \
    && /usr/bin/go build -ldflags="-s -w" -o dnsdist-autoconf

# production stage
FROM debian:10-slim
LABEL maintainer="docker@public.swineson.me"

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update -y \
    && apt-get install -y --no-install-recommends curl ca-certificates supervisor cron

# add PowerDNS repo
COPY docker/pdns.list.buster /etc/apt/sources.list.d/pdns.list
COPY docker/dnsdist.perference /etc/apt/preferences.d/dnsdist
RUN curl https://repo.powerdns.com/FD380FBB-pub.asc -o /etc/apt/trusted.gpg.d/pdns.asc

RUN apt-get update -y \
	&& apt-get install -y --no-install-recommends dnsdist \
	&& apt-get clean -y \
	&& rm -rf /var/lib/apt/lists/*

# copy executables
RUN mkdir -p /usr/local/bin
COPY --from=builder /root/dnsdist-autoconf/dnsdist-autoconf /usr/local/bin/
COPY docker/docker-entrypoint.sh /usr/local/bin/
COPY docker/reload.sh /usr/local/bin/
COPY docker/update-remote-config.sh /usr/local/bin/
# for Windows filesystem compatibility, set executable flag
RUN chmod +x /usr/local/bin/*

# setup crontab
COPY docker/crontab.txt /tmp/
RUN crontab /tmp/crontab.txt \
	&& rm /tmp/crontab.txt \
	&& chmod 600 /etc/crontab

# setup supervisor
COPY docker/supervisord.conf /etc/supervisor/

# setup dnsdist-autoconf default config file
COPY examples/autoconf.toml /etc/dnsdist/autoconf.toml

EXPOSE 53/udp 53/tcp 80/tcp
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["supervisord", "-c", "/etc/supervisor/supervisord.conf"]
