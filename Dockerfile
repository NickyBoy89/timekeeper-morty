FROM alpine:latest as builder

RUN apk add --update go

# Build the bot
COPY . .

RUN go build

# Run the bot
FROM alpine:latest

COPY --from=builder timekeeper-morty .

CMD ["/timekeeper-morty"]
