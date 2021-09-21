# Timekeeper-Morty

A simple discord bot that I wrote in Go to keep track of all my friend's timezones when in college

## Requirements

1. A version of `Go` installed

## Running

1. Run `go build`
1. Run the created binary
  * MacOS/Linux: `./timekeeper-morty`

## Commands

* `!settime` Sets the current user's time to a [TZ timezone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
* `!timefor` Gets the mentioned user/s times, from the sender's perspective
