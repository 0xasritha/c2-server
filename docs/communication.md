- #TODO: have chat reword this to be better
- #TODO: put all DNS files through chat, fix any errors
# C2Server <-> Implant Communcation 

- two main functionalities: 
    - beacon: responsible for both posting results, getting commands 

## HTTPS 

- new: `https://www.cloud-docker.net/implant/new` -> change back to register
// TODO: add JSON Payloads here 
- beacon: `https://www.cloud-docker.net/implant/beacon`
// TODO: add JSON Payloads here 

## DNS
# TODO: add in improvesments that DNS can be extended to support multiple task types
- everywhere, using "command" instead of "cmd"

While creating DNS code, we kept in mind the following limits: 
	DNS limits:
		- full QNAME (∑ (1 + len(labelᵢ))  +  1  ≤  255) -> this one not using 
			- each label is prefixed by a single length-octet
			- terminated by zero length-octet (0x00)
		- each label must be maximum 63 bytes


- DNS only supports ShellCommand Tasks
- all Implants must initially reigster as HTTPS (we are assuming outbound HTTPS traffic from the implant is allowed as the beachead?/stage 1 dropper was able to download the implant binary from a GCP storage bucket). In order to be more stealthy, the operator can later change the implant to use any protocol, and switch between communication protocols 

- all results + command payloads exchanged between the servers will be base32 encoded

- beacon: 
    - post results: `<implantID>.<taskID>.<chunkCount>.<chunkIndex>.<chunkData>.r.cloud-docker.net`
    - get commands: `<implantID>.c.cloud-docker.net`
        
