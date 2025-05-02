package main

/*
- before key exchange:
	- both parties agree on parameters:
		- modulus (p)
		- base (g)
	- each party generates its own private key
	- each party generates their own public key by computing PubK_A = (g ^ PrivK_A) mod p
	- each party exchanges their public keys
	- each party can compute the shared key using the other party's public key
	- can use shared key as secret symmetric encryption key
https://eminmuhammadi.com/articles/diffiehellman-key-exchange-example-in-golang
*/
