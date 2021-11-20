package connection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)

type session struct {
	AccessToken     string      `json:"accessToken"`
	ClientId        string      `json:"clientId"`
	TokenExpiration int         `json:"accessTokenExpirationTimestampMs"`
	RefreshToken    string      `json:"refreshToken"`
	HttpClient      http.Client `json:"-"`
}

var BASEURL string = "api.spotify.com"
var USERAGENT string = "User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:71.0) Gecko/20100101 Firefox/71.0"
var FILENAME string = "creds.json"

//TODO find a way to not do this
var globalSession session
var httpSession http.Client

//TODO save this to file

func GetSession() *session {

	if reflect.DeepEqual(globalSession, session{}) {

		loaded, err := loadFromFile()
		if err != nil {
			newSession := NewSession()
			globalSession = *newSession
			saveToFile(*newSession)
			return newSession
		}

		RefreshToken(&loaded)
		globalSession = loaded
		saveToFile(loaded)
		return &loaded

	}
	if reflect.DeepEqual(globalSession.HttpClient, http.Client{}) {
		globalSession.HttpClient = *http.DefaultClient
	}

	if time.Now().Unix() > int64(globalSession.TokenExpiration/1000) {
		RefreshToken(&globalSession)

	}
	return &globalSession
}

func PrintSession(s session) {

	fmt.Printf("AccessToken: %s\nClientID: %s\nTokenExp: %d\nRefreshToken: %s",
		globalSession.AccessToken,
		globalSession.ClientId,
		globalSession.TokenExpiration,
		globalSession.RefreshToken,
	)
}

func saveToFile(s session) {

	file, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(FILENAME, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func loadFromFile() (session, error) {
	file, err := ioutil.ReadFile(FILENAME)
	data := session{}
	if err != nil {
		return data, err
	}
	_ = json.Unmarshal([]byte(file), &data)
	return data, nil
}

func NewSession() *session {

	fmt.Println("Getting new token. Wait a few seconds.")
	newSession := session{}
	GetToken(&newSession)
	fmt.Println("Successfully Logged In.")
	return &newSession
}

func CheckStatusCode(statusCode int, s *session) {
	switch statusCode {
	case 200:
		return
	case 401:
		RefreshToken(s)
		return
	default:
		log.Fatal("Got " + fmt.Sprintf("%d", statusCode) + " Code Response!")
	}
}

func RefreshToken(s *session) {
	query := "https://open.spotify.com/get_access_token?reason=transport"
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", USERAGENT)
	req.Header.Add("cookie", "sp_dc="+s.RefreshToken)
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-site")

	res, err := s.HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var resultRefresh session

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	if err := json.Unmarshal(body, &resultRefresh); err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatal("Got " + res.Status + " Status")
	}

	s.AccessToken = resultRefresh.AccessToken

}

func addHeaders(r *http.Request, s session) {
	r.Header.Add("Accept", "*/*")
	r.Header.Add("User-Agent", USERAGENT)
	r.Header.Add("Authorization", "Bearer "+s.AccessToken)
	r.Header.Add("Sec-Fetch-Mode", "cors")
	r.Header.Add("Sec-Fetch-Site", "same-site")
}

func (s *session) GetRequest(apiCall string) []byte {

	query := "https://" + BASEURL + apiCall
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
	return body
}
