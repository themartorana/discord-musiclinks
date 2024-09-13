package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/bigspawn/go-odesli"
)

type OdesliProvider struct {
	odesliClient *odesli.Client
}

type PlatformLinkResponse struct {
	SongName      string
	Artist        string
	PlatformLinks []PlatformLink
}

type PlatformLink struct {
	Platform odesli.Platform
	Link     string
}

func NewOdesliProvider() *OdesliProvider {
	odesliClient, _ := odesli.NewClient(odesli.ClientOption{})
	return &OdesliProvider{
		odesliClient: odesliClient,
	}
}

func (d *OdesliProvider) GetLinksForUrl(url string) (PlatformLinkResponse, error) {
	ctx := context.Background()
	resp, err := d.odesliClient.GetLinks(
		ctx,
		odesli.GetLinksRequest{URL: url},
	)
	if err != nil {
		err = fmt.Errorf(
			"error getting links for url %s: %s",
			url,
			err,
		)
		return PlatformLinkResponse{}, err
	}

	platformLinks := make([]PlatformLink, 0)
	for platform, links := range resp.LinksByPlatform {
		url := d.cleanURL(links.Url)
		platformLinks = append(platformLinks, PlatformLink{
			Platform: platform,
			Link:     url,
		})
	}

	// return platformLinks, nil
	return PlatformLinkResponse{
		SongName:      d.getSongFromResponse(&resp),
		Artist:        d.getArtistFromResponse(&resp),
		PlatformLinks: platformLinks,
	}, nil
}

func (*OdesliProvider) cleanURL(url string) string {
	if strings.Contains(url, "listen.tidal.com") {
		return strings.Replace(url, "listen.tidal.com", "tidal.com", 1)
	}
	return url
}

// getSongFromResponse returns the song name from the response
// We're opting to trust Spotify the most
func (d *OdesliProvider) getSongFromResponse(resp *odesli.GetLinksResponse) string {
	for key, val := range resp.EntitiesByUniqueId {
		if strings.HasPrefix(key, "SPOTIFY") {
			return val.Title
		}
	}

	return ""
}

// getArtistFromResponse returns the artist name from the response
// We're opting to trust Spotify the most
func (d *OdesliProvider) getArtistFromResponse(resp *odesli.GetLinksResponse) string {
	for key, val := range resp.EntitiesByUniqueId {
		if strings.HasPrefix(key, "SPOTIFY") {
			return val.ArtistName
		}
	}
	return ""
}
