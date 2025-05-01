- https://eminmuhammadi.com/articles/diffiehellman-key-exchange-example-in-golang -> use this for the symmetric key exchange

By Friday: 
- working implant
    - rn the implant only works if someone triggers exec /bin/bash on the container -> https://github.com/Frichetten/CVE-2019-5736-PoC/tree/master this works. (ideally would fix this so its upon running container), then write script to add inject this any container, test w nginx. create a script to also install older version of Docker so can easily test this even after runC is replaced
    - instead of /bin/bash like the POC, it should do the wget to pull the implant in
- working base Dockerfile w/ exploit code hidden + working script to generate more poisonous dockerfiles

So: 
    - get the older version of docker installed + install script to show Kat 
    - create a test evil nginx + private container registry
        - have terraform to spin this up (container registry, storage bucket for implant, ec2 instance w/ domain pointed to it)
    - implant should work and run a script that makes runc look legit, downloads implant and runs it, so that more containers can still be ran 
    - working implant + c2 w/ dns 
    - run implant in memory, not in disk - small C file thing) 
    - fix the symmetric key exchange to use RSA symmetric key exchange (+ also fix the shitty handcoded implementation of the symmetric key encoding/decoding) 

---
## Clean up/ Don't Forget 
- have Makefile to build go binary w/ go-obfuscation + go minificaiton  
- standardize error handling 
- Operator CLI
- Kill Switch
- https://github.com/amenzhinsky/go-memexec  -> see chat 564 project
---

# Insane Nice To Haves
- concurrent execution
- self-mutating malware 
- inject into process? 
- Custom Obfuscation/ Other obfuscation methods? (look at MITRE) 
- persist through hijacking a docker command? 
