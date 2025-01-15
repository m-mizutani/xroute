package model

type SlackMessage struct {
	Emoji   string              `json:"emoji"`
	Icon    string              `json:"icon"`
	Channel string              `json:"channel"`
	Color   string              `json:"color"`
	Title   string              `json:"title"`
	Body    string              `json:"body"`
	Fields  []SlackMessageField `json:"fields"`
}

type SlackMessageField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Link  string `json:"link"`
}
