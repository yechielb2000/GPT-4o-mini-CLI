package builders

import (
	"encoding/json"
	"gpt4omini/events"
)

func BuildEvent(rawEvent []byte) (events.Event, error) {
	var event = events.Event{}
	var err error
	err = json.Unmarshal(rawEvent, &event)
	return event, err
}
