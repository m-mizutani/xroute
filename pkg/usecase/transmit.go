package usecase

import (
	"context"

	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/transmith/pkg/domain/model"
	"github.com/m-mizutani/transmith/pkg/utils/logging"
)

func (x *UseCases) Transmit(ctx context.Context, msg model.Message) error {
	logger := logging.Extract(ctx)
	logger.Debug("Run usecase")
	eb := goerr.NewBuilder(goerr.V("message", msg))

	input := model.PolicyTransmitInput{
		Message: msg,
	}
	var output model.PolicyTransmitOutput

	if err := x.adaptors.Policy().Query(ctx, "data.transmit", input, &output); err != nil {
		return eb.Wrap(err, "Failed to query policy")
	}
	logger.Debug("Query result", "input", input, "output", output)

	for _, slackMsg := range output.Slack {
		if err := transmitSlack(ctx, slackMsg, x.adaptors.Slack()); err != nil {
			return eb.Wrap(err, "Failed to transmit slack message")
		}
	}

	return nil
}
