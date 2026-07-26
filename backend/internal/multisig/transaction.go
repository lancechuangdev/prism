package multisig

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/lancechuangdev/prism/backend/internal/contracts"
)

var ErrInvalidProposal = errors.New("invalid multisig proposal")

const (
	OperationAddOwner        = "add_owner"
	OperationRemoveOwner     = "remove_owner"
	OperationReplaceOwner    = "replace_owner"
	OperationChangeThreshold = "change_threshold"
	OperationCreatePool      = "create_pool"
	OperationSettlePool      = "settle_pool"
	OperationRepayPool       = "repay_pool"
	OperationLiquidatePool   = "liquidate_pool"
)

type ConfigChangeParams struct {
	ChainID         string
	MultisigAddress string
	Operation       string
	Owner           string
	OldOwner        string
	NewOwner        string
	Threshold       uint64
	Nonce           string
}

type ProposalParams struct {
	ChainID         string
	MultisigAddress string
	Operation       string
	Target          string
	Value           string
	Data            string
	Nonce           string
}

type Proposal struct {
	TransactionHash string `json:"transactionHash"`
	Operation       string `json:"operation"`
	Target          string `json:"target"`
	Value           string `json:"value"`
	Data            string `json:"data"`
	Nonce           string `json:"nonce"`
}

type PreparedTransaction struct {
	To      string `json:"to"`
	Data    string `json:"data"`
	Value   string `json:"value"`
	ChainID string `json:"chainId"`
}

type PreparedProposal struct {
	Proposal             Proposal            `json:"proposal"`
	ApprovalTransaction  PreparedTransaction `json:"approvalTransaction"`
	ExecutionTransaction PreparedTransaction `json:"executionTransaction"`
}

type ProposalPreparer interface {
	PrepareConfigChange(ctx context.Context, params ConfigChangeParams) (PreparedProposal, error)
	PrepareProposal(ctx context.Context, params ProposalParams) (PreparedProposal, error)
}

type TransactionBuilder struct {
	chainID string
}

func NewTransactionBuilder(chainID string) (*TransactionBuilder, error) {
	parsedChainID, ok := new(big.Int).SetString(chainID, 10)
	if !ok || parsedChainID.Sign() <= 0 {
		return nil, fmt.Errorf("invalid chain ID %q", chainID)
	}
	return &TransactionBuilder{chainID: parsedChainID.String()}, nil
}

func (b *TransactionBuilder) PrepareConfigChange(ctx context.Context, params ConfigChangeParams) (PreparedProposal, error) {
	multisigABI, err := contracts.ThresholdMultiSigMetaData.GetAbi()
	if err != nil {
		return PreparedProposal{}, fmt.Errorf("load ThresholdMultiSig ABI: %w", err)
	}
	data, err := packConfigChange(multisigABI, params)
	if err != nil {
		return PreparedProposal{}, err
	}
	return b.PrepareProposal(ctx, ProposalParams{
		ChainID:         params.ChainID,
		MultisigAddress: params.MultisigAddress,
		Operation:       params.Operation,
		Target:          params.MultisigAddress,
		Value:           "0x0",
		Data:            hexutil.Encode(data),
		Nonce:           params.Nonce,
	})
}

