FROM golang:1.16 as builder

RUN apt-get update && apt-get -y upgrade

# Build the bot
COPY . timekeeper

WORKDIR timekeeper

RUN go build

# Run the bot
FROM golang:1.16

COPY --from=builder /go/timekeeper/timekeeper-morty .

CMD ["./timekeeper-morty"]
