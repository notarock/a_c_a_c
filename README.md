# `a_c_a_c`, a cheap a_n_i_v clone

![Dev time](https://img.shields.io/endpoint?url=https://wakapi.notarock.lol/api/compat/shields/v1/notarock/interval:any/project:a_c_a_c&label=Dev%20time&color=blue)
![GitHub Tag](https://img.shields.io/github/v/tag/notarock/a_c_a_c)
![GitHub contributors](https://img.shields.io/github/contributors/notarock/a_c_a_c)
![GitHub License](https://img.shields.io/github/license/notarock/a_c_a_c)
[![Tuyauterie](https://github.com/notarock/a_c_a_c/actions/workflows/main.yml/badge.svg)](https://github.com/notarock/a_c_a_c/actions/workflows/main.yml)

A Twitch chatbot that learns from chat and generates messages using a Markov chain.
Similar to BinyotBot and a_n_i_v.

## Want it in your chat?

I host an instance of the bot. To request it for your channel, [open a hosting request](https://github.com/notarock/a_c_a_c/issues/new?template=hosting-request.md).

## Community

Want to request `a_c_a_c`, need help, want to suggest features, or report bugs? Join my Discord:
[![Discord](https://img.shields.io/discord/1534046853383983186?label=Discord&logo=discord)](https://discord.com/invite/YjNAs9k9b8)

## Features

- Learns from chat messages as they arrive.
- Saves chat history and generated messages to files. No database needed.
- Supports a different message frequency and ignored-bot list for each channel.
- Can ignore users copying the bot's messages.
- Filters generated messages using blocked strings or full messages.
- Sometimes generates funny responses.

## Requirements

- A Twitch account with an OAuth token.
- Docker or a [release binary](https://github.com/notarock/a_c_a_c/releases).
- Go 1.25+ if building from source.

## Configuration

Use a `.env` file or environment variables for the bot's settings, and a YAML file for the channels it joins.

### Environment variables

For Twitch chat, set `BASE_PATH`, `TWITCH_USER`, `TWITCH_OAUTH_STRING`, `CHANNEL_CONFIG`, and `COUNTDOWN`. Set `ENV=production` to enable sending messages.

| Variable | Description | Example |
| --- | --- | --- |
| `ENV` | Only `production` enables sending messages to chat. | `production` |
| `BASE_PATH` | Directory for message history: `channel.txt`, `channel-sent.txt`, and `channel-rejected.txt`. | `./data` |
| `COUNTDOWN` | Default number of chat messages between responses. Must be an integer, even if every channel has its own frequency. | `200` |
| `IGNORE_PARROTS` | Set to `true` to avoid learning from copies of the bot's messages. | `true` |
| `TWITCH_USER` | Twitch username used by the bot. | `a_c_a_c` |
| `TWITCH_OAUTH_STRING` | OAuth token for that account. | `oauth:123123123123123` |
| `PROHIBITED_STRINGS` | Comma-separated strings to block in generated messages. | `https://,twitch.tv,@` |
| `PROHIBITED_MESSAGES` | Comma-separated full messages to block. | `acac` |
| `CHANNEL_CONFIG` | Path to the channel configuration file. | `./channels.yaml` |

### Channel configuration

`bots` lists usernames to ignore in every channel. Each channel can add its own ignored usernames through `extra_bots`.

`frequency` sets the number of chat messages between responses. If omitted or set to `0`, it uses `COUNTDOWN`.

```yaml
bots:
  - "nightbot"
  - "streamelements"

channels:
  - name: "my_favorite_streamer"
    frequency: 200
    extra_bots:
      - "mod_helper"
      - "my_favorite_moderation_bot"
  - name: "streamer_two"
```

The [example configuration](./example-channels.yaml) also includes `allow_bits`. That setting currently has no effect because the bits-filtering code is disabled.

## Usage

### Binary

Download a binary from [releases](https://github.com/notarock/a_c_a_c/releases).

Create a `data` directory and save your channel configuration as `channels.yaml`. Create a `.env` file in the directory where you run the bot:

```dotenv
ENV=production
BASE_PATH=./data
COUNTDOWN=200
IGNORE_PARROTS=true
TWITCH_USER=your_bot_account
TWITCH_OAUTH_STRING=oauth:your_token
CHANNEL_CONFIG=./channels.yaml
```

Then run:

```sh
./acac
```

Use `ENV=dev` to run without sending messages to chat.

### Docker

With the `.env`, `data` directory, and `channels.yaml` from above:

```sh
docker run --env-file .env \
  -v "$(pwd)/data:/data" \
  -v "$(pwd)/channels.yaml:/channels.yaml:ro" \
  -e BASE_PATH=/data \
  -e CHANNEL_CONFIG=/channels.yaml \
  ghcr.io/notarock/a_c_a_c:latest
```

Message history is saved in the mounted `data` directory.

## Before adding the bot

Get the streamer's permission before adding it to a channel, and follow Twitch's Terms of Service.

The bot learns from whatever people say in chat. Expect gibberish and potentially offensive messages.

I am not responsible for whatever it says. Everyone in chat is.

## Contributing

Open an issue for bugs or feature requests. For code changes, fork the repository and submit a pull request.

## License

MIT License.