func (b *TransactionBuilder) PrepareProposal(_ context.Context, params ProposalParams) (PreparedProposal, error) {
	if params.ChainID != b.chainID {
		return PreparedProposal{}, fmt.Errorf("%w: chain_id must be %s", ErrInvalidProposal, b.chainID)
	}
	multisigAddress, err := parseAddress("multisig address", params.MultisigAddress)
	if err != nil {
		return PreparedProposal{}, err
	}
	target, err := parseAddress("target", params.Target)
	if err != nil {
		return PreparedProposal{}, err
	}
	value, err := hexutil.DecodeBig(params.Value)
	if err != nil || value.Sign() < 0 {
		return PreparedProposal{}, fmt.Errorf("%w: value must be a non-negative hex quantity", ErrInvalidProposal)
	}
	data, err := hexutil.Decode(params.Data)
	if err != nil || len(data) < 4 {
		return PreparedProposal{}, fmt.Errorf("%w: data must be hex-encoded contract calldata", ErrInvalidProposal)
	}
	nonce, ok := new(big.Int).SetString(params.Nonce, 10)
	if !ok || nonce.Sign() < 0 {
		return PreparedProposal{}, fmt.Errorf("%w: nonce must be a non-negative decimal integer", ErrInvalidProposal)
	}

	multisigABI, err := contracts.ThresholdMultiSigMetaData.GetAbi()
	if err != nil {
		return PreparedProposal{}, fmt.Errorf("load ThresholdMultiSig ABI: %w", err)
	}
	approvalData, err := multisigABI.Pack("approveTransaction", target, value, data, nonce)
	if err != nil {
		return PreparedProposal{}, fmt.Errorf("encode approveTransaction: %w", err)
	}
	executionData, err := multisigABI.Pack("executeTransaction", target, value, data, nonce)
	if err != nil {
		return PreparedProposal{}, fmt.Errorf("encode executeTransaction: %w", err)
	}

	proposal := Proposal{
		Operation: params.Operation,
		Target:    target.Hex(),
		Value:     hexutil.EncodeBig(value),
		Data:      hexutil.Encode(data),
		Nonce:     nonce.String(),
	}
	outerTarget := multisigAddress.Hex()
	return PreparedProposal{
		Proposal: proposal,
		ApprovalTransaction: PreparedTransaction{
			To: outerTarget, Data: hexutil.Encode(approvalData), Value: "0x0", ChainID: b.chainID,
		},
		ExecutionTransaction: PreparedTransaction{
			To: outerTarget, Data: hexutil.Encode(executionData), Value: "0x0", ChainID: b.chainID,
		},
	}, nil
}

func packConfigChange(contractABI *abi.ABI, params ConfigChangeParams) ([]byte, error) {
	switch params.Operation {
	case OperationAddOwner:
		owner, err := parseAddress("owner", params.Owner)
		if err != nil {
			return nil, err
		}
		return packMethod(contractABI, "addOwner", owner)
	case OperationRemoveOwner:
		owner, err := parseAddress("owner", params.Owner)
		if err != nil {
			return nil, err
		}
		return packMethod(contractABI, "removeOwner", owner)
	case OperationReplaceOwner:
		oldOwner, err := parseAddress("old_owner", params.OldOwner)
		if err != nil {
			return nil, err
		}
		newOwner, err := parseAddress("new_owner", params.NewOwner)
		if err != nil {
			return nil, err
		}
		if oldOwner == newOwner {
			return nil, fmt.Errorf("%w: old_owner and new_owner must differ", ErrInvalidProposal)
		}
		return packMethod(contractABI, "replaceOwner", oldOwner, newOwner)
	case OperationChangeThreshold:
		if params.Threshold == 0 {
			return nil, fmt.Errorf("%w: threshold must be positive", ErrInvalidProposal)
		}
		return packMethod(contractABI, "changeThreshold", new(big.Int).SetUint64(params.Threshold))
	default:
		return nil, fmt.Errorf("%w: unsupported config operation %q", ErrInvalidProposal, params.Operation)
	}
}

func packMethod(contractABI *abi.ABI, method string, args ...any) ([]byte, error) {
	data, err := contractABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("encode %s: %w", method, err)
	}
	return data, nil
}

func parseAddress(name string, raw string) (common.Address, error) {
	if !common.IsHexAddress(raw) || common.HexToAddress(raw) == (common.Address{}) {
		return common.Address{}, fmt.Errorf("%w: %s must be a non-zero address", ErrInvalidProposal, name)
	}
	return common.HexToAddress(raw), nil
}
