// Package sentry implements the sentry backend.
package sentry

import (
	"context"
	"fmt"
	"sync"

	"github.com/oasisprotocol/oasis-core/go/common/identity"
	"github.com/oasisprotocol/oasis-core/go/common/logging"
	consensus "github.com/oasisprotocol/oasis-core/go/consensus/api"
	"github.com/oasisprotocol/oasis-core/go/sentry/api"
)

var _ api.Backend = (*backend)(nil)

type backend struct {
	sync.RWMutex

	logger *logging.Logger

	consensus consensus.Service
	identity  *identity.Identity
}

func (b *backend) GetAddresses(context.Context) (*api.SentryAddresses, error) {
	// Consensus addresses.
	consensusAddrs, err := b.consensus.GetAddresses()
	if err != nil {
		return nil, fmt.Errorf("sentry: error obtaining consensus addresses: %w", err)
	}
	b.logger.Debug("successfully obtained consensus addresses",
		"addresses", consensusAddrs,
	)

	return &api.SentryAddresses{
		Consensus: consensusAddrs,
	}, nil
}

// New constructs a new sentry backend instance.
func New(
	consensus consensus.Service,
	identity *identity.Identity,
) (api.Backend, error) {
	if consensus == nil {
		return nil, fmt.Errorf("sentry: consensus backend is nil")
	}

	b := &backend{
		logger:    logging.GetLogger("sentry"),
		consensus: consensus,
		identity:  identity,
	}

	return b, nil
}
