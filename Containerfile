FROM golang:1.24.2-alpine

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
RUN go build -o bot *.go

ENV DISCORD_TOKEN "YOUR TOKEN HERE"

# Run the bot
CMD [ \
    "./bot", \
    "-b", \
    "-p", "tidal", \
    "-p", "spotify", \
    "-p", "youtube", \
    "-p", "youtubeMusic", \
    "-p", "appleMusic", \
    "-t", "MTI4MzA5OTczMTE3MzU3MjY5OQ.Gc256W.-0GBH0ut3iK2b-dz1twMJ83aWfO4vTVyev5EDM" \
    ]
