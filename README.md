# Requirements
- [Go](https://go.dev/doc/install)

# How to build
```bash 
git clone https://github.com/Jilwer/RemoteControlMe
cd ./RemoteControlMe
go build
```

# Usage
1. Update and config.example.toml file to your liking
2. Run your compiled binary `Ex: RemoteControlMe.exe`
3. Your remote is now hosted locally!
4. I suggest using [Cloudflare tunnels](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/downloads/) to get online publically
 - Install from the above link
 - For a quick tunnel with a ephemeral domain just run `cloudflared tunnel --url localhost:8080`
 - To link it to your own domain follow this guide: https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/configure-tunnels/remote-management/


Example use of https://github.com/Jilwer/VRChatOscInput

Powered by: https://github.com/jfyne/live
