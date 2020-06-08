# Triggered bot

A mostly silly bot that gets triggered by random words.

Provided with a source set of words, triggered-bot will randomly select a configurable percentage as trigger words.
If anyone makes a post that contains a trigger word, triggered-bot will kick the creator from the server, but send
them a dm with an apology and an invitation to rejoin. The messages posted to the channel and text in the apology are
also configurable and chosen at random from a set of templates.

**Note:** Triggered bot will not attempt to kick the server owner, nor does it have the ability to do so.
Currently, Triggered bot is deployed to my raspberry pi and randomly selected about 850 of the top 10000 most common english words longer than 3 letters.

## Requirements
* go 1.14

## Building

```bash
    go build
```

## Usage

```bash
    TOKEN=<<YOUR-BOT-TOKEN-HERE>> ./triggered-bot \
        --source words.txt `# source of words to be randomly sampled` \
        --sample-ratio 0.1 `# percentage of words to be randomly sampled` \
        --include include.txt `# words to be sampled unconditionally`
```
