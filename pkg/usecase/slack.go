package usecase

import (
	"context"
	"fmt"

	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/transmith/pkg/domain/interfaces"
	"github.com/m-mizutani/transmith/pkg/domain/model"
	"github.com/m-mizutani/transmith/pkg/utils/logging"
	"github.com/slack-go/slack"
)

func transmitSlack(ctx context.Context, msg model.SlackMessage, client interfaces.Slack) error {
	logger := logging.Extract(ctx)
	logger.Debug("Transmit slack message", "message", msg)

	attachment := buildSlackMessage(msg)
	options := []slack.MsgOption{
		slack.MsgOptionAttachments(attachment),
	}

	if msg.Emoji != "" { // Emoji has higher priority than Icon
		options = append(options, slack.MsgOptionIconEmoji(msg.Emoji))
	} else if msg.Icon != "" {
		options = append(options, slack.MsgOptionIconURL(msg.Icon))
	}

	if _, _, err := client.PostMessageContext(ctx, msg.Channel, options...); err != nil {
		return goerr.Wrap(err, "failed to post slack message", goerr.V("message", msg))
	}

	return nil
}

var preservedColors = map[string]string{
	"info":    "#2EB67D",
	"warning": "#FFA500",
	"error":   "#FF0000",
}

func buildSlackMessage(msg model.SlackMessage) slack.Attachment {
	color := "#2EB67D"
	if msg.Color != "" {
		if preserved, ok := preservedColors[msg.Color]; ok {
			color = preserved
		} else {
			color = msg.Color
		}
	}

	var blockSet []slack.Block

	if msg.Title != "" {
		txt := slack.NewTextBlockObject("plain_text", msg.Title, false, false)
		blockSet = append(blockSet, slack.NewHeaderBlock(txt))
	}

	var body *slack.TextBlockObject
	if msg.Body != "" {
		body = slack.NewTextBlockObject("mrkdwn", msg.Body, false, false)
	}

	fields := make([]*slack.TextBlockObject, len(msg.Fields))
	for i, field := range msg.Fields {
		value := field.Value
		if field.Link != "" {
			value = "<" + field.Link + "|" + field.Value + ">"
		}
		mrkdwn := fmt.Sprintf("*%s*\n%s", field.Name, value)
		fields[i] = slack.NewTextBlockObject("mrkdwn", mrkdwn, false, false)
	}

	blockSet = append(blockSet, slack.NewSectionBlock(body, fields, nil))

	return slack.Attachment{
		Color: color,
		Blocks: slack.Blocks{
			BlockSet: blockSet,
		},
	}
}
