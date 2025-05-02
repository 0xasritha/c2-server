# Current Features
- **Exploit:**
    - working lab env (`Ubuntu 18.04`)
    - working exploit binary, triggered when user pulls evil image and runs `docker exec /bin/sh` (supply chain attack)
- **Implant:**
    - implant binary done w/ setup for both HTTPS + DNS communication          
        - HTTPS fully functional 
    - designed so its very easy to add new tasks + communication protocols 
        - can modify implant configuration at any time 
- **C2 Server**:
    - hooked up to a Postgres DB, so: 
        - can scale multiple instances of the C2 server 
        - can restart, modify, upgrade C2 server at any time w/out loosing data  
    - c2 server done w/ support for keeping track of MULTIPLE implants + their corresponding tasks + results 
        - easily queriable endpoints for operator + implant
        - setup to easily extend and add new communication protocols
            - HTTPS fully functional

- **Infra**: 
    - working Terraform (Infrastructure as Code) file to easily spin up our resources in a GCP cloud env (C2 server + it's firewall rules, Storage bucket to host implant binaries + for file transfer tasks, private container registry to test/demo evil images)

---
# TODO: 
- **Exploit:**
    - write script to pull any Docker image off Dockerhub, inject it w/ exploit binary, create new evil image (*NICE TO HAVE*)
    - write script to install Docker v 18.09.1 to quickly setup lab env
    - write "replace runC script": should redownload runC for running future containers, execute dropper script
    - write dropper script to pull evil image 2 w/ full implant binary, extract it, run it in memory (https://github.com/amenzhinsky/go-memexec) 

- **Implant:**
    - write Makefile to build implant binary w/ `garble` (obfuscation), strip debug info/symbol tables, etc.
    - Kill switch
    - Implement Tasks: 
        - Update Config 
        - File Operations (*NICE TO HAVE*)
            - Need to update C2-server to create folder in GCP bucket for each implant + change GCP bucket permissions to publically writable? (temporarily)
        - Cron Job (Persistence)
        - `Rc.local` modifications? (Persistence) (*NICE TO HAVE*)
        - Delete Docker container/daemon, Docker logging software logs (*NICE TO HAVE*)

- **C2-Server:**
   - add CLI for operators (*NICE TO HAVE*)
   - encryption (need to add corresponding implant code too): 
        - add JWTs for implant authorization, for HTTPS only (working Proof of Concept in `./secure`) 
        - check over/ add in AES GCM for symmetric key encryption (working Proof of Concept in `./secure`)
        - add in symmetric key exchange via/ DiffieHellman https://eminmuhammadi.com/articles/diffiehellman-key-exchange-example-in-golang) 
   - integrate DNS tunneling as second communication protocol option (Proof of Concept in `./dns-tunneling`)
        - tunneling 
        - tsig (Kat)  

- **Infra:**
    - setup proxy (*NICE TO HAVE*)
- **Presentation**
