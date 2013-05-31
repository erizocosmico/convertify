package main

import (
	"net/http"
	"html/template"
	"strings"
)

var (
	staticHttp = http.NewServeMux()
	appPort = ":8888"
	rootUrl = "http://localhost:8888"
)

type ResultsPage struct {
	Grid template.HTML
	Links string
}

func songToHtml(dst string, song Song) string {
	return dst + "<div class='song'><div class='album_left'><div class='album' style=\"background-image: url('" + song.Image + "')\"></div>"+
						"</div><div class='song_info'><h1>"+ song.Title +"</h1><h2>by "+ song.Artist +"</h2><a href='"+ song.Url +"' class='youtube_link'>Listen on youtube</a>"+
						"</div><div class='clear'></div></div>"
}

func main() {
	staticHttp.Handle("/static/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
    	t, _ := template.ParseFiles("tpl/main.html")
		t.Execute(w, rootUrl)
	})
	http.HandleFunc("/submit", func (w http.ResponseWriter, r *http.Request) {
    	if r.Method == "POST" {
			links := strings.Split(r.FormValue("spotify_links"), "\n")
			for i, link := range links {
				links[i] = strings.Trim(link, " ")
			}
			spotifyLinks := make([]string, len(links))
			counter := 0
			for _, link := range links {
				if link != "" && (strings.Contains(link, "open.spotify.com/track/") || strings.Contains(link, "play.spotify.com/track/")) && !strings.HasSuffix(link, "spotify.com/track/") {
					if strings.Contains(link, "play.spotify.com/track/") {
						link = strings.Replace(link, "play", "open", 1)
					}
					spotifyLinks[counter] = link
					counter++
				}
			}
			songsHtmlResult := ""
			songUrls := ""
			songs := GetSongs(spotifyLinks)
			for _, song := range songs {
				if song.Title != "" {
					if songUrls == "" {
						songUrls = song.Url
					} else {
						songUrls = songUrls + "\n" + song.Url
					}
					songsHtmlResult = songToHtml(songsHtmlResult, song)
				}
			}
			data := ResultsPage{template.HTML(songsHtmlResult), songUrls}
			t, _ := template.ParseFiles("tpl/results.html")
			t.Execute(w, data)
		} else {
	    	t, _ := template.ParseFiles("tpl/main.html")
			t.Execute(w, rootUrl)
		}
	})
	http.HandleFunc("/static/", func (w http.ResponseWriter, r *http.Request) {
		staticHttp.ServeHTTP(w, r)
	})
	http.ListenAndServe(appPort, nil)
}