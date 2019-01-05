# dnsdist-autoconf

Simple [dnsdist](https://dnsdist.org) config generator made for human.

[![Build Status](https://dev.azure.com/nekomimiswitch/General/_apis/build/status/dnsdist-autoconf?branchName=master)](https://dev.azure.com/nekomimiswitch/General/_build/latest?definitionId=39?branchName=master)

## Features

* Set different DNS servers for different domains
* Integrated [felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list) rules
* Automatically probe Active Directory 

## Usage

An example config file is at [examples/config.toml](examples/config.toml).

```shell
dnsdist-autoconf -config config.toml -output dnsdist.conf
dnsdist -C dnsdist.conf
```

## Caveats

### ulimit

The generated config might cause dnsdist to use a lot file descriptors.

```
Warning, this configuration can use more than 1220 file descriptors, web server and console connections not included, and the current limit is 1024.
You can increase this value by using LimitNOFILE= in the systemd unit file or ulimit.
```

Quick fix if you are running directly in a shell:

```shell
sudo su
# ulimit -u unlimited
# dnsdist -C dnsdist.conf
```

Fix if you are running in systemd:

```shell
mkdir -p /etc/systemd/system/dnsdist.service.d
echo -e "[Service]\nLimitNOFILE=16384\n" > /etc/systemd/system/dnsdist.service.d/ulimit.conf
systemctl daemon-reload
```