package multisig

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/lancechuangdev/prism/backend/internal/contracts"
)

var ErrInvalidTransactionHash = errors.New("invalid multisig transaction hash")

type Config struct {
	ChainID         string   `json:"chain_id"`
	ContractAddress string   `json:"contract_address"`
	Owners          []string `json:"owners"`
	Threshold       uint64   `json:"threshold"`
}

type OwnerApproval struct {
	Address  string `json:"address"`
	Approved bool   `json:"approved"`
}

type ProposalStatus struct {
	TransactionHash              string          `json:"transactionHash"`
	CurrentConfigurationVersion  string          `json:"currentConfigurationVersion"`
	ProposalConfigurationVersion *string         `json:"proposalConfigurationVersion"`
	ConfigurationValid           bool            `json:"configurationValid"`
	ApprovalCount                uint64          `json:"approvalCount"`
	Threshold                    uint64          `json:"threshold"`
	ReadyToExecute               bool            `json:"readyToExecute"`
	Executed                     bool            `json:"executed"`
	Owners                       []OwnerApproval `json:"owners"`
}

type ChainReader interface {
	Config(ctx context.Context) (Config, error)
	TransactionHash(ctx context.Context, proposal Proposal) (string, error)
	ProposalStatus(ctx context.Context, transactionHash string) (ProposalStatus, error)
	Close()
}

type RPCReader struct {
	client   *ethclient.Client
	contract *contracts.ThresholdMultiSigCaller
	address  common.Address
	chainID  string
}

func NewRPCReader(ctx context.Context, rpcURL string, contractAddress string) (*RPCReader, error) {
	if strings.TrimSpace(rpcURL) == "" {
		return nil, fmt.Errorf("RPC URL is required")
	}
	if !validAddress(contractAddress) {
		return nil, fmt.Errorf("invalid ThresholdMultiSig address %q", contractAddress)
	}

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("connect chain RPC: %w", err)
	}
	address := common.HexToAddress(contractAddress)
	code, err := client.CodeAt(ctx, address, nil)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("read ThresholdMultiSig bytecode: %w", err)
	}
	if len(code) == 0 {
		client.Close()
		return nil, fmt.Errorf("no contract deployed at ThresholdMultiSig address %s", address.Hex())
	}
	chainID, err := client.ChainID(ctx)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("read RPC chain ID: %w", err)
	}
	contract, err := contracts.NewThresholdMultiSigCaller(address, client)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("bind ThresholdMultiSig: %w", err)
	}
	return &RPCReader{
		client: client, contract: contract, address: address, chainID: chainID.String(),
	}, nil
}

func (r *RPCReader) Close() {
	r.client.Close()
}

func (r *RPCReader) Config(ctx context.Context) (Config, error) {
	opts := &bind.CallOpts{Context: ctx}
	ownerCount, err := r.contract.OwnerCount(opts)
	if err != nil {
		return Config{}, fmt.Errorf("call ownerCount: %w", err)
	}
	if !ownerCount.IsInt64() {
		return Config{}, fmt.Errorf("owner count %s exceeds int64", ownerCount)
	}
	threshold, err := r.contract.Threshold(opts)
	if err != nil {
		return Config{}, fmt.Errorf("call threshold: %w", err)
	}
	if !threshold.IsUint64() {
		return Config{}, fmt.Errorf("threshold %s exceeds uint64", threshold)
	}

	owners := make([]string, ownerCount.Int64())
	for index := int64(0); index < ownerCount.Int64(); index++ {
		owner, err := r.contract.GetOwner(opts, big.NewInt(index))
		if err != nil {
			return Config{}, fmt.Errorf("call getOwner(%d): %w", index, err)
		}
		owners[index] = owner.Hex()
	}
	return Config{
		ChainID: r.chainID, ContractAddress: r.address.Hex(),
		Owners: owners, Threshold: threshold.Uint64(),
	}, nil
}

