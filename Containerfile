FROM golang:1-alpine

# Set the working directory
WORKDIR /app

# Install git
# RUN apt-get update && apt-get install -y git
RUN apk update
RUN apk add git

# Clone the repository
RUN git clone https://github.com/themartorana/discord-musiclinks.git .
RUN git pull --all

# Install the required packages
RUN go build -ldflags "-X main.buildHash=$(git rev-parse --short HEAD)" -o bot *.go

ENV DISCORD_TOKEN ""
ENV MUSICLINKS_PLATFORMS "tidal,spotify,youtube,youtubeMusic,appleMusic"

# Run the bot
CMD ["./bot", "-b"]
