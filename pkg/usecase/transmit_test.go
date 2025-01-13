package usecase_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/opac"
	"github.com/m-mizutani/xroute/pkg/adapter"
	"github.com/m-mizutani/xroute/pkg/domain/model"
	"github.com/m-mizutani/xroute/pkg/mock"
	"github.com/m-mizutani/xroute/pkg/usecase"
	"github.com/slack-go/slack"
)

//go:embed testdata/transmit.rego
var policyTransmitRego string

func TestTransmitSlack(t *testing.T) {
	slackMock := mock.SlackMock{
		PostMessageContextFunc: func(ctx context.Context, channelID string, options ...slack.MsgOption) (string, string, error) {
			return "", "", nil
		},
	}

	policy, err := opac.New(opac.Data(map[string]string{
		"transmit.rego": policyTransmitRego,
	}))
	gt.NoError(t, err)

	adapters := adapter.New(adapter.WithSlack(&slackMock), adapter.WithPolicy(policy))
	uc := usecase.New(adapters)

	msg := model.Message{
		Schema: "for_slack",
		Data:   "Hello, Slack",
	}
	gt.NoError(t, uc.Route(context.Background(), msg))
}
