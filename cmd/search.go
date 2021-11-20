/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Gabulhas/spotify_controller/connection"
	"github.com/Gabulhas/spotify_controller/model"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

var queryTypes = []string{"artist", "album", "playlist", "track"}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for music.",
	Long:  ``,
	Run:   search_run,
}

func init() {
	rootCmd.AddCommand(searchCmd)

}

//ARGS: Type (Artist, Playlist, ...) Name
func search_run(cmd *cobra.Command, args []string) {
	query := ""
	queryType := ""
	if len(args) == 0 {
		idx, err := fuzzyfinder.Find(
			queryTypes,
			func(i int) string {
				return queryTypes[i]
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		queryType = queryTypes[idx]
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your search: ")
		query, _ = reader.ReadString('\n')
	} else {
		queryType = args[0]
		for i, word := range args[1:] {
			if i != 0 {
				query = query + " "
			}
			query = query + word
		}
	}

	query = strings.TrimSpace(query)
	queryType = strings.TrimSpace(queryType)
	query_api(query, queryType)

}

//TODO: change this so it can also display in console (not on fuzzy finder)
func query_api(query, queryType string) {

	query = strings.ReplaceAll(query, " ", "+")
	resp := connection.GetSession().GetRequest(
		"/v1/search?q=" + query + "&type=" + queryType + "&limit=30",
	)
	result_uri := ""

	switch queryType {
	case "track":
		var trackResponse model.TrackResponse
		err := json.Unmarshal(resp, &trackResponse)
		if err != nil {
			log.Fatal(err)
		}
		result_uri = track_display(trackResponse)

	case "artist":
		var artistResponse model.ArtistResponse
		err := json.Unmarshal(resp, &artistResponse)
		if err != nil {
			log.Fatal(err)
		}
		result_uri = artist_display(artistResponse)

	case "album":
		var albumResponse model.AlbumResponse
		err := json.Unmarshal(resp, &albumResponse)
		if err != nil {
			log.Fatal(err)
		}
		result_uri = album_display(albumResponse)

	case "playlist":
		var playlistResponse model.PlaylistResponse
		err := json.Unmarshal(resp, &playlistResponse)
		if err != nil {
			log.Fatal(err)
		}
		result_uri = playlist_display(playlistResponse)
	default:
		fmt.Printf(string(resp))
	}

	if IsInteractive {
		PlayByUri(result_uri)
	}
}

func track_display(resp model.TrackResponse) string {

	if !IsInteractive {
		result_string := ""
		for _, item := range resp.Tracks.Items {
			if item.Name == "" {
				continue
			}
			result_string = result_string + fmt.Sprintf("%s - %s - %s - %s\n", item.Name, artists_to_list(item.Artists), item.Album.Name, item.URI)
		}
		fmt.Println(result_string)
		return ""
	}

	idx, err := fuzzyfinder.Find(
		resp.Tracks.Items,
		func(i int) string {
			item := resp.Tracks.Items[i]
			artistList := artists_to_list(item.Artists)
			return fmt.Sprintf("%s - %s", artistList, item.Name)
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}

			//TODO: remove duplicated code
			item := resp.Tracks.Items[i]
			artistList := artists_to_list(item.Artists)

			return fmt.Sprintf(
				"Artist(s): %s\nTrack: %s\nAlbum: %s\nDuration %s\n",
				artistList,
				item.Name,
				item.Album.Name,
				fmt.Sprintf("%d:%d", (item.DurationMs/(1000*60))%60, (item.DurationMs/1000)%60),
			)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Tracks.Items[idx].URI
}

func artist_display(resp model.ArtistResponse) string {
	if !IsInteractive {
		result_string := ""
		for _, item := range resp.Artists.Items {
			if item.Name == "" {
				continue
			}
			result_string = result_string + fmt.Sprintf("%s - %s\n", item.Name, item.URI)
		}
		fmt.Println(result_string)
		return ""
	}
	idx, err := fuzzyfinder.Find(
		resp.Artists.Items,
		func(i int) string {
			item := resp.Artists.Items[i]
			return fmt.Sprintf("%s", item.Name)
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}

			//TODO: remove duplicated code
			item := resp.Artists.Items[i]

			return fmt.Sprintf(
				"Name: %s\nPopularity: %d\nFollowers: %d",
				item.Name,
				item.Popularity,
				item.Followers.Total,
			)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Artists.Items[idx].URI
}

func album_display(resp model.AlbumResponse) string {
	if !IsInteractive {
		result_string := ""
		for _, item := range resp.Albums.Items {
			if item.Name == "" {
				continue
			}
			result_string = result_string + fmt.Sprintf("%s - %s - %s - %s\n", item.Name, artists_to_list(item.Artists), item.Type, item.URI)
		}
		fmt.Println(result_string)
		return ""
	}
	idx, err := fuzzyfinder.Find(
		resp.Albums.Items,
		func(i int) string {
			item := resp.Albums.Items[i]
			artists := artists_to_list(item.Artists)
			return fmt.Sprintf("%s - %s", item.Name, artists)
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}

			//TODO: remove duplicated code
			item := resp.Albums.Items[i]
			artists := artists_to_list(item.Artists)

			return fmt.Sprintf(
				"Name: %s\nArtists: %s\nAlbum type: %s\nRelease Date: %s\nType: %s",
				item.Name,
				artists,
				item.AlbumType,
				item.ReleaseDate,
				item.Type,
			)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Albums.Items[idx].URI
}

func playlist_display(resp model.PlaylistResponse) string {
	if !IsInteractive {
		result_string := ""
		for _, item := range resp.Playlists.Items {
			if item.Name == "" {
				continue
			}
			result_string = result_string + fmt.Sprintf("%s - %s - %s\n", item.Name, item.Owner.DisplayName, item.URI)
		}
		fmt.Println(result_string)
		return ""
	}

	idx, err := fuzzyfinder.Find(
		resp.Playlists.Items,
		func(i int) string {
			item := resp.Playlists.Items[i]

			return fmt.Sprintf("%s - %s", item.Name, item.Owner.DisplayName)
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}

			//TODO: remove duplicated code
			item := resp.Playlists.Items[i]

			return fmt.Sprintf(
				"Name: %s\nDescription: %s\nOwner: %s\nTotal Tracks%d",
				item.Name,
				item.Description,
				item.Owner.DisplayName,
				item.Tracks.Total,
			)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Playlists.Items[idx].URI

}

func artists_to_list(artists []model.Artist) string {
	artistList := ""
	for i, artist := range artists {
		if i != 0 {
			artistList = artistList + ","
		}
		artistList = artistList + artist.Name
	}
	return artistList
}
