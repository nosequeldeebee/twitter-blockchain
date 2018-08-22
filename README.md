[![Build Status](https://travis-ci.com/mycoralhealth/twitter-blockchain.svg?branch=master)](https://travis-ci.com/mycoralhealth/twitter-blockchain)

# twitter-blockchain
Build your own Blockchain Twitter Recorder in Go

Read this [tutorial](https://medium.com/@mycoralhealth/build-your-own-blockchain-twitter-recorder-in-go-4fa504e912c3) first.

### Setup

- Clone this repo and navigate to it
- [Get your API keys](https://developer.twitter.com/en/apply-for-access.html) from Twitter. In particular you will need: APIKey, APISecret, Token, TokenSecret
- Rename `example.env` to `.env` and input your Twitter API credentials you just got
- Pick a port number e.g. 9090 and put that in `.env`

### Usage

- `go run main.go`
- open a browser and visit `localhost:{port number}/{twitter handle}` e.g. localhost:9090/nosequeldeebee

![ScreenShot](https://s15.postimg.cc/4y9hujn8b/Screen_Shot_2018-08-20_at_9.18.41_AM.png)

- copy the text and paste into a [SHA256 converter](https://passwordsgenerator.net/sha256-hash-generator/) to get the hash

![ScreenShot](https://s15.postimg.cc/jiv630cor/Screen_Shot_2018-08-20_at_9.21.43_AM.png)

- Enter the hash as the block data as seen in the [tutorial](https://medium.com/@mycoralhealth/build-your-own-blockchain-twitter-recorder-in-go-4fa504e912c3)
