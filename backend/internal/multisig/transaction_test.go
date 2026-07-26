package multisig

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/lancechuangdev/prism/backend/internal/contracts"
)

const (
	testMultisigAddress = "0x1000000000000000000000000000000000000001"
	testOwnerAddress    = "0x2000000000000000000000000000000000000002"
	testOldOwnerAddress = "0x3000000000000000000000000000000000000003"
	testNewOwnerAddress = "0x4000000000000000000000000000000000000004"
)

func TestPrepareConfigChange(t *testing.T) {
	builder, err := NewTransactionBuilder("31337")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		params ConfigChangeParams
		method string
	}{
		{
			name: "add owner",
			params: ConfigChangeParams{
				ChainID: "31337", MultisigAddress: testMultisigAddress, Operation: OperationAddOwner,
				Owner: testOwnerAddress, Nonce: "1",
			},
			method: "addOwner",
		},
		{
			name: "remove owner",
			params: ConfigChangeParams{
				ChainID: "31337", MultisigAddress: testMultisigAddress, Operation: OperationRemoveOwner,
				Owner: testOwnerAddress, Nonce: "2",
			},
			method: "removeOwner",
		},
		{
			name: "replace owner",
			params: ConfigChangeParams{
				ChainID: "31337", MultisigAddress: testMultisigAddress, Operation: OperationReplaceOwner,
				OldOwner: testOldOwnerAddress, NewOwner: testNewOwnerAddress, Nonce: "3",
			},
			method: "replaceOwner",
		},
		{
			name: "change threshold",
			params: ConfigChangeParams{
				ChainID: "31337", MultisigAddress: testMultisigAddress, Operation: OperationChangeThreshold,
				Threshold: 2, Nonce: "4",
			},
			method: "changeThreshold",
		},
	}

	contractABI, err := contracts.ThresholdMultiSigMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := builder.PrepareConfigChange(context.Background(), tt.params)
			if err != nil {
				t.Fatalf("prepare transactions: %v", err)
			}
			if result.Proposal.Target != common.HexToAddress(testMultisigAddress).Hex() ||
				result.Proposal.Operation != tt.params.Operation ||
				result.Proposal.Nonce != tt.params.Nonce {
				t.Fatalf("unexpected proposal: %+v", result.Proposal)
			}
			adminData, err := hexutil.Decode(result.Proposal.Data)
			if err != nil {
				t.Fatal(err)
			}
			if len(adminData) < 4 || string(adminData[:4]) != string(contractABI.Methods[tt.method].ID) {
				t.Fatalf("proposal does not call %s", tt.method)
			}
			assertOuterTransaction(t, contractABI, result.ApprovalTransaction, "approveTransaction", adminData, tt.params.Nonce)
			assertOuterTransaction(t, contractABI, result.ExecutionTransaction, "executeTransaction", adminData, tt.params.Nonce)
		})
	}
}

func TestPrepareConfigChangeRejectsInvalidInput(t *testing.T) {
	builder, err := NewTransactionBuilder("31337")
	if err != nil {
		t.Fatal(err)
	}

	tests := []ConfigChangeParams{
		{ChainID: "1", MultisigAddress: testMultisigAddress, Operation: OperationAddOwner, Owner: testOwnerAddress, Nonce: "1"},
		{ChainID: "31337", MultisigAddress: "invalid", Operation: OperationAddOwner, Owner: testOwnerAddress, Nonce: "1"},
		{ChainID: "31337", MultisigAddress: testMultisigAddress, Operation: OperationAddOwner, Owner: "invalid", Nonce: "1"},
		{ChainID: "31337", MultisigAddress: testMultisigAddress, Operation: OperationReplaceOwner, OldOwner: testOwnerAddress, NewOwner: testOwnerAddress, Nonce: "1"},
		{ChainID: "31337", MultisigAddress: testMultisigAddress, Operation: OperationChangeThreshold, Threshold: 0, Nonce: "1"},
		{ChainID: "31337", MultisigAddress: testMultisigAddress, Operation: "unknown", Nonce: "1"},
		{ChainID: "31337", MultisigAddress: testMultisigAddress, Operation: OperationAddOwner, Owner: testOwnerAddress, Nonce: "-1"},
	}

	for _, params := range tests {
		_, err := builder.PrepareConfigChange(context.Background(), params)
		if !errors.Is(err, ErrInvalidProposal) {
			t.Fatalf("expected invalid transaction error for %+v, got %v", params, err)
		}
	}
}

func TestPrepareProposalWrapsExternalTarget(t *testing.T) {
	builder, err := NewTransactionBuilder("31337")
	if err != nil {
		t.Fatal(err)
	}
	result, err := builder.PrepareProposal(context.Background(), ProposalParams{
		ChainID: "31337", MultisigAddress: testMultisigAddress,
		Operation: OperationCreatePool,
		Target:    testOwnerAddress, Value: "0x0", Data: "0x12345678", Nonce: "9",
	})
	if err != nil {
		t.Fatalf("prepare proposal: %v", err)
	}
	if result.Proposal.Target != common.HexToAddress(testOwnerAddress).Hex() ||
		result.ApprovalTransaction.To != common.HexToAddress(testMultisigAddress).Hex() ||
		result.Proposal.Operation != OperationCreatePool {
		t.Fatalf("unexpected external proposal: %+v", result)
	}
}

func assertOuterTransaction(t *testing.T, contractABI *abi.ABI, transaction PreparedTransaction, methodName string, adminData []byte, nonce string) {
	t.Helper()
	if transaction.To != common.HexToAddress(testMultisigAddress).Hex() ||
		transaction.Value != "0x0" ||
		transaction.ChainID != "31337" {
		t.Fatalf("unexpected %s transaction: %+v", methodName, transaction)
	}
	data, err := hexutil.Decode(transaction.Data)
	if err != nil {
		t.Fatal(err)
	}
	method := contractABI.Methods[methodName]
	if len(data) < 4 || string(data[:4]) != string(method.ID) {
		t.Fatalf("transaction does not call %s", methodName)
	}
	values, err := method.Inputs.Unpack(data[4:])
	if err != nil {
		t.Fatalf("decode %s: %v", methodName, err)
	}
	if values[0].(common.Address) != common.HexToAddress(testMultisigAddress) ||
		values[1].(*big.Int).Sign() != 0 ||
		string(values[2].([]byte)) != string(adminData) ||
		values[3].(*big.Int).String() != nonce {
		t.Fatalf("unexpected %s arguments: %+v", methodName, values)
	}
}
