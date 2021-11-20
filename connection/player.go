package connection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Gabulhas/spotify_controller/model"
)

var PLAYERURL string = "api.spotify.com/v1/me/player"

// TODO: change parameter to interface, since function input could be string or []string
func (s *session) PlayType(uri string) {

	var bodyStruct interface{}

	if strings.Contains(uri, "track") {
		bodyStruct = struct {
			Uris []string `json:"uris"`
		}{
			Uris: []string{uri},
		}

	} else {
		bodyStruct = struct {
			Context_uri string `json:"context_uri"`
		}{
			Context_uri: uri,
		}

	}

	bodyString, err := json.Marshal(bodyStruct)
	fmt.Println("bodyString", string(bodyString))
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(
		"PUT",
		"https://"+PLAYERURL+"/play",
		bytes.NewBuffer(bodyString),
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+s.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	_, err = s.HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

}

/*
   var headers = {
       'Authorization': `Bearer ${this.auth}`,
       'Content-Type': 'application/json'
   };
   var action='play'
   await this.getStatus()
   if (this.isPlaying) action='pause'
   var options = {
       url: 'https://api.spotify.com/v1/me/player/'+action,
       method: 'PUT',
       headers: headers,
   };

   function callback(error, response, body) {
       if (!error && response.statusCode == 200) {
           console.log(body);
       }
   }

   request(options, callback);

*/

func (s *session) PlayPause() {
	action := "play"
	if s.GetIsPlaying() {
		action = "pause"
	}
	s.Action(action, "PUT")
}

func (s *session) GetStatus() model.State {

	query := "https://" + PLAYERURL
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		log.Fatal(err)
	}
	addHeaders(req, *s)

	res, err := s.HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	if res.StatusCode != 200 {
		log.Fatal("Got " + res.Status + " Status")
	}
	var responseState model.State
	err = json.Unmarshal(body, &responseState)
	if err != nil {
		log.Fatal(err)
	}
	return responseState

}

func (s *session) Skip() {
	s.Action("next", "POST")
}

func (s *session) Back() {
	s.Action("previous", "POST")
}

func (s *session) Action(action, reqtype string) {
	query := "https://" + PLAYERURL + "/" + action
	fmt.Println(query)
	req, err := http.NewRequest(
		reqtype,
		query,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	addHeaders(req, *s)

	_, err = s.HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

}

func (s *session) GetDuration() int {
	return s.GetStatus().Item.DurationMs

}

func (s *session) GetPlaybackName() string {
	return s.GetStatus().Item.Name
}

func (s *session) GetPlaybackAlbumName() string {
	return s.GetStatus().Item.Album.Name
}

func (s *session) GetPlaybackArtists() string {
	return model.Artists_to_list(s.GetStatus().Item.Artists)
}

func (s *session) GetProgress() int {
	return s.GetStatus().ProgressMs
}

func (s *session) GetProgressPercent() int {
	status := s.GetStatus()
	return (status.ProgressMs * 100) / status.Item.DurationMs
}

func (s *session) GetRepeateState() string {
	return s.GetStatus().RepeatState
}

func (s *session) GetShuffleState() bool {
	return s.GetStatus().ShuffleState
}

func (s *session) GetVolume() int {
	return s.GetStatus().Device.VolumePercent
}

func (s *session) GetIsPlaying() bool {
	return s.GetStatus().IsPlaying
}
