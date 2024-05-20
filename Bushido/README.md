
# Bushido
A implant designed to work with Demolitron_C2. It is based on the splinter implant.  

## Usage/Examples
From within the Demolitron Console run generate --ip <ip4Listener> -p <port4Listener>
The server will look for a scripts directory. enter it and run a script to build the implant in the ../Bushido directory.
If one of these dirs is not present this won't work !!!
So far the implant is only windows x64 compatible. Will probably write a Linux one later down the line.

## TO-DO
 - [x] cd
 - [x] ls
 - [x] pwd
 - [x] hostinfo
 - [x] Load and execute shellcode (Not perfect but works)
 - [ ] Replace fmt with an actual logging library
 - [ ] Upload/Download
 - [ ] Persistence via service creation
 - [ ] Keylogging
 - [ ] RickRoll

## Notes:
I will not be responsible for any illegal activity conducted with this code.
The implant uses raw TCP. OPSEC is a non-concern in this project. IDS/IPS systems should sniff it out.
