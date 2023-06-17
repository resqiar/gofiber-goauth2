package main

import (
	"fiber-goauth/libs"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	// initialize godotenv to read all .env files
	godotenv.Load()

	// initialize new instance of fiber
	server := fiber.New()

	server.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Fiber!")
	})

	// create a config for google config
	conf := &oauth2.Config{
		ClientID:     os.Getenv("G_CLIENT_ID"),
		ClientSecret: os.Getenv("G_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("G_REDIRECT"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}

	server.Get("/google", func(c *fiber.Ctx) error {
		// create url for auth process.
		// we can pass state as someway to identify
		// and validate the login process.
		URL := conf.AuthCodeURL("not-implemented-yet")

		// redirect to the google authentication URL
		return c.Redirect(URL)
	})

	server.Get("/auth/callback", func(c *fiber.Ctx) error {
		// get state and from the google callback
		code := c.Query("code")

		// exchange code that retrieved from google via
		// URL query parameter into an access token.
		token, err := conf.Exchange(c.Context(), code)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// convert token to get the user data
		profile, err := libs.ConvertToken(token.AccessToken)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(profile)
	})

	// bind server to listen to port 8000
	// change the port as you like.
	server.Listen(":8000")
}
