
# Demolitron_C2
A C2 inspired by sliver and metasploit. Post exploitation only.

![image](https://github.com/user-attachments/assets/279f25bf-63d1-4db5-a734-7cea95a435bd)

## Requirements
Before building the implant don't forget to edit the code and insert your TCP listner address.

 1. Go
 2. chmod +x on all the scripts in the scripts directory

## Usage
From within the directory containing main.go

    go get <any_and_all_import>
    go mod tidy
    go build
    sudo ./Demolitron
If you are having trouble with the Bushido implant, feel free to use generateDebug with the same syntax as the generate command from withing the Demolitron console to create a debug version of the implant.
If you want to use push notifications, you could link this to the ntfy service by supplying the command line flag that is you ntfy: ./server myChannelURI
This will have the server notify you when the server is started and will notify when an implant successfully calls back to the server for the first time.


## TO-DO

 - [x]  Add a menu a la meterpreter
 - [ ] Install Script that satisfies most of the dependencies automatically
 - [ ] Encrypt the tcp connections. Maybe some Rot13 ?
 - [x] Better error handling. TO-DO: Sending the "/" or "\" over the wire crashes the agent and server. FIX IT.
 - [x] Load and execute shellcode with CreateThread
 - [x] Inject Shellcode via Process Hollowing //Removed since kinda pointless having multiple techniques
 - [x] Inject Shellcode with CreateRemoteThread
 - [ ] A cross platform builder script
 - [ ] Upload/Download (File transfer via RSync or SCP perhaps ?)
 - [x] Persistence
 - [ ] Keylogging
 - [ ] RickRoll

## Notes:
The project is under "active" development. I am just a monkey with a keyboard. I am working on this when I have the time and learning as I go.
There is also a theoretical hardlimit to how many agents you can have deployed since it relys on IDs generated from 1 to 9000. You can modify this if you want to raise the limit of the amount of agents you expect to connect to the server. The first time a implant is compiled, it could take a while if you are missing a few go packages. Afterwards implant generation should be quick. RUN THE SERVER BINARY FROM WITHIN THE PROJECT DIRECTORY OTHERWISE IT WON'T BE ABLE TO BUILD THE IMPLANT OR START HTTP SERVERS AND SUCH.
