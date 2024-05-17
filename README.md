
# Demolitron_C2
A C2 inspired by sliver and metasploit. Post exploitation only.
  

## Requirements
Before building the implant don't forget to edit the code and insert your TCP listner address.

 1. Go
 2. chmod +x on all the scripts in the scripts directory

## Usage
From within the directory containing main.go

    go get <any_and_all_import>
    go mod tidy
    go build
    sudo ./server
If you are having trouble with the Bushido implant, feel free to use generateDebug with the same syntax as the generate command from withing the Demolitron console to create a debug version of the implant.


## TO-DO

 - [x]  Add a menu a la meterpreter
 - [ ] Install Script that satisfies most of the dependencies automatically
 - [ ] Encrypt the tcp connections
 - [x] Better error handling
 - [x] Load and execute shellcode
 - [x] Inject Shellcode via Process Hollowing
 - [ ] Upload/Download
 - [ ] Persistence
 - [ ] Keylogging
 - [ ] RickRoll

## Notes:
The project is under "active" development. I am just a monkey with a keyboard. I am working on this when I have the time and learning as I go.
There is also a theoretical hardlimit to how many agents you can have deployed since it relys on IDs generated from 1 to 9000. This can be solved by trying to enumerate the host and get the host id as well and have it stored in the session struct.
