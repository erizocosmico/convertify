package main

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"regexp"
)

type Song struct {
	Title, Artist, Image, Url string
}

func findInContent(regex string, content []byte, minMatches int, desiredMatch int) string {
	re := regexp.MustCompile(regex)
	matches := re.FindSubmatch(content)
	if len(matches) >= minMatches {
		return string(matches[desiredMatch][:len(matches[desiredMatch])])
	} else {
		return ""
	}
}

func GetSongs(links []string) []Song {
	songs := make([]Song, len(links))
	current := 0
	for _, link := range links {
		if link != "" {
			resp, _ := http.Get(link)
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			title := findInContent("<h1 itemprop=\"name\">(.*)</h1>", body, 2, 1)
			artist := findInContent("<h2>(.?)by <a href=\"(.*)\">(.*)</a></h2>", body, 4, 3)
			image := findInContent("<img src=\"(.*)\" border=\"0\" alt=\"(.*)\" id=\"big-cover\"", body, 3, 1)
			ytUrl := "http://gdata.youtube.com/feeds/api/videos?max-results=1&alt=json&q=" + url.QueryEscape(title + " - " + artist)
			resp, _ = http.Get(ytUrl)
			defer resp.Body.Close()
			responseBody, _ := ioutil.ReadAll(resp.Body)
			videoId := findInContent("v=([a-zA-Z0-9_-]+)&feature=youtube_gdata_player\"}", responseBody, 2, 1)
			songUrl := "http://www.youtube.com/watch?v=" + videoId
			if image != "" && artist != "" && title != "" && songUrl != "" {
				songs[current] = Song {title, artist, image, songUrl}
				current++
			}
		}
	}
	return songs
}