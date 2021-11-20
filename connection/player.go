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
	if s.GetStatus().IsPlaying {
		action = "pause"
	}
	query := "https://" + PLAYERURL + "/" + action
	fmt.Println(query)
	req, err := http.NewRequest(
		"PUT",
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
