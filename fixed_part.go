package main

// const strings

const globalConfigPrependString = `
-- global static config start --

-- fix up possibly badly truncated answers from pdns 2.9.22
truncateTC(true)

-- switch the server balancing policy to round robin,
-- the default being least outstanding queries
-- setServerPolicy(roundrobin)

-- predefined IP masks

-- All private IPs
PrivateIPs = newNMG()
PrivateIPs:addMask("10.0.0.0/8")
PrivateIPs:addMask("172.16.0.0/12")
PrivateIPs:addMask("192.168.0.0/16")
PrivateIPs:addMask("100.64.0.0/10")
PrivateIPs:addMask("169.254.0.0/16")
PrivateIPs:addMask("fc00::/7")
PrivateIPs:addMask("fe80::/10")

-- default block rules

-- refuse all queries not having exactly one question
addAction(NotRule(RecordsCountRule(DNSSection.Question, 1, 1)), RCodeAction(dnsdist.REFUSED))

-- global static config end --

`
