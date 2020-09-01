package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

// StravaBase Strava base url
const StravaBase string = "http://www.strava.com/oauth"

// ErrorResponse Struct ...
type ErrorResponse struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"code"`
}

// AthleteResponse Struct ...
type AthleteResponse struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	ResourceState int    `json:"resource_state"`
	FirstName     string `json:"firstname"`
	LastName      string `json:"lastname"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	Sex           string `json:"sex"`
	Premium       bool   `json:"premium"`
	Summit        bool   `json:"summit"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	BadgeTypeID   int    `json:"badge_type_id"`
	ProfileMedium string `json:"profile_medium"`
	Profile       string `json:"profile"`
	Friend        int    `json:"friend"`
	Follower      int    `json:"follower"`
}

// OauthResponse Struct ...
type OauthResponse struct {
	Athlete      AthleteResponse `json:"athlete"`
	TokenType    string          `json:"token_type"`
	ExpiresAt    int             `json:"expires_at"`
	ExpriresIn   int             `json:"expires_in"`
	RefreshToken string          `json:"refresh_token"`
	AccessToken  string          `json:"access_token"`
	Message      string          `json:"message"`
	Errors       []ErrorResponse `json:"errors"`
}

// init runs before main function
func init() {
	environment := os.Getenv("ENVIRONMENT")
	if environment != "production" || environment == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// set PORT
	if os.Getenv("PORT") == "" {
		log.Fatal("$PORT must be set")
	}

	// set STRAVA_CLIENT_ID
	if os.Getenv("STRAVA_CLIENT_ID") == "" {
		log.Fatal("$STRAVA_CLIENT_ID must be set")
	}

	// set STRAVA_CLIENT_ID
	if os.Getenv("STRAVA_SECRET") == "" {
		log.Fatal("$STRAVA_SECRET must be set")
	}
}

// main function
func main() {
	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")
	engine.Reload(true)
	engine.Debug(true)

	app := fiber.New(&fiber.Settings{
		Views: engine,
	})

	// index
	app.Get("/", func(c *fiber.Ctx) {
		loginURI := fmt.Sprintf(
			"%s/authorize?client_id=%s&response_type=code&redirect_uri=%s&approval_prompt=force&scope=%s",
			StravaBase,
			os.Getenv("STRAVA_CLIENT_ID"),
			os.Getenv("STRAVA_REDIRECT_URI"),
			os.Getenv("STRAVA_SCOPE"),
		)

		// Render index template
		_ = c.Render("index", fiber.Map{
			"Title": "Login to Strava",
			"Login": loginURI,
		}, "layout/main")
	})

	// strava-oauth
	app.Get("/strava-oauth", func(c *fiber.Ctx) {
		var oathResp OauthResponse

		// build token url
		oauthURI := fmt.Sprintf(
			"%s/token?client_id=%s&client_secret=%s&code=%s&grant_type=%s",
			StravaBase,
			os.Getenv("STRAVA_CLIENT_ID"),
			os.Getenv("STRAVA_SECRET"),
			c.Query("code"),
			"authorization_code",
		)
		// send empty post request and params
		resp, err := http.Post(oauthURI, "application/json", nil)
		if err != nil {
			print(err)
		}
		defer resp.Body.Close()

		jsonData, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(jsonData, &oathResp)
		if err != nil {
			fmt.Println(err)
		}

		if len(oathResp.Errors) > 0 {
			// Render login_error template
			_ = c.Render("login_error.html", fiber.Map{
				"Title":    "Authorization Failure",
				"Code":     oathResp.Errors[0].Code,
				"Resource": oathResp.Errors[0].Resource,
				"Field":    oathResp.Errors[0].Field,
			}, "layout/main")
		} else {
			// Render login_results template
			_ = c.Render("login_results", fiber.Map{
				"Title":        "Authorization Success",
				"Athlete":      oathResp.Athlete,
				"ExpiresAt":    oathResp.ExpiresAt,
				"ExpiresIn":    oathResp.ExpriresIn,
				"AccessToken":  oathResp.AccessToken,
				"RefreshToken": oathResp.RefreshToken,
			}, "layout/main")
		}
	})

	app.Listen(os.Getenv("PORT"))
}
