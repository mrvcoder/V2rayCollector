# Introduction

**First project on the internet** that crawls v2ray configs from Telegram channels. And the list will update every 5 hours. 😋

# How to use this ?! 🤔


- ‼ Github banned github actions on my account :( so you can use the list below Or you can fork this repo and enable github actions on your account and use your own subs links :) 
-------------------------------

It is so easy just go ahead and download a V2ray Client App that **supports subscription link** and use these links as subscription link 🤩
Config Type|subscription link
-------------------------------|-----------------------------|
Vmess         |https://raw.githubusercontent.com/youfoundamin/V2rayCollector/main/vmess_iran.txt      |
ShadowSocks        |https://raw.githubusercontent.com/youfoundamin/V2rayCollector/main/ss_iran.txt  |
Trojan |https://raw.githubusercontent.com/youfoundamin/V2rayCollector/main/trojan_iran.txt|
Vless|https://raw.githubusercontent.com/youfoundamin/V2rayCollector/main/vless_iran.txt|
Mixed (configs of this are different)|https://raw.githubusercontent.com/youfoundamin/V2rayCollector/main/mixed_iran.txt|


## Todos
 - [x] Adding comments to functions
 - [x] Getting messges modular so it can be easy to edit
 - [ ] Add feature to only stores configs from present until x days ago
 - [x] Sort the stored configs (from latest to oldest)
 - [x] Optimze config exctraction (only get config and remove the dsc and other things)
 - [x] Read Channels from channels.csv (it should support {all_messages} flag)
 - [ ] Update README (add usage of configs in different os and move channels list to channels.csv)
 - [ ] Add support for v2ray configs that posted in json data
 - [ ] Add support for configing script to limit configs count in each files
 - [ ] Add support for testing v2ray configs and adding only correct and working configs in files
 - [x] Fix issue at removing duplicate lines ( duplicates won't create by script , some channels put duplicate configs in their chats :D )

# Telegram channels list that used as source 😉 
click [here](https://github.com/mrvcoder/V2rayCollector/blob/main/channels.csv) to see the list

If you know other telegram channels which they put V2ray Configs feel free to add pull request :)

