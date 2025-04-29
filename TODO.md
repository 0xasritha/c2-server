# Capstone Project

Due Date: May 7 

## Rubric 
- **Exploit**: 
    - CVE-2019-5763
- **Implant for execution + data exfiltration**: 
- **C2 Server for implant tasking and data reception**: 
- **Encryption/Obfuscation for implant, C2 communication, data exfiltration**
- **Clever Injection Method**

## MVP 

By Friday: 
- working implant 
- working C2 server w/ DNS tunneling + hooked up to postgres DB 
- working base Dockerfile w/ exploit code hidden + working script to generate more poisonous dockerfiles


---
- Listener API hooked to Postgres DB
    - basic authentication for implant + operator 
    - sample requests file w/ curl requests operator can subtitute JSON payloads 
- Working implant that does the 4 Task Types in blog 
- Working base Dockerfile w/ exploit code hidden for supply chain attack  
- Script to poison images + push to Dockerhb 
- DNS Tunneling to mask API requests from implant to listener 
- XOR obfuscation for implant binary 
- Persistence tasks? Create cronjob, ssh key, etc.? (look up linode red team persistence)
    - can hijack some docker command?
    - have task to kill auditd + logging software
- Run implant in memory, not in disk (small C file thing)


---
run implant in memory, not in disk - small C file thing) https://www.linode.com/docs/guides/linux-red-team-persistence-techniques/ - have tasks for this
      Use MITRE Attack Framework in presentation for each step
      Proxy + TF to spin up c2 infra 
     Persist through hijacking some docker command 
     kill auditd + logging software (have this as task)?
 Encryption obfuscation techniques (look at MITRE)
minification
https://github.com/amenzhinsky/go-memexec  -> see chat 564 project

---

## C2 Server ("listener")

- Authentication/Authorization for implants + operators (they should each be able to access only their endpoints)
- Operator CLI 
