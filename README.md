<div align="center">
<pre>
 ______   __  __  __  __    
/\__  _\ /\ \/\ \/\ \/\ \   
\/_/\ \/ \ \ `\\ \ \ \ \ \  
   \ \ \  \ \ , ` \ \ \ \ \ 
    \_\ \__\ \ \`\ \ \ \_/ \
    /\_____\\ \_\ \_\ `\___/
    \/_____/ \/_/\/_/`\/__/ 
<br>
CS2 inventory tracker in the terminal
<br>
<img alt="GitHub License" src="https://img.shields.io/github/license/ItzAfroBoy/inv"> <a href="https://www.codefactor.io/repository/github/itzafroboy/inv"><img src="https://www.codefactor.io/repository/github/itzafroboy/inv/badge" alt="CodeFactor" /></a> <img alt="GitHub code size in bytes" src="https://img.shields.io/github/languages/code-size/ItzAfroBoy/inv">
</pre>
</div>

## Installation

### Install with Go

```shell
go install github.com/ItzAfroBoy/inv@latest
inv ...
```

### Build from source

```shell
git clone https://github.com/ItzAfroBoy/inv
cd inv
go install
inv ...
```

## Usage

`Usage: inv [--user NAME | STEAMID] [--use-csf] [--csf-key KEY] [--export] [--import] [--notify] [--prices] [--print-config] [--sort TYPE]`  

- `--use-csf`: Uses CSFloat API to retrieve inventory data. Use `--csf-key` if not in config file
- `--export`: Exports inventory as JSON file
- `--import`: Import inventory from JSON file
- `--notify`: Send notification to a NTFY endpoint on inventory retrieval completion
- `--prices`: Fetch item prices from Steam Market
- `--print-config`: Prints config file to the terminal
- `--sort`: Sort table by `price`, `item`, `collection` or `float`

Powered by:

- [SteamID](https://steamid.io)
- [Steam Community API](https://steamcommunity.com/inventory)
- [Steam Market API](https://steamcommunity.com/market/priceoverview)
- [CSFloat API](https://api.csfloat.com)
