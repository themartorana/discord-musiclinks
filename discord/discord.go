package discord

import (
	"log"
	"musiclinks/provider"
	"regexp"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls/v2"
)

type DiscordBot struct {
	token     string
	platforms []string

	session       *discordgo.Session
	xurlsSearcher *regexp.Regexp
	provider      *provider.OdesliProvider

	end chan struct{}
}

var patterns = []string{
	"open.spotify.com",
	"tidal.com",
	"music.apple.com/us/album",
	"music.youtube.com",
	"music.amazon.com",
	"pandora.com",
	"deezer.com",
	"soundcloud.com",
	"napster.com",
}

func StartBot(token string, platforms ...string) *DiscordBot {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	b := &DiscordBot{
		token:         token,
		session:       session,
		platforms:     platforms,
		end:           make(chan struct{}),
		xurlsSearcher: xurls.Strict(),
		provider:      provider.NewOdesliProvider(),
	}

	session.AddHandler(b.HandleMessage)
	err = session.Open()
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func (b *DiscordBot) HandleMessage(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// Ignore own messages
	if msg.Author.Bot {
		return
	}

	go b.processMessage(msg.Message)
}

func (b *DiscordBot) processMessage(msg *discordgo.Message) {
	urls := b.xurlsSearcher.FindAllString(msg.Content, -1)
	if len(urls) == 0 {
		return
	}

	// Check if any of the links are music links
	var musicLinks []string
	for _, url := range urls {
		for _, pattern := range patterns {
			if strings.Contains(url, pattern) && !slices.Contains(musicLinks, url) {
				musicLinks = append(musicLinks, url)
			}
		}
	}

	for _, link := range musicLinks {
		go b.processMusicLink(link, msg)
	}
}

func (b *DiscordBot) processMusicLink(link string, msg *discordgo.Message) {
	response, err := b.provider.GetLinksForUrl(link)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf(
		"Got %d links for %s by %s %s\n",
		len(response.PlatformLinks),
		response.SongName,
		response.Artist,
		link,
	)

	// Parsed links
	parsedLinks := make(map[string]string)
	for _, platformLink := range response.PlatformLinks {
		if slices.Contains(b.platforms, string(platformLink.Platform)) {
			parsedLinks[platformLink.Platform.ReadableName()] = platformLink.Link
		}
	}

	if len(parsedLinks) > 0 {
		b.respondToMessageWithLinks(
			response.Artist,
			response.SongName,
			parsedLinks,
			msg,
		)
	}
}

func (b *DiscordBot) respondToMessageWithLinks(
	artist,
	songName string,
	links map[string]string,
	msg *discordgo.Message,
) {
	var sb strings.Builder
	sb.WriteString("Links for ")
	if songName != "" {
		sb.WriteString("_" + songName + "_")
		if artist != "" {
			sb.WriteString(" by ")
			sb.WriteString("_" + artist + "_")
		}
		sb.WriteString(":\n")
	} else {
		sb.WriteString("the song you requested:\n")
	}

	for platform, link := range links {
		sb.WriteString(platform)
		sb.WriteString(": ")
		sb.WriteString(link)
		sb.WriteString("\n")
	}

	message := discordgo.MessageSend{
		Content:   sb.String(),
		Reference: msg.Reference(),
		Flags: discordgo.MessageFlagsSuppressEmbeds |
			discordgo.MessageFlagsSuppressNotifications,
	}

	_, _ = b.session.ChannelMessageSendComplex(
		msg.ChannelID,
		&message,
	)
}

func (b *DiscordBot) Close() {
	b.session.Close()
	b.session = nil
	b.end <- struct{}{}
	close(b.end)
}

func (b *DiscordBot) EndChan() <-chan struct{} {
	return b.end
}
