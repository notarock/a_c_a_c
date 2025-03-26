# `a_c_a_c`, a cheap a_n_i_v clone

A simple, lightweight Twitch chatbot that learns from chat messages and generates responses using a Markov chain algorithm—similar to BinyotBot and a_n_i_v.

> [!WARNING]  
> Before adding this bot to any channel, **get the streamer's permission.**  
> The bot generates gibberish about **50% of the time**, and people *will* figure out that it's not human sooner than you think.  
> This will likely get you **banned** in most cases.  
> **No, really—don't add this to someone's channel without their consent.**


## Features  
- Reads and learns from chat messages in real-time  
- Stores chat history of received and sent messages  
- Generates responses using a Markov chain model  
- Ignores parrots (users who repeatedly copy the bot’s messages)  
- Configurable message frequency to prevent spam  
- Lightweight and easy to set up (Docker or binary, no database needed)  
- Sometimes generates funny responses  

## Getting Started  

### Want `a_c_a_c` in your chat?  

I'm currently hosting a copy of the bot and can add your Twitch channel upon request.  

- Open an issue under `[Hosting Request]` (coming soon :tm: )
- Or reach out via email or discord if you find me.

## Requirements  
- **Golang 1.22**  
- **A Twitch account** with an OAuth token  

## Configuration

The bot can be configured using a `.env` file or via environment variables.

| ENV Variable | Description | Example |
| -------- | ------- | --- |
| `ENV` | Environment where the software is being executed from. Anything outside of "production" will not send any real messages. | `"dev"` / `"production"` |
| `BASE_PATH` | Path where the recieved/sent messages will be stored to. Files are stored as `channel.txt` and `channel-sent.txt`. | `"./"` |
| `COUNTDOWN` | Number of messages to read from chat before sending a message. | `0` |
| `IGNORE_PARROTS` | Flag to ignore users who copy the last bot's messages. Prevent learning from the bot's own gibberish. | `"true"`/`"false"` |
| `TWITCH_USER` | Username of the account which this bot operates under. | `"a_c_a_c"` |
| `TWITCH_OAUTH_STRING` | Your account's oauth string to authenticate with twitch chat. | `"oauth:123123123123123"` |
| `TWITCH_CHANNELS` | Comma separated list of twitch channels to connect to. Every channel gets their own message db and as a result, their own "chat personality". | `"bozo,thelegend27"` |
| `TWITCH_BOT_USERNAMES` | Comma separated list of twitch bot in the channel. Bot users are added to the ignore list. | `"nightbot,myownbot,funtoon"` |
| `PROHIBITED_STRINGS` | Comma separated list of strings that will not be sent by the bot. Use this to filter out links, user mentions, etc. | `"https://,twitch.tv,@"` |
| `PROHIBITED_MESSAGES` | Comma separated list of messages that will not be sent by the bot. Use this to filter out full messages that should not be sent. | `"acac"` |

## Usage

### CLI

Download and run the binary from the release tab:

```sh
./acac
```

### Docker

The bot is available as a Docker image:
```sh
docker run -v /your/local/path:/data -e BASE_PATH="/data" ghcr.io/notarock/a_c_a_c:latest
```
Make sure to mount `BASE_PATH` as a volume.

## Disclaimer
This bot is intended for entertainment purposes. Use it responsibly and adhere to Twitch’s Terms of Service.

And again:

> [!CAUTION]
> Only use this with a streamer's consent.

## Contributing

Contributions are welcome! 🚀

- Feature requests & bug reports → Open an issue.
- Code contributions → Fork, create a branch, and submit a PR.
- Want `a_c_a_c` in your chat? → Open an issue.

## License
MIT License

