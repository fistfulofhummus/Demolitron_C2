
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
 - [ ] Install Script that satisfies all the dependencies automatically
 - [ ] Encrypt the tcp connections
 - [ ] Better error handling
 - [ ] Load and execute shellcode
 - [ ] Upload/Download
 - [ ] Persistence
 - [ ] Keylogging
 - [ ] RickRoll

## Notes:
The project is under "active" development. I am just a monkey with a keyboard y all working on this when I have the time and learning as I go.
