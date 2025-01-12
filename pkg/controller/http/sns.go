package http

/*
type snsMessage struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Subject          string `json:"Subject"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
	SubscribeURL     string `json:"SubscribeURL"`
}

func handleSNSSchema(r *http.Request, uc *usecase.UseCase) error {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		return goerr.New("Invalid HTTP method", goerr.V("method", r.Method))
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return goerr.Wrap(err, "Unable to read request body")
	}

	var message snsMessage
	err = json.Unmarshal(body, &message)
	if err != nil {
		return goerr.Wrap(err, "Invalid JSON", goerr.V("body", string(body)))
	}

	logging.Extract(ctx).Debug("Received SNS message", "message", message)

	// Here you can handle different SNS message types like subscription confirmation, notification, etc.
	switch message.Type {
	case "SubscriptionConfirmation":
		if err := confirmSnsSubscription(ctx, message.SubscribeURL); err != nil {
			return err
		}

	case "Notification":
		if err := uc.Transmit(ctx, message.Message); err != nil {
			return err
		}
	}

	return nil
}

func confirmSnsSubscription(ctx context.Context, subscribeURL string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, subscribeURL, nil)
	if err != nil {
		return goerr.Wrap(err, "Failed to send GET request to subscription URL")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return goerr.Wrap(err, "Failed to send GET request to subscription URL")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return goerr.New("Failed to confirm subscription", goerr.V("status_code", resp.StatusCode), goerr.V("body", string(body)))
	}

	return nil
}
*/
