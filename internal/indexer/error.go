package indexer

import "errors"

var (
	// ErrChainNotSupported is used to indicate that chain is not supported
	ErrChainNotSupported = errors.New("chain not supported")
)
