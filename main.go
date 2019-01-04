package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

var (
	address = ":8183"

	// Get Verification Token in Basic Information profile
	verificationToken = "aaaaaaaa"

	// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"
	accessToken = "xoxb-xxxxxx"
)

// create global slack client
var api = slack.New(accessToken)

func main() {
	http.HandleFunc("/events-endpoint", EventsHandler)

	log.Printf("Listen and serve %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

//  EventsHandler
func EventsHandler(w http.ResponseWriter, r *http.Request) {

	// create event object
	raw, _ := ioutil.ReadAll(r.Body)
	opts := slackevents.OptionVerifyToken(&slackevents.TokenComparator{
		VerificationToken: verificationToken,
	})
	eventsAPIEvent, err := slackevents.ParseEvent(raw, opts)
	if err != nil {
		log.Printf("slackevents.ParseEvent error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(raw, &r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent

		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			messageEvent := innerEvent.Data.(*slackevents.MessageEvent)
			reply := "may I help you?"

			if messageEvent.SubType == "" { // check if the message not from the bot itself
				api.PostMessage(ev.Channel, slack.MsgOptionText(reply, false)) // post the reply
			}
		}
	}
}
