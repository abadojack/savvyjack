package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/abadojack/anaconda"
)

type twitterAPI struct {
	api *anaconda.TwitterApi
}

//truncateString truncates string and adds 3 dots at the end.
func truncateString(s string, n int) string {
	if len(s) <= n {
		return s
	}

	n -= 4

	return s[:n] + "..."
}

// correctPeopleOnTwitter corrects people on Twitter :) by replying tweets which
// contain trackKey with one of the replies.
// e.g. You can use it to correct people who say 'Linux' instead of GNU/Linux.
// Ignore tweets that contain ignoreKey ... in the case above this would be GNU.
func (t *twitterAPI) correctPeopleOnTwitter(trackKey, ignoreKey string, replies []string) {
	v := url.Values{
		"track": []string{trackKey},
	}

	stream := t.api.PublicStreamFilter(v)
	defer stream.Stop()

	tweet := anaconda.Tweet{}

	self, err := t.api.GetSelf(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("bot account handle: ", self.ScreenName)

	for {
		log.Println("Waiting for relevant tweet")
		i := <-stream.C
		tweet = i.(anaconda.Tweet)

		//Do not reply to my own tweets
		if strings.Compare(tweet.User.ScreenName, self.ScreenName) == 0 {
			continue
		}

		//Ignore tweets that already contain ignoreKey
		if !strings.Contains(strings.ToLower(tweet.Text), ignoreKey) {
			rand.Seed(time.Now().Unix())
			replyStr := replies[rand.Intn(len(replies))] + " RT @" + strings.ToLower(tweet.User.ScreenName) + " " + tweet.Text

			log.Println("Posting reply: ", replyStr)

			_, err := t.api.PostTweet(truncateString(replyStr, 140), nil)
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println("Reply posted successfuly.")

				//Tweet only once every 2 minutes
				time.Sleep(2 * time.Minute)
			}
		}

	}
}

func main() {
	//Replies is a wrapper for slice Replies
	type Replies struct {
		Replies []string
	}

	b, err := ioutil.ReadFile("replies.json")
	if err != nil {
		panic(err)
	}

	var replies Replies
	err = json.Unmarshal(b, &replies)
	if err != nil {
		panic(err)
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY_SAVVY_JACK"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET_SAVVY_JACK"))
	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_KEY_SAVVY_JACK"), os.Getenv("TWITTER_ACCESS_SECRET_SAVVY_JACK"))

	t := twitterAPI{api}
	t.correctPeopleOnTwitter("Jack Sparrow", "captain", replies.Replies)
}
