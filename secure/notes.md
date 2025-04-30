# Layers of encryption/security: 
- TLS 
- JWT 

# JWT
- https://www.youtube.com/watch?v=7Q17ubqLfaM&ab_channel=WebDevSimplified 
- used to replace cookies + sessions (w/ this method sessions are managed server-side, not good for scaling)
- JWT - used to verify authorization per user 
- crypto signature: signature on message/token (JWT), verifies data hasn't been altered

- HS256: 
  - H: HMAC (Hash-based Message Authentication Code)
  - S256: Secure Hash Algo 256-bit 
  - HS256 is algo that uses shared secret key + SHA-256 to generate a crypto signatures
  - verifies that: 
    - the token has not been tampered with (integrity)
    - sender knows the secret key (authenticity)

- only parties that know shared secret can: 
  - generate valid tokens 
  - verify tokens haven't been tampered w/

- format: header.payload.signature
  - header is base64 encoded
  - payload: data - can put anything  
    - common fields: 
      - `sub`: ID of user/thing
      - `iat`: issued at - when token was created; good for knowing when to expire 
      - `exp`: data that token expires at
  - (verify) signature: hashed content of (base64Encode(header) + base64Encode(payload))

- server checks if the header + payload matches the signature part -> if match is valid and JWT has not been tampered with; hence can store user information in client; client cannot view or tamper w/ any of the information unless they have secret key 

## AES-GCM encryption

- block cipher: type of symmetric-key encryption algo that encrypts data in fixed-size blocks (64, 128 bits); uses same key for encryption + decryption 

- AES-GCM: symmetric key encryption algo; needs key to encrypt/decrypt data 
  - AES: can only encrypt/decrypt 16 byte long data (its block size) by itself
    - w/ block cipher mode (ex. GCM block cipher mode): can encrypt/decrypt arbitrary sized data; also adds message authentication (integrity)

- https://dev.to/breda/secret-key-encryption-with-go-using-aes-316d



----
Here’s a high-level design for standing up an HTTPS-based C2 listener that can track multiple implants, authenticate them, encrypt payloads end-to-end, and serve everything over TLS.

---

## 1. Implant Identification & Session Tracking

1. **Unique Implant ID**  
   On first check-in, each implant generates a UUID (v4) and a public key (for asymmetric payload encryption; see below). It POSTs this to `/register`.
2. **Server-side Session Store**  
   Keep a map of `implantID → { last_seen, public_key, IP, status, nonce_counter }`. This can live in an in-memory DB (Redis) or lightweight embedded DB (BoltDB).
3. **Heartbeat / Polling**  
   Implants poll `/task` every N seconds. Each request includes:
   ```http
   GET /task
   Cookie: session=<JWT>
   X-Implant-ID: <uuid>
   X-Nonce: <incrementing counter>
   ```
   Server updates `last_seen` on each valid poll.

---

## 2. Authentication & Session Management

### a) TLS Mutual Authentication (mTLS)
- **Client Certs**: Issue each implant its own client certificate signed by your CA.
- **Server Config**: Require `ClientAuth: tls.RequireAndVerifyClientCert` in Go’s `tls.Config`.
- **Advantage**: No cookies/JWTs; TLS handles auth.

### b) JWT Sessions (Fallback)
- On `/register`, server returns a signed JWT (HMAC-SHA256 or RSA).
- Implant stores JWT in memory and sends it in the `Authorization: Bearer <token>` header (or as an HTTP‐only Secure cookie).
- Server verifies signature on each request and checks `implantID` claim against the session store.

---

## 3. Payload Encryption

Even though TLS encrypts in transit, you can add an extra layer:

1. **Asymmetric Key Exchange**
  - Implant generates an RSA or EC keypair.
  - Sends its public key in the `/register` request (encrypted under the server’s public key).
2. **Per-Message Symmetric Keys**
  - Server uses the implant’s public key to encrypt a random AES-GCM key per callback.
  - Payloads are then encrypted with AES-GCM under that key.
3. **Replay Protection**
  - Include an increasing nonce or timestamp inside the GCM‐AAD.
  - Server increments a counter in the session store; implant rejects any message with a stale nonce.

---

## 4. Setting Up HTTPS on the Listener

### a) Obtain a Certificate
- **Let’s Encrypt (ACME)**:
  ```bash
  certbot certonly --standalone -d c2.example.com
  ```
- **Self-Signed (for lab use)**:
  ```bash
  openssl req -x509 -newkey rsa:4096 -days 365 \
    -nodes -keyout server.key -out server.crt \
    -subj "/CN=c2.example.com"
  ```

### b) Server Configuration Examples

#### Go HTTP Server with mTLS
```go
cert, _ := tls.LoadX509KeyPair("server.crt", "server.key")
caCert, _ := ioutil.ReadFile("ca.crt")
caPool := x509.NewCertPool()
caPool.AppendCertsFromPEM(caCert)

tlsCfg := &tls.Config{
  Certificates: []tls.Certificate{cert},
  ClientCAs:    caPool,
  ClientAuth:   tls.RequireAndVerifyClientCert,
}

srv := &http.Server{
  Addr:      ":443",
  TLSConfig: tlsCfg,
  Handler:   myC2Handler(), // your mux/router
}

log.Fatal(srv.ListenAndServeTLS("", ""))
```

#### Nginx Reverse-Proxy (with Let’s Encrypt)
```nginx
server {
  listen 80;
  server_name c2.example.com;
  location /.well-known/acme-challenge/ { root /var/www/certbot; }
  location / { return 301 https://$host$request_uri; }
}

server {
  listen 443 ssl http2;
  server_name c2.example.com;

  ssl_certificate     /etc/letsencrypt/live/c2.example.com/fullchain.pem;
  ssl_certificate_key /etc/letsencrypt/live/c2.example.com/privkey.pem;

  # Optional mTLS
  ssl_client_certificate /etc/ssl/ca.crt;
  ssl_verify_client on;

  location / {
    proxy_pass       http://127.0.0.1:8080;
    proxy_set_header X-Forwarded-For $remote_addr;
    proxy_set_header Host            $host;
    proxy_set_header X-Implant-ID    $http_x_implant_id;
    proxy_set_header Authorization   $http_authorization;
  }
}
```

---

## 5. Putting It All Together

1. **Registration**
  - Implant POSTs `{ uuid, pubkey }` to `/register`
  - Server stores session and returns JWT (or issues client cert out-of-band)
2. **Tasking Loop**
  - Implant polls `GET /task` with JWT or mTLS cert
  - Server validates, updates `last_seen`, returns encrypted task blob
3. **Result Upload**
  - Implant POSTs to `/result` with encrypted data and nonce
  - Server decrypts with AES key, verifies nonce, stores result

---

### Key Takeaways

- **mTLS** is your strongest auth: each implant proves identity at the TLS layer.
- **JWT sessions** are simpler to prototype but live in your app layer.
- **AES-GCM per-implant keys** give you payload confidentiality and integrity beyond TLS.
- **A robust session store** ties it all together: implant IDs, pubs keys, nonces, last seen, statuses.

With this foundation you can scale to dozens—or hundreds—of implants, all securely talking to a single HTTPS listener.