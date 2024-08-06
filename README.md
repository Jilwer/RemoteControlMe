# RemoteControlMe
A WIP web app to allow users to remote control you on VRChat

## Screenshot
Please note that this does not portray a finished application

<img src="https://github.com/user-attachments/assets/252b6a84-a03a-4c73-a5f3-b8c4c690ec91" width="75%" alt="Screenshot">

## Requirements
- [Go](https://go.dev/doc/install)

## How to build
```bash 
git clone https://github.com/Jilwer/RemoteControlMe
cd ./RemoteControlMe
go build
```


## Usage
1. To configure, either modify config.example.toml to your preferences and rename it to config.toml, or run the program once to generate a default configuration file which you can then adjust as needed before restarting the application.
2. Run your compiled binary `Ex: RemoteControlMe.exe`
3. Your remote is now hosted locally!

4. I suggest using [Cloudflare tunnels](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/downloads/) to get online publicly
   - Install from the above link
   - For a quick tunnel with an ephemeral domain run `cloudflared tunnel --url localhost:8080`
   - To link it to your own domain follow this guide: [Configure Tunnels](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/configure-tunnels/remote-management/)


## Powered By
- https://github.com/jfyne/live
- https://github.com/Jilwer/VRChatOscInput
