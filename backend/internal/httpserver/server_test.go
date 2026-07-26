package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lancechuangdev/prism/backend/internal/auth"
	"github.com/lancechuangdev/prism/backend/internal/chain"
	"github.com/lancechuangdev/prism/backend/internal/config"
	"github.com/lancechuangdev/prism/backend/internal/multisig"
	"github.com/lancechuangdev/prism/backend/internal/price"
	"github.com/lancechuangdev/prism/backend/internal/store"
)

type testPoolTransactionPreparer struct {
	result       chain.PreparedTransaction
	err          error
	params       chain.CreatePoolParams
	settlePoolID string
	repayParams  chain.RepayPoolParams
}

type testMultisigTransactionPreparer struct {
	result         multisig.PreparedProposal
	err            error
	configParams   multisig.ConfigChangeParams
	proposalParams multisig.ProposalParams
}

type testMultisigReader struct {
	config multisig.Config
	status multisig.ProposalStatus
	hash   string
	err    error
}

func (r *testMultisigReader) Config(_ context.Context) (multisig.Config, error) {
	return r.config, r.err
}

func (r *testMultisigReader) ProposalStatus(_ context.Context, _ string) (multisig.ProposalStatus, error) {
	return r.status, r.err
}

func (r *testMultisigReader) TransactionHash(_ context.Context, _ multisig.Proposal) (string, error) {
	return r.hash, r.err
}

func (r *testMultisigReader) Close() {}

func (p *testMultisigTransactionPreparer) PrepareConfigChange(_ context.Context, params multisig.ConfigChangeParams) (multisig.PreparedProposal, error) {
	p.configParams = params
	return p.result, p.err
}

func (p *testMultisigTransactionPreparer) PrepareProposal(_ context.Context, params multisig.ProposalParams) (multisig.PreparedProposal, error) {
	p.proposalParams = params
	return p.result, p.err
}

func (c *testPoolTransactionPreparer) PrepareCreatePool(_ context.Context, params chain.CreatePoolParams) (chain.PreparedTransaction, error) {
	c.params = params
	return c.result, c.err
}

func (c *testPoolTransactionPreparer) PrepareSettlePool(_ context.Context, poolID string) (chain.PreparedTransaction, error) {
	c.settlePoolID = poolID
	return c.result, c.err
}

func (c *testPoolTransactionPreparer) PrepareRepayPool(_ context.Context, params chain.RepayPoolParams) (chain.PreparedTransaction, error) {
	c.repayParams = params
	return c.result, c.err
}

func TestHealthz(t *testing.T) {
	server := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Status != "ok" {
		t.Fatalf("expected healthy response, got %+v", body)
	}
}

func TestPoolBaseInfo(t *testing.T) {
	server := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/poolBaseInfo?chainId=31337", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body listResponse[store.PoolBase]
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(body.Data) != 1 {
		t.Fatalf("expected one pool, got %d", len(body.Data))
	}
	if body.Data[0].Data.Key.PoolID != 1 {
		t.Fatalf("expected pool id 1, got %+v", body.Data[0])
	}
}

func TestPoolDataInfoRequiresChainID(t *testing.T) {
	server := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/poolDataInfo", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestTokenList(t *testing.T) {
	server := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/token?chainId=31337", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body tokenListResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(body.Data) != 2 {
		t.Fatalf("expected two tokens, got %d", len(body.Data))
	}
}

func TestPrice(t *testing.T) {
	server := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/price?symbol=PRM", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body priceResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Data.Symbol != "PRM" || body.Data.Price != "0.0027" {
		t.Fatalf("unexpected price response: %+v", body.Data)
	}
}

func TestLoginAndProtectedSession(t *testing.T) {
	server := newTestServer(t)

	token := loginForTest(t, server)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/session", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body sessionResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Username != "admin" {
		t.Fatalf("expected admin username, got %s", body.Username)
	}
}

