- client-server architecutre
    - client: "operator"/"players" connect to same server 
    - `~/.silver/go`: version of the Go toolchain dedicated to silver, used to compile implants/C2 agents
    - uses `garble` to obfuscate Go binaries https://github.com/burrowers/garble
    - operators authenticate to server w/ mutual TLS (server issues client certificates)
    - can use silver authentication process for clients

- 4 C2 protocols: 
    - mTLS (mutual TLS) - messaging layer security 
    - Wireguard
    - for more restricted env: 
        - HTTP
        - DNS
    - use DNS server for DNS names for C2 endpoints instead of IP addresses
    - C2 server serves implant binaries on port 80

