package model

type PolicyAuthzInput struct {
}

type PolicyAuthzOutput struct {
	Allow bool `json:"allow"`
}

type PolicyTransmitInput struct {
	Message
}

type PolicyTransmitOutput struct {
	Slack []SlackMessage `json:"slack"`
}
