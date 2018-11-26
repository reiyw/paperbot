# paperbot

[![Build Status](https://travis-ci.org/reiyw/paperbot.svg?branch=master)](https://travis-ci.org/reiyw/paperbot)

## Features

- Extract paper information from URL.
    - Simple formatting to avoid it takes much space.
    - More information as a thread with Japanese translation.
- Show top-10 trending papers on arXiv every day.
    - Powered by [Arxiv Sanity Preserver](http://www.arxiv-sanity.com/).
- And translation, btw.

## Usage

Fill `.env` file:

```.env
PAPERBOT_SLACK_TOKEN=
ARXIV_TREND_CHANNEL_ID=
BOT_USER_ID=
BOT_USER_NAME=
BOT_ICON_URL=
```

Run:

```bash
go build
./paperbot
```
