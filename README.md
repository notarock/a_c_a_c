# a cheap a_n_i_v clone

A simple and lightweight Twitch chat bot that learns from chat messages and generates responses using a Markov chain.

## Features
- Reads and learns from chat messages in real-time
- Ignore parroted messages and message from other bot users
- Generates responses based on a Markov chain model
- Configurable message frequency to prevent spam
- Lightweight and easy to set up (Docker or binary provided)

## Requirements

- Golang 1.22
- A twitch user account with a oauth token

## Configuration

Create a `.env` file with the following:

```ini
BASE_PATH="path/to/files"
ENV="dev"  # production will enable sending messages

COUNTDOWN=5 # Messages count untill next message
TWITCH_USER="a_c_a_c"
TWITCH_OAUTH_STRING="oauth:<token>"
TWITCH_CHANNEL="channelone,channeltwo"
# TWITCH_CHANNEL="eddieuce"
```

Alternatively, you can provision this configuration via Environment variables.

## Usage

Run the bot by executing the provided binary:
```sh
./acac
```

Or by using the provided docker image, `ghcr.io/notarock/a_c_a_c:latest`

## Disclaimer
This bot is intended for entertainment purposes. Use it responsibly and adhere to Twitch’s Terms of Service.

## License
MIT License

