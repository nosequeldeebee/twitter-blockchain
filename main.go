package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Tweet is a collection of important info in each Tweet
type Tweet struct {
	Date string `json:"created_at"`
	Text string `json:"text"`
	ID   string `json:"id_str"`
}

var config *oauth1.Config
var token *oauth1.Token

// API page limit
const pages = 18

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Load API credentials
	config = oauth1.NewConfig(os.Getenv("APIKEY"), os.Getenv("APISECRET"))
	token = oauth1.NewToken(os.Getenv("TOKEN"), os.Getenv("TOKENSECRET"))

	s := &http.Server{
		Addr:           os.Getenv("PORT"),
		Handler:        makeMuxRouter(),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// makeMuxRouter defines and creates routes
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/{id}", handleGetTweets).Methods("GET")
	return muxRouter
}

func respondWithError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}

func handleGetTweets(w http.ResponseWriter, r *http.Request) {
	var maxIDQuery string
	var tweets []Tweet
	vars := mux.Vars(r)
	userID := vars["id"]

	// httpClient will automatically authorize http.Requests
	httpClient := config.Client(oauth1.NoContext, token)

Outer:
	for i := 0; i < pages; i++ {
		// Twitter API request
		// userID is the Twitter handle
		// maxIDQuery is the last result on each page, so the API knows what the next page is
		path := fmt.Sprintf("https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name=%v&include_rts=false&count=200%v", userID, maxIDQuery)
		if strings.Contains(path, "favicon.ico") { // API returns this so skip it, not needed
			break
		}

		resp, err := httpClient.Get(path)
		if err != nil {
			respondWithError(err, w)
			break
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			respondWithError(err, w)
			break
		}

		var gotTweets []Tweet
		err = json.Unmarshal(body, &gotTweets)
		if err != nil {
			respondWithError(err, w)
			break
		}

		// range through Tweets to clear out unneeded info
		for i, t := range gotTweets {

			if i == len(gotTweets)-1 {
				// this is the logic to tell Twitter API where the pages are
				if maxIDQuery == fmt.Sprintf("&max_id=%v", t.ID) {
					break Outer
				}
				maxIDQuery = fmt.Sprintf("&max_id=%v", t.ID)
			}

			// remove @ mentions and links from returned Tweets
			regAt := regexp.MustCompile(`@(\S+)`)
			t.Text = regAt.ReplaceAllString(t.Text, "")
			regHTTP := regexp.MustCompile(`http(\S+)`)
			t.Text = regHTTP.ReplaceAllString(t.Text, "")
			tweets = append(tweets, t)
		}
	}

	var result []string

	for _, t := range tweets {
		result = append(result, t.Text)
	}

	stringResult := strings.Join(result, "\n")

	w.WriteHeader(200)
	w.Write([]byte(stringResult))
}
