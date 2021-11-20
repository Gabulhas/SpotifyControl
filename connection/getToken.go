package connection

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/spf13/viper"
)

var LOGINURL string = `https://accounts.spotify.com/en/login?continue=https:%2F%2Fopen.spotify.com%2F`
var USERNAMEINPUT string = "#login-username"
var PASSWORDINPUT string = "#login-password"

func GetToken(s *session) {

	// create context
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
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
		),
	)
	if err != nil {
		log.Fatal(err)
	}

}

func login(username, password string) chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.Navigate(LOGINURL),
		chromedp.WaitVisible(USERNAMEINPUT),
		chromedp.SendKeys(USERNAMEINPUT, username),
		chromedp.SendKeys(PASSWORDINPUT, password),
		chromedp.Click("#login-button"),
	}
}