func TestProtectedSessionRejectsMissingToken(t *testing.T) {
	server := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/session", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestLogoutRevokesToken(t *testing.T) {
	server := newTestServer(t)
	token := loginForTest(t, server)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected logout status %d, got %d", http.StatusOK, rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/admin/session", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected revoked token status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func loginForTest(t *testing.T, server *http.Server) string {
	t.Helper()

	body := bytes.NewBufferString(`{"name":"admin","password":"password"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", body)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected login status %d, got %d", http.StatusOK, rec.Code)
	}

	var response loginResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if response.TokenID == "" {
		t.Fatal("expected login token")
	}
	return response.TokenID
}

func TestGetMultisig(t *testing.T) {
	server := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/multisig", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var response dataResponse[multisig.Config]
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Data.ContractAddress != "0x1000000000000000000000000000000000000001" ||
		len(response.Data.Owners) != 2 || response.Data.Threshold != 2 {
		t.Fatalf("unexpected multisig config: %+v", response.Data)
	}
}

func TestGetMultisigProposalStatus(t *testing.T) {
	reader := defaultTestMultisigReader()
	reader.status = multisig.ProposalStatus{
		TransactionHash: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		ApprovalCount:   2, Threshold: 2, ReadyToExecute: true,
	}
	server := newTestServerWithDependencies(t, &testPoolTransactionPreparer{}, &testMultisigTransactionPreparer{}, reader)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/multisig/proposals/0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var response dataResponse[multisig.ProposalStatus]
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !response.Data.ReadyToExecute || response.Data.ApprovalCount != 2 {
		t.Fatalf("unexpected proposal status: %+v", response.Data)
	}
}

func TestPrepareMultisigConfigChangeProposal(t *testing.T) {
	preparer := &testMultisigTransactionPreparer{
		result: multisig.PreparedProposal{
			Proposal: multisig.Proposal{
				Operation: multisig.OperationAddOwner,
				Target:    "0x1000000000000000000000000000000000000001",
				Value:     "0x0",
				Data:      "0x1234",
				Nonce:     "7",
			},
			ApprovalTransaction: multisig.PreparedTransaction{
				To: "0x1000000000000000000000000000000000000001", Data: "0x5678", Value: "0x0", ChainID: "31337",
			},
			ExecutionTransaction: multisig.PreparedTransaction{
				To: "0x1000000000000000000000000000000000000001", Data: "0x9abc", Value: "0x0", ChainID: "31337",
			},
		},
	}
	server := newTestServerWithPreparers(t, &testPoolTransactionPreparer{}, preparer)
	token := loginForTest(t, server)

	body := bytes.NewBufferString(`{
		"chain_id":"31337",
		"nonce":"7",
		"operation":{
			"type":"add_owner",
			"params":{
				"owner":"0x3000000000000000000000000000000000000003"
			}
		}
	}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/multisig/proposals", body)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if preparer.configParams.ChainID != "31337" ||
		preparer.configParams.MultisigAddress != "0x1000000000000000000000000000000000000001" ||
		preparer.configParams.Operation != multisig.OperationAddOwner ||
		preparer.configParams.Owner != "0x3000000000000000000000000000000000000003" ||
		preparer.configParams.Nonce != "7" {
		t.Fatalf("unexpected prepare params: %+v", preparer.configParams)
	}
	var response dataResponse[multisig.PreparedProposal]
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Data.ApprovalTransaction.Data != "0x5678" ||
		response.Data.ExecutionTransaction.Data != "0x9abc" {
		t.Fatalf("unexpected prepared transactions: %+v", response.Data)
	}
}

func TestPrepareMultisigProposalRejectsParamsFromAnotherOperation(t *testing.T) {
	server := newTestServer(t)
	token := loginForTest(t, server)
	body := bytes.NewBufferString(`{
		"chain_id":"31337",
		"nonce":"7",
		"operation":{
			"type":"add_owner",
			"params":{
				"owner":"0x3000000000000000000000000000000000000003",
				"settleTime":"2000000000"
			}
		}
	}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/multisig/proposals", body)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func TestPrepareMultisigCreatePoolProposal(t *testing.T) {
	poolCreator := &testPoolTransactionPreparer{result: chain.PreparedTransaction{
		To: "0x4000000000000000000000000000000000000004", Data: "0x12345678",
		Value: "0x0", ChainID: "31337",
	}}
	preparer := &testMultisigTransactionPreparer{result: multisig.PreparedProposal{
		Proposal: multisig.Proposal{Operation: multisig.OperationCreatePool},
	}}
	server := newTestServerWithPreparers(t, poolCreator, preparer)
	token := loginForTest(t, server)

	body := bytes.NewBufferString(`{
		"chain_id":"31337",
		"nonce":"8",
		"operation":{
			"type":"create_pool",
			"params":{
				"settleTime":"2000000000",
				"maturityTime":"2000600000",
				"interestRate":"1000000",
				"maxLendSupply":"1000000000000000000000",
				"collateralizationRatio":"200000000",
				"lendToken":"0x1000000000000000000000000000000000000001",
				"collateralToken":"0x2000000000000000000000000000000000000002",
				"lenderPositionToken":"0x3000000000000000000000000000000000000003",
				"borrowerPositionToken":"0x4000000000000000000000000000000000000004",
				"liquidateRate":"20000000"
			}
		}
	}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/multisig/proposals", body)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if poolCreator.params.MaxLendSupply != "1000000000000000000000" {
		t.Fatalf("unexpected pool params: %+v", poolCreator.params)
	}
	if preparer.proposalParams.Operation != multisig.OperationCreatePool ||
		preparer.proposalParams.MultisigAddress != "0x1000000000000000000000000000000000000001" ||
		preparer.proposalParams.Target != "0x4000000000000000000000000000000000000004" ||
		preparer.proposalParams.Data != "0x12345678" ||
		preparer.proposalParams.Nonce != "8" {
		t.Fatalf("unexpected proposal params: %+v", preparer.proposalParams)
	}
}

func TestPrepareMultisigSettlePoolProposal(t *testing.T) {
	poolTransactions := &testPoolTransactionPreparer{result: chain.PreparedTransaction{
		To: "0x4000000000000000000000000000000000000004", Data: "0x12345678",
		Value: "0x0", ChainID: "31337",
	}}
	preparer := &testMultisigTransactionPreparer{result: multisig.PreparedProposal{
		Proposal: multisig.Proposal{Operation: multisig.OperationSettlePool},
	}}
	server := newTestServerWithPreparers(t, poolTransactions, preparer)
	token := loginForTest(t, server)

	body := bytes.NewBufferString(`{
		"chain_id":"31337",
		"nonce":"9",
		"operation":{
			"type":"settle_pool",
			"params":{"poolId":"0"}
		}
	}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/multisig/proposals", body)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if poolTransactions.settlePoolID != "0" {
		t.Fatalf("unexpected pool ID: %q", poolTransactions.settlePoolID)
	}
	if preparer.proposalParams.Operation != multisig.OperationSettlePool ||
		preparer.proposalParams.MultisigAddress != "0x1000000000000000000000000000000000000001" ||
		preparer.proposalParams.Target != "0x4000000000000000000000000000000000000004" ||
		preparer.proposalParams.Data != "0x12345678" ||
		preparer.proposalParams.Nonce != "9" {
		t.Fatalf("unexpected proposal params: %+v", preparer.proposalParams)
	}
}

func TestPrepareMultisigRepayPoolProposal(t *testing.T) {
	poolTransactions := &testPoolTransactionPreparer{result: chain.PreparedTransaction{
		To: "0x4000000000000000000000000000000000000004", Data: "0x12345678",
		Value: "0x0", ChainID: "31337",
	}}
	preparer := &testMultisigTransactionPreparer{result: multisig.PreparedProposal{
		Proposal: multisig.Proposal{Operation: multisig.OperationRepayPool},
	}}
	server := newTestServerWithPreparers(t, poolTransactions, preparer)
	token := loginForTest(t, server)

	body := bytes.NewBufferString(`{
		"chain_id":"31337",
		"nonce":"10",
		"operation":{
			"type":"repay_pool",
			"params":{
				"poolId":"1",
				"maxCollateralAmount":"5000000000000000000"
			}
		}
	}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/multisig/proposals", body)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if poolTransactions.repayParams.PoolID != "1" ||
		poolTransactions.repayParams.MaxCollateralAmount != "5000000000000000000" {
		t.Fatalf("unexpected repayment params: %+v", poolTransactions.repayParams)
	}
	if preparer.proposalParams.Operation != multisig.OperationRepayPool ||
		preparer.proposalParams.MultisigAddress != "0x1000000000000000000000000000000000000001" ||
		preparer.proposalParams.Target != "0x4000000000000000000000000000000000000004" ||
		preparer.proposalParams.Data != "0x12345678" ||
		preparer.proposalParams.Nonce != "10" {
		t.Fatalf("unexpected proposal params: %+v", preparer.proposalParams)
	}
}

func TestPrepareMultisigProposalRequiresAuth(t *testing.T) {
	server := newTestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/multisig/proposals", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func newTestServer(t *testing.T) *http.Server {
	return newTestServerWithPoolCreator(t, &testPoolTransactionPreparer{})
}

func newTestServerWithPoolCreator(t *testing.T, poolCreator chain.PoolTransactionPreparer) *http.Server {
	return newTestServerWithPreparers(t, poolCreator, &testMultisigTransactionPreparer{})
}

func newTestServerWithPreparers(t *testing.T, poolCreator chain.PoolTransactionPreparer, multisigPreparer multisig.ProposalPreparer) *http.Server {
	return newTestServerWithDependencies(t, poolCreator, multisigPreparer, defaultTestMultisigReader())
}

func newTestServerWithDependencies(t *testing.T, poolCreator chain.PoolTransactionPreparer, multisigPreparer multisig.ProposalPreparer, multisigReader multisig.ChainReader) *http.Server {
	t.Helper()

	repo := store.NewMemoryStore()
	if err := chain.SyncPools(context.Background(), chain.NewFakeReader(), repo, "31337"); err != nil {
		t.Fatalf("sync demo contract data: %v", err)
	}

	auth := auth.NewService(auth.Config{
		AdminUsername: "admin",
		AdminPassword: "password",
		TokenSecret:   "test-secret",
		TokenTTL:      time.Hour,
	})
	return New(
		config.Config{Env: "test", Port: "0", APIVersion: "1"},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		chain.NewQueryService(repo),
		poolCreator,
		multisigPreparer,
		multisigReader,
		auth,
		price.NewService(price.NewDemoProvider()),
	)
}

func defaultTestMultisigReader() *testMultisigReader {
	return &testMultisigReader{config: multisig.Config{
		ChainID: "31337", ContractAddress: "0x1000000000000000000000000000000000000001",
		Owners: []string{
			"0x2000000000000000000000000000000000000002",
			"0x3000000000000000000000000000000000000003",
		},
		Threshold: 2,
	}, hash: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
}
