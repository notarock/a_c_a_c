# a cheap a_n_i_v clone

A simple and lightweight Twitch chat bot that learns from chat messages and generates responses using the Markov chain algorithm.

> [!WARNING]  
> Before adding this bot to any channel, make sure that you have the streamer's permission. The messages this bot generates are clearly gibberish 50% of the time, and people *will* find out the user is not human sooner than you think. This *will* leads to you getting banned in most cases.
> 
> No, really. *Don't add this to someone's channel without their consent.*

## Features
- Reads and learns from chat messages in real-time
- Store a chat history of messages recieved and sent
- Generates responses based on a Markov chain model
- Ignore messages from "parrots" i.e. people who copy the bot's messages over and over.
- Configurable message frequency to prevent spam
- Lightweight and easy to set up (Docker or binary provided, no database needed)
- Is sometime funny

## I want this thing on my channel. What do I do?

I am currently hosting a copy of it and I have no problem adding your twitch channel to the list of joined channel.

 Ask politely. Open an issue to get in touch, or reach out in some ways or another idk.

There will be an issue template for [Hosting Requests] soon :tm:  

## Requirements

- Golang 1.22
- A twitch user account with a oauth token

## Configuration

Create a `.env` file with the following (`.env.example`):

```ini
BASE_PATH="path/to/files"
ENV="dev"  # production will enable sending messages

COUNTDOWN=5 # Messages count untill next message
TWITCH_USER="a_c_a_c"
TWITCH_OAUTH_STRING="oauth:<token>"
TWITCH_CHANNEL="channelone,channeltwo"
```

Alternatively, you can provision this configuration using Environment variables.

## Usage

### CLI

Run the bot by downloading and executing the provided binary:
```sh
./acac
```

### Docker

There is also a Docker image available at `ghcr.io/notarock/a_c_a_c:latest`.

To make this run properly, make sure to mount the `BASE_PATH` as a volume within the container. 


## Disclaimer
This bot is intended for entertainment purposes. Use it responsibly and adhere to Twitch’s Terms of Service.

And again:

> [!CAUTION]
> Only use this with a streamer's consent.

## License
MIT License

