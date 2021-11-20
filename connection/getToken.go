package connection

import (
	"context"
	"log"

	"github.com/chromedp/cdproto/storage"
	"github.com/chromedp/chromedp"
	"github.com/spf13/viper"
)

var LOGINURL string = `https://accounts.spotify.com/en/login?continue=https:%2F%2Fopen.spotify.com%2F`
var USERNAMEINPUT string = "#login-username"
var PASSWORDINPUT string = "#login-password"

type configType struct {
	AccessToken                      string `json:"accessToken"`
	AccessTokenExpirationTimestampMs int    `json:"accessTokenExpirationTimestampMs"`
}

func GetToken(s *session) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	username, okUsername := viper.Get("USERNAME").(string)
	password, okPassword := viper.Get("PASSWORD").(string)

	if !okPassword || !okUsername {
		log.Fatal("Both USERNAME and PASSWORD variables must be defined in the .env file.")
	}

	err := chromedp.Run(ctx,
		login(
			username,
			password,
			s,
		),
	)

	if err != nil {
		log.Fatal(err)
	}

}

func login(username, password string, s *session) chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.Navigate(LOGINURL),
		chromedp.WaitVisible(USERNAMEINPUT),
		chromedp.SendKeys(USERNAMEINPUT, username),
		chromedp.SendKeys(PASSWORDINPUT, password),
		chromedp.Click("#login-button"),
		chromedp.WaitVisible(".Root__main-view"),
		chromedp.Evaluate("JSON.parse(document.getElementById('config').text.trim())", s),
		chromedp.ActionFunc(func(ctx context.Context) error {
			result, err := storage.GetCookies().Do(ctx)

			if err != nil {
				log.Fatal(err)
			} else {
				for _, cookie := range result {
					if cookie.Name == "sp_dc" {
						s.RefreshToken = cookie.Value
						break
					}
				}
			}
			return err
		}),
	}
}
