# `a_c_a_c`, a cheap a_n_i_v clone 

![](https://img.shields.io/endpoint?url=https://wakapi.notarock.lol/api/compat/shields/v1/notarock/interval:any/project:a_c_a_c&label=Dev%20time&color=blue)
![GitHub Tag](https://img.shields.io/github/v/tag/notarock/a_c_a_c)
![GitHub contributors](https://img.shields.io/github/contributors/notarock/a_c_a_c)
![GitHub License](https://img.shields.io/github/license/notarock/a_c_a_c)


A simple, lightweight Twitch chatbot that learns from chat messages and generates responses using a Markov chain algorithm—similar to BinyotBot and a_n_i_v.

If you enjoy a_c_a_c, consider leaving a ⭐ on the repository to show your support!

> [!TIP]
> ### Want `a_c_a_c` to join your Twitch chat?
> I'm currently hosting an instance of the bot and can add your Twitch channel upon request.
>
> To request `a_c_a_c` in your channel, simply [open an issue using the Hosting Request template](https://github.com/notarock/a_c_a_c/issues/new?template=hosting-request.md) and provide the required information. Once submitted, I’ll take care of the rest and the bot should join your chat shortly after.

## Features  
- Reads and learns from chat messages in real-time  
- Stores chat history of received and sent messages  
- Generates responses using a Markov chain model  
- Ignores parrots (users who repeatedly copy the bot’s messages)  
- Configurable message frequency to prevent spam  
- Lightweight and easy to set up (Docker or binary, no database needed)  
- Sometimes generates funny responses  

## Getting Started  

## Requirements  
- **Golang 1.22**  
- **A Twitch account** with an OAuth token  

## Configuration

The bot can be configured using a `.env` file or via environment variables. 

In addition, channels are managed via a configuration file to allow better customization of behaviour in different channels.

### Env Variables

Here's a table with all the environment variables required to run this properly:

| ENV Variable | Description | Example |
| -------- | ------- | --- |
| `ENV` | Environment where the software is being executed from. Anything outside of "production" will not send any real messages. | `"dev"` / `"production"` |
| `BASE_PATH` | Path where the recieved/sent messages will be stored to. Files are stored as `channel.txt` and `channel-sent.txt`. | `"./"` |
| `COUNTDOWN` | Number of messages to read from chat before sending a message. | `0` |
| `IGNORE_PARROTS` | Flag to ignore users who copy the last bot's messages. Prevent learning from the bot's own gibberish. | `"true"`/`"false"` |
| `TWITCH_USER` | Username of the account which this bot operates under. | `"a_c_a_c"` |
| `TWITCH_OAUTH_STRING` | Your account's oauth string to authenticate with twitch chat. | `"oauth:123123123123123"` |
| `PROHIBITED_STRINGS` | Comma separated list of strings that will not be sent by the bot. Use this to filter out links, user mentions, etc. | `"https://,twitch.tv,@"` |
| `PROHIBITED_MESSAGES` | Comma separated list of messages that will not be sent by the bot. Use this to filter out full messages that should not be sent. | `"acac"` |
| `CHANNEL_CONFIG` | Path to channels configuration file. See `Channel Configuration File` section bellow. | `"./channels.yaml"` |

### Channel Configuration File

The configuration file defines how a_c_a_c behaves across different Twitch channels. It includes a global list of bot usernames to ignore . This prevents the bot from mimicking repetitive, bot-like behavior and ensures cleaner, more human-like interactions.

Each channel entry can specify its own settings, such as message frequency, whether to respond to bits, and additional bots to ignore. Any omitted settings will fall back to sensible defaults—for example, frequency defaults to the value of the COUNTDOWN environment variable, and allow_bits defaults to false.

See also: [Example config file](./example-channels.yaml)

```yaml
bots:
  - "nightbot"
  - "streamerelements"

channels:
  - name: my_favorite_streamer
    frequency: 200
    allow_bits: true
    extra_bots:
      - "mod_helper"
      - "my_favorite_moderation_bot"
  - name: "streamer_two"
```

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

It learns from the chatroom. As a result, I am not responsible for whatever it says. *Everyone in chat is*

> [!WARNING]  
> Before adding this bot to any channel, **get the streamer's permission.**  
> The bot generates gibberish about **50% of the time**, and people *will* figure out that it's not human sooner than you think.  
> This will likely get you **banned** in most cases.  
> **No, really—don't add this to someone's channel without their consent.**

## Contributing

Contributions are welcome! 🚀

- Feature requests & bug reports → Open an issue.
- Code contributions → Fork, create a branch, and submit a PR.
- Want `a_c_a_c` in your chat? → [Use the `Hosting Request` issue template](https://github.com/notarock/a_c_a_c/issues/new?template=hosting-request.md).

## License
MIT License

