package main

import (
	"log"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/abadojack/anaconda"
)

var replies = []string{
	"CAPTAIN Jack Sparrow. Savvy?",
	"CAPTAIN Jack Sparrow, if you please.",
	"Captain... CAPTAIN Jack Sparrow.",
	"Captain... Captain Jack Sparrow.",
	"Captain Jack Sparrow, if you please.",
	"I'm Captain Jack Sparrow. Savvy?",
	"Captain Jack Sparrow.",
	"There should be a 'Captain' in there somewhere.",
	"I'm Captain Jack Sparrow!",
	"That's not nice, you didn't call me 'Captain.'",
}

var api *anaconda.TwitterApi

func init() {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY_SAVVY_JACK"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET_SAVVY_JACK"))
	api = anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_KEY_SAVVY_JACK"), os.Getenv("TWITTER_ACCESS_SECRET_SAVVY_JACK"))
}

//truncateString truncates string and adds 3 dots at the end.
func truncateString(s string, n int) string {
	if len(s) <= n {
		return s
	}

	n -= 4
	for !utf8.ValidString(s[:n]) {
		n--
	}
	return s[:n] + "..."
}

func main() {
	v := url.Values{
		"track": []string{"Jack Sparrow"},
	}

	stream := api.PublicStreamFilter(v)
	defer stream.Stop()

	tweet := anaconda.Tweet{}

	self, err := api.GetSelf(nil)
	if err != nil {
		log.Println(err)
	}

	for {
		i := <-stream.C
		tweet = i.(anaconda.Tweet)
		//Do not reply to my own tweets
		if strings.Compare(tweet.User.ScreenName, self.ScreenName) == 0 {
			continue
		}

		//Ignore tweets that already contain 'captain'
		if !strings.Contains(strings.ToLower(tweet.Text), "captain") {
			rand.Seed(time.Now().Unix())
			replyStr := replies[rand.Intn(len(replies))] + " RT @" + strings.ToLower(tweet.User.ScreenName) + " " + tweet.Text

			t, err := api.PostTweet(truncateString(replyStr, 140), nil)
			if err != nil {
				log.Println(t, err)
			}
		}
	}
}
