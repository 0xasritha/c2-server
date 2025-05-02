# DNS Tunneling

- attack adds malicious payloads/commands to trusted DNS
    - some DNS tunneling software needs client software too (on victim machine, distribute through phising, etc.)
        - attacker will comprimise a client, run software on victim's machine w/ a specific domain name; all requests will be sent to local DNS server, then to authoritative DNS server for that domain (controlled by attacker as/ C2 server)
- attack: 
    - 1) attacker runs offensive DNS tunneling software on public server
    - 2) attacker registers public domain w/ authoritative DNS set to the above server ^
        - authoritative DNS server: source of truth (final authority) for a domain's DNS records

---
- NS record (nameserver record): maps a domain name to list of DNS servers (if website is hosted on multiple servers)

1) client encodes data in DNS request (ex. prepend piece of data in domain of request) 
2) DNS request goes out to DNS server 
3) DNS server finds A register of domain w/ IP address of your server 
4) request for that domain is forwrard to the server 
5) server processes whatever mypieceofdata was supposed to do (assuming it was HTTP request) 
6) server replies back over DNS

---

- 4 DNS servers involved in loading a webpage: 
  - 1) DNS recursor: server that receives queries from clients through apps (ex. web browser); makes additional requests to satisfy client's DNS query
  - 2) root nameserver: first step in translating host names -> IP addresses; typically serves as reference to other more specific locations 
   - 3) TLD (top level domain) nameserver: next step in search for specific IP address; hosts last portion of hostname (ex. `com` in `example.com`)
   - 4) authoritative nameserver: final nameserver; if it has access to requested record, will return IP addresses of requested hostname back to DNS recursor 

- recursive DNS resolver: computer that responds to recursive request from client, tracks down DNS record by making series of requests until it reaches authoritative DNS server for requested record
  - can slow lookup process via caching 


- authoritative DNS server: server that holds the DNS resource record (final source of truth, no more further queries) (ex. Squarespace, Namecheap)
  - when query is subdomain, additional nameserver can be added to sequence after authoritative nameserver (to store subdomains CNAME record)


- steps in DNS lookup (will skip some steps if cached): 
  - 1) user types domain into browser; DNS recursive resolver receives it 
  - 2) resolve queries DNS root nameserver (.). The DNS root nameservers contain records for all top level domain name servers
   - 3) root server responds to resolver w/ address of TLD DNS server (ex. `.com`)
  - 4) resolver makes request to `.com` TLD
  - 5) TLD server responds w/ IP address of domain's nameserver (`example.com`)
  - 6) recursive resolver sends query to the domain's nameserver
  - 7) IP address for `example.com` is returned to resolver from nameserver 
  - 8) DNS resolver responds to browser w/ IP of the domain requested initially, now browser can make request to web page


- types of DNS queries:
  - recursive query: DNS client requires that DNS server (so the DNS recursive resolver) responds to client w/ either requested resource or error message 
  - iterative query: DNS client allows a DNS server to return best possible answer - so if queried DNS server doesn't have match for query name, it returns referral to DNS server authoratitve for a lower level of the domain namespace. client queries there, process continues until found or error 
  - non-recursive query: happens when DNS resolver client queries DNS server for a record that it has access to bc its authoritative for the record or the record exists within cache


- NS record (Name Server): specifies which DNS servers are authoritative for a domain
  - each NS record points to hostname (not IP) of authoritative DNS server 
  - can have multiple NS records for redundancy (primary + secondary nameservers)
  - NS records must exist at both parent zone + within child zone 
  - critical during delegation (ex. when parent zone for `.com` tells clients which servers manage `example.com`)
  - `example.com.   IN   NS    ns1.example1dns.com`


 
---
# DNS Tunneling 

- how it exchanges information: 
  - encodes data inside subdomain part of DNS query as text (Base32, Base64, etc.)
  - implant encodes data (ex. commands, stolen data) into domain (ex. sends DNS query like `<encoded_data>.attacker-controlled-domain.com`)


- most common: 
  - TXT records (bc they allow arbitrary text data)
  - A records: implant encodes data into subdomains, response gives dummy IP)



