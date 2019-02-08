# dnsdist-autoconf

Simple [dnsdist](https://dnsdist.org) config generator made for human.

Prebuilt binaries might be found in [releases](https://github.com/Jamesits/dnsdist-autoconf/releases) or from the CI below.

Integrated Docker image: [Docker Cloud](https://cloud.docker.com/repository/docker/jamesits/dnsdist-autoconf)

[![Build Status](https://dev.azure.com/nekomimiswitch/General/_apis/build/status/dnsdist-autoconf?branchName=master)](https://dev.azure.com/nekomimiswitch/General/_build/latest?definitionId=39?branchName=master) [![](https://images.microbadger.com/badges/image/jamesits/dnsdist-autoconf.svg)](https://microbadger.com/images/jamesits/dnsdist-autoconf "Get your own image badge on microbadger.com")

## Features

* Set different DNS servers for different domains
* Integrated [felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list) rules
* Automatically probe Active Directory 

## Usage

### [I Just Wanna Run](https://www.youtube.com/watch?v=HrWnfx8uRPw)

An example config file is at [examples/autoconf.toml](examples/autoconf.toml).

```shell
# generate the config
dnsdist-autoconf -config autoconf.toml -output dnsdist.conf
# check the config grammar (important, since the author is not very confident)
dnsdist -C dnsdist.conf --check-config
# run it!
dnsdist -C dnsdist.conf
```

### Use Hosted Dnsdist in Docker

The docker image will rerun `dnsdist-autoconf` every night to update dynamic config.

Set `REMOTE_CONFIG` and `autoconf.toml` will be updated too.

#### Option 1: directly run the docker image

0. Put a [dnsdist-autoconf config toml file](examples/autoconf.toml) into `/etc/dnsdist`
1. Run an instance of `jamesits/dnsdist-autoconf:latest`

```shell
docker pull jamesits/dnsdist-autoconf:latest
docker run --rm --name=dnsdist-autoconf_1 -p=53:53/udp -p=53:53/tcp -p=8083:80/tcp -v=/etc/dnsdist:/etc/dnsdist jamesits/dnsdist-autoconf:latest
```

#### Option 2: manage it with systemd

0. Make sure your current OS have good DNS (at least can connect to the Docker registry and let dnsdist-autoconf finish probing services)
1. Put [`dnsdist-autoconf.service`](docker/dnsdist-autoconf.service) in this repo to `/usr/lib/systemd/system`
2. Put a [dnsdist-autoconf config toml file](examples/autoconf.toml) into `/etc/dnsdist`
3. Start (and optionally enable) the `dnsdist-autoconf.service` systemd unit

Example:

```shell
mkdir -p /usr/lib/systemd/system
mkdir -p /etc/dnsdist
wget https://github.com/Jamesits/dnsdist-autoconf/raw/master/docker/dnsdist-autoconf.service -O /usr/lib/systemd/system/dnsdist-autoconf.service
wget https://github.com/Jamesits/dnsdist-autoconf/raw/master/examples/autoconf.toml -O /etc/dnsdist/autoconf.toml
systemctl daemon-reload
systemctl enable --now dnsdist-autoconf.service
```

## Building

Use Go 1.10 or higher.

## Caveats

### dnsdist version compatibility

We only support dnsdist version 1.3 and later. Although there are some cases running dnsdist 1.2 with it, these cases will less likely be supported.

### Disable systemd-resolved

`systemd-resolved` will take up port 53 on Ubuntu 17.04 onwards. To disable it:

0. Make sure your hostname resolves in `/etc/hosts`
1. If your `/etc/resolv.conf` is a symlink, delete it and recreate it, type in `nameserver 8.8.8.8` or any other working DNS server
2. `systemctl disable --now systemd-resolved.service`
3. `systemctl mask systemd-resolved.service`
4. Config your DHCP client to use the appropriate DNS config. For example, if using NetworkManager, add `dns=default` under `[main]` section of `/etc/NetworkManager/NetworkManager.conf`.

### Active Directory

We make a simple assumption that every DC have DNS roles installed, since we can only get LDAP/Kerberos server list from DNS queries, and quering any other config requires much more complex protocols. 

### ulimit (Too many open files)

The generated config might cause dnsdist to use a lot file descriptors.

```
Warning, this configuration can use more than 1220 file descriptors, web server and console connections not included, and the current limit is 1024.
You can increase this value by using LimitNOFILE= in the systemd unit file or ulimit.
```

Quick fix if you are running directly in a shell:

```shell
# you might need root privilege
ulimit -u unlimited
dnsdist -C dnsdist.conf
```

Fix if you are running in systemd:

```shell
mkdir -p /etc/systemd/system/dnsdist.service.d
echo -e "[Service]\nLimitNOFILE=16384\n" > /etc/systemd/system/dnsdist.service.d/ulimit.conf
systemctl daemon-reload
```
