package multisig

import (
	"context"
	"errors"
	"testing"
)

func TestProposalStatusRejectsInvalidHash(t *testing.T) {
	reader := &RPCReader{}

	_, err := reader.ProposalStatus(context.Background(), "0x1234")

	if !errors.Is(err, ErrInvalidTransactionHash) {
		t.Fatalf("expected invalid transaction hash, got %v", err)
	}
}

func TestTransactionHashRejectsInvalidProposal(t *testing.T) {
	reader := &RPCReader{}

	_, err := reader.TransactionHash(context.Background(), Proposal{
		Target: "invalid", Value: "0x0", Data: "0x12345678", Nonce: "1",
	})

	if !errors.Is(err, ErrInvalidProposal) {
		t.Fatalf("expected invalid proposal, got %v", err)
	}
}