func (r *RPCReader) TransactionHash(ctx context.Context, proposal Proposal) (string, error) {
	target, err := parseAddress("target", proposal.Target)
	if err != nil {
		return "", err
	}
	value, err := hexutil.DecodeBig(proposal.Value)
	if err != nil || value.Sign() < 0 {
		return "", fmt.Errorf("%w: value must be a non-negative hex quantity", ErrInvalidProposal)
	}
	data, err := hexutil.Decode(proposal.Data)
	if err != nil || len(data) < 4 {
		return "", fmt.Errorf("%w: data must be hex-encoded contract calldata", ErrInvalidProposal)
	}
	nonce, ok := new(big.Int).SetString(proposal.Nonce, 10)
	if !ok || nonce.Sign() < 0 {
		return "", fmt.Errorf("%w: nonce must be a non-negative decimal integer", ErrInvalidProposal)
	}
	hash, err := r.contract.GetTransactionHash(
		&bind.CallOpts{Context: ctx}, target, value, data, nonce,
	)
	if err != nil {
		return "", fmt.Errorf("call getTransactionHash: %w", err)
	}
	return common.BytesToHash(hash[:]).Hex(), nil
}

func (r *RPCReader) ProposalStatus(ctx context.Context, transactionHash string) (ProposalStatus, error) {
	decodedHash, err := hexutil.Decode(transactionHash)
	if err != nil || len(decodedHash) != common.HashLength {
		return ProposalStatus{}, fmt.Errorf("%w: must be a 32-byte hex value", ErrInvalidTransactionHash)
	}
	hash := common.BytesToHash(decodedHash)
	opts := &bind.CallOpts{Context: ctx}

	approvalCount, err := r.contract.ApprovalCount(opts, hash)
	if err != nil {
		return ProposalStatus{}, fmt.Errorf("call approvalCount: %w", err)
	}
	configurationVersion, err := r.contract.ConfigurationVersion(opts)
	if err != nil {
		return ProposalStatus{}, fmt.Errorf("call configurationVersion: %w", err)
	}
	storedProposalVersion, err := r.contract.TransactionConfigurationVersion(opts, hash)
	if err != nil {
		return ProposalStatus{}, fmt.Errorf("call transactionConfigurationVersion: %w", err)
	}
	executed, err := r.contract.Executed(opts, hash)
	if err != nil {
		return ProposalStatus{}, fmt.Errorf("call executed: %w", err)
	}
	config, err := r.Config(ctx)
	if err != nil {
		return ProposalStatus{}, err
	}

	owners := make([]OwnerApproval, len(config.Owners))
	for index, rawOwner := range config.Owners {
		owner := common.HexToAddress(rawOwner)
		approved, err := r.contract.HasApproved(opts, hash, owner)
		if err != nil {
			return ProposalStatus{}, fmt.Errorf("call hasApproved for %s: %w", owner.Hex(), err)
		}
		owners[index] = OwnerApproval{Address: owner.Hex(), Approved: approved}
	}
	if !approvalCount.IsUint64() {
		return ProposalStatus{}, fmt.Errorf("approval count %s exceeds uint64", approvalCount)
	}
	count := approvalCount.Uint64()
	var proposalVersion *string
	configurationValid := false
	if storedProposalVersion.Sign() > 0 {
		decodedVersion := new(big.Int).Sub(storedProposalVersion, big.NewInt(1))
		value := decodedVersion.String()
		proposalVersion = &value
		configurationValid = decodedVersion.Cmp(configurationVersion) == 0
	}
	return ProposalStatus{
		TransactionHash: hash.Hex(), CurrentConfigurationVersion: configurationVersion.String(),
		ProposalConfigurationVersion: proposalVersion, ConfigurationValid: configurationValid,
		ApprovalCount: count, Threshold: config.Threshold,
		ReadyToExecute: count >= config.Threshold && configurationValid && !executed,
		Executed:       executed, Owners: owners,
	}, nil
}

func validAddress(raw string) bool {
	return common.IsHexAddress(raw) && common.HexToAddress(raw) != (common.Address{})
}
