package adapter

import "github.com/m-mizutani/xroute/pkg/domain/interfaces"

type Adapters struct {
	slack  interfaces.Slack
	policy interfaces.Policy
}

func New(options ...Option) *Adapters {
	adapter := &Adapters{}
	for _, opt := range options {
		opt(adapter)
	}
	return adapter
}

func (x *Adapters) Slack() interfaces.Slack {
	return x.slack
}
func (x *Adapters) Policy() interfaces.Policy {
	return x.policy
}

type Option func(*Adapters)

func WithSlack(slack interfaces.Slack) Option {
	return func(a *Adapters) {
		a.slack = slack
	}
}

func WithPolicy(policy interfaces.Policy) Option {
	return func(a *Adapters) {
		a.policy = policy
	}
}
