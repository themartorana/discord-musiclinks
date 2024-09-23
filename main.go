package main

import (
	"log"
	"musiclinks/discord"
	"musiclinks/provider"
	"os"
	"os/signal"
	"strings"

	"github.com/bigspawn/go-odesli"
	"github.com/spf13/pflag"
	"golang.org/x/exp/slices"
)

var (
	token     string
	debug     bool
	platforms []string

	url      string
	startBot bool
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if url == "" && !startBot {
		log.Fatal("URL or start bot flag must be specified")
	}
	if startBot && token == "" {
		log.Fatal("Token must be specified to start the bot")
	}
	if url != "" {
		log.Println("Getting links for", url)
		getLinksForUrl(url)
	} else if startBot {
		log.Println("Starting Discord bot...")
		if startBot {
			b := discord.StartBot(token, platforms...)

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)
			<-quit
			b.Close()
		}
	}
}

func getLinksForUrl(url string) {
	provider := provider.NewOdesliProvider()
	resp, err := provider.GetLinksForUrl(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Song:", resp.SongName)
	log.Println("Artist:", resp.Artist)
	for _, platformLink := range resp.PlatformLinks {
		if len(platforms) > 0 && !slices.Contains(
			platforms,
			string(platformLink.Platform),
		) {
			continue
		}
		log.Printf(
			"%s: %s\n",
			platformLink.Platform.ReadableName(),
			platformLink.Link,
		)
	}
}

func helpForPlatformsFlag() string {
	sb := strings.Builder{}
	sb.WriteString(`Platforms to return links for.
If no platforms are specified, all are returned.
Options are:
`)
	four := 0
	for _, platform := range odesli.AvailablePlatforms() {
		sb.WriteString(string(platform) + ", ")
		four++
		if four == 4 {
			sb.WriteString("\n")
			four = 0
		}
	}
	msg := sb.String()
	return msg[:len(msg)-2]
}

func init() {
	pflag.StringVarP(&token, "token", "t", "", "Discord bot token")
	pflag.StringVarP(&url, "url", "u", "", "Original song or album URL")

	pflag.BoolVarP(&startBot, "bot", "b", false, "Start the discord bot")
	pflag.BoolVarP(&debug, "debug", "d", false, "Enable debug mode")

	pflag.StringArrayVarP(
		&platforms,
		"platform",
		"p",
		[]string{},
		helpForPlatformsFlag(),
	)
	pflag.Parse()
}
