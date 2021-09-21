FROM ubuntu:latest as builder

RUN apt-get update && apt-get -y upgrade

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install libxss1 golang-go

# Build the bot
COPY . .

RUN go build

# Run the bot
FROM ubuntu:latest

COPY --from=builder timekeeper-morty .

CMD ["/timekeeper-morty"]
