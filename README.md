# RingBot

RingBot is a Discord bot that can call phone numbers using [Twilio](https://twilio.com/).

The bot provides bi-directional audio, allowing either side to hear and speak to the other.

## Disclaimer

**This bot is in no way meant to be used on an already public bot. This is meant to be a proof of concept, as exposing it to other users can result in unintended consequences with the numbers they could call.**

## Getting Started

### Twilio Setup

<img align="right" alt="Twilio Logo" src="https://upload.wikimedia.org/wikipedia/commons/c/c0/Twilio_logo.png" width="300">

To use the bot, you need a [Twilio](https://twilio.com/) account with [funds added](https://console.twilio.com/us1/billing/manage-billing/billing-overview).

You also need an [active number](https://console.twilio.com/us1/develop/phone-numbers/manage/incoming), purchasable [here](https://console.twilio.com/us1/develop/phone-numbers/manage/search). The bot will currently select the first number in the list.

Once setup, fill out `username` and `password` in settings.json with the Account SID and Auth Token respectfully. Both are found on the [Twilio console](https://console.twilio.com/).

### Websocket Setup

To allow Twilio to contact the `/mediastream` websocket within the bot (to hear voice audio and speak), you will need a publically exposed IP.

This is available locally by port-forwarding `:8000` and using your [IPv4](https://www.whatismyip.com/) IP.

Another available tool is downloading [ngrok](https://ngrok.com/download), running `ngrok http 8000` and copying everything except `https://` from the URL.

### Bot Token

Create a Discord bot, or use an existing one from the [Developer Portal](https://discord.com/developers/applications) and grab the token from there.

## Running the bot

### Starting a Call

Starting a call is as easy as typing `/call {number}`. Any normal phone formatting will work (e.g. (123) 456-7890 or 123-456-7890).

To start a call, you must be in a voice channel the bot can see.

### Speakerphone System

To prevent overlapping audio, the bot only listens to the person who typed the call command.

If another member of the voice channel wants to speak through the phone, clicking 'Get Speaker' will then only listen to the user who clicked it.

### Key Pad

If you wish to send a keypad number over the phone, clicking on 'Key pad' will popup a modal to enter any number of digits to play.

`01234567890*#w` are all acceptable to input. `w` gives a 0.5 second delay. All inputted numbers already have this 0.5 second delay between each.

### Mute/unmute Call

If you wish for the other side to stop hearing you, clicking the 'Mute/unmute call' button will stop the other side from hearing you until you unmute.

### End Call

Self-explanatory. The call will end and the bot will leave the voice channel.

## Contributing

**PRs are highly encouraged.** This code was a result of a roughly 48-hour challenge to create a tool like this, and with that multiple bugs are sure to show.

If you are less tech-savvy or don't understand Golang that well, please open an issue and a contributor or I will attempt to solve the issue. Include any relevant logging found in the console.
