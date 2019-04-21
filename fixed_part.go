package main

// const strings

const globalConfigPrependString = `
-- global static config start --

-- fix up possibly badly truncated answers from pdns 2.9.22
truncateTC(true)

-- the default is leastOutstanding, which part of ignores the "order" of the server
setServerPolicy(firstAvailable)
setServFailWhenNoServer(true)

-- predefined IP masks

-- All private IPs
PrivateIPs = newNMG()
PrivateIPs:addMask("0.0.0.0/8")
PrivateIPs:addMask("10.0.0.0/8")
PrivateIPs:addMask("100.64.0.0/10")
PrivateIPs:addMask("127.0.0.0/8")
PrivateIPs:addMask("169.254.0.0/16")
PrivateIPs:addMask("172.16.0.0/12")
PrivateIPs:addMask("192.0.2.0/24")
PrivateIPs:addMask("192.88.99.0/24")
PrivateIPs:addMask("192.168.0.0/16")
PrivateIPs:addMask("198.18.0.0/15")
PrivateIPs:addMask("198.51.100.0/24")
PrivateIPs:addMask("203.0.113.0/24")
PrivateIPs:addMask("224.0.0.0/4")
PrivateIPs:addMask("240.0.0.0/4")
PrivateIPs:addMask("::/8")
PrivateIPs:addMask("0100::/64")
PrivateIPs:addMask("2001:2::/48")
PrivateIPs:addMask("2001:10::/28")
PrivateIPs:addMask("2001:db8::/32")
PrivateIPs:addMask("2002::/16")
PrivateIPs:addMask("3ffe::/16")
PrivateIPs:addMask("fc00::/7")
PrivateIPs:addMask("fe80::/10")
PrivateIPs:addMask("fec0::/10")
PrivateIPs:addMask("ff00::/8")

-- default block rules

-- refuse all queries not having exactly one question
addAction(NotRule(RecordsCountRule(DNSSection.Question, 1, 1)), RCodeAction(DNSRCode.REFUSED))

-- global static config end --

`
