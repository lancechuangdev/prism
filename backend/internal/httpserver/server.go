package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/lancechuangdev/prism/backend/internal/auth"
	"github.com/lancechuangdev/prism/backend/internal/chain"
	"github.com/lancechuangdev/prism/backend/internal/config"
	"github.com/lancechuangdev/prism/backend/internal/multisig"
	"github.com/lancechuangdev/prism/backend/internal/price"
	"github.com/lancechuangdev/prism/backend/internal/store"
)

type healthResponse struct {
	Status string `json:"status"`
	App    string `json:"app"`
}

type listResponse[T any] struct {
	Data []indexedItem[T] `json:"data"`
}

type indexedItem[T any] struct {
	Index int `json:"index"`
	Data  T   `json:"pool_data"`
}

type tokenListResponse struct {
	Data []store.TokenInfo `json:"data"`
}

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type loginResponse struct {
	TokenID string `json:"tokenId"`
}

type sessionResponse struct {
	Username string `json:"username"`
}

type priceResponse struct {
	Data price.Quote `json:"data"`
}

type dataResponse[T any] struct {
	Data T `json:"data"`
}

type prepareMultisigProposalRequest struct {
	ChainID   string                   `json:"chain_id"`
	Nonce     string                   `json:"nonce"`
	Operation multisigOperationRequest `json:"operation"`
}

type multisigOperationRequest struct {
	Type                   string `json:"type"`
	Owner                  string `json:"owner"`
	OldOwner               string `json:"old_owner"`
	NewOwner               string `json:"new_owner"`
	Threshold              uint64 `json:"threshold"`
	SettleTime             string `json:"settleTime"`
	MaturityTime           string `json:"maturityTime"`
	InterestRate           string `json:"interestRate"`
	MaxLendSupply          string `json:"maxLendSupply"`
	CollateralizationRatio string `json:"collateralizationRatio"`
	LendToken              string `json:"lendToken"`
	CollateralToken        string `json:"collateralToken"`
	LenderPositionToken    string `json:"lenderPositionToken"`
	BorrowerPositionToken  string `json:"borrowerPositionToken"`
	LiquidateRate          string `json:"liquidateRate"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func New(cfg config.Config,
	logger *slog.Logger,
	chainQueryService *chain.QueryService,
	poolTransactions chain.PoolTransactionPreparer,
	multisigTransactions multisig.ProposalPreparer,
	multisigReader multisig.ChainReader,
	authService *auth.Service,
	priceService *price.Service) *http.Server {
	mux := http.NewServeMux()
	apiPrefix := "/api/v" + strings.TrimPrefix(cfg.APIVersion, "v")

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{
			Status: "ok",
			App:    "prism-backend",
		})
	})

	mux.HandleFunc("GET "+apiPrefix+"/poolBaseInfo", func(w http.ResponseWriter, r *http.Request) {
		chainID, ok := requireChainID(w, r)
		if !ok {
			return
		}

		pools, err := chainQueryService.ListPoolBases(r.Context(), chainID)
		if err != nil {
			logger.Error("list pool base info failed", slog.Any("error", err))
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "list pool base info failed"})
			return
		}

		items := make([]indexedItem[store.PoolBase], len(pools))
		for i, pool := range pools {
			items[i] = indexedItem[store.PoolBase]{
				Index: int(pool.Key.PoolID - 1),
				Data:  pool,
			}
		}
		writeJSON(w, http.StatusOK, listResponse[store.PoolBase]{Data: items})
	})

	mux.HandleFunc("GET "+apiPrefix+"/poolDataInfo", func(w http.ResponseWriter, r *http.Request) {
		chainID, ok := requireChainID(w, r)
		if !ok {
			return
		}

		pools, err := chainQueryService.ListPoolData(r.Context(), chainID)
		if err != nil {
			logger.Error("list pool data info failed", slog.Any("error", err))
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "list pool data info failed"})
			return
		}

		items := make([]indexedItem[store.PoolData], 0, len(pools))
		for _, pool := range pools {
			items = append(items, indexedItem[store.PoolData]{
				Index: int(pool.Key.PoolID - 1),
				Data:  pool,
			})
		}
		writeJSON(w, http.StatusOK, listResponse[store.PoolData]{Data: items})
	})

	mux.HandleFunc("GET "+apiPrefix+"/token", func(w http.ResponseWriter, r *http.Request) {
		chainID, ok := requireChainID(w, r)
		if !ok {
			return
		}

		tokens, err := chainQueryService.ListTokens(r.Context(), chainID)
		if err != nil {
			logger.Error("list tokens failed", slog.Any("error", err))
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "list tokens failed"})
			return
		}
		writeJSON(w, http.StatusOK, tokenListResponse{Data: tokens})
	})

	mux.HandleFunc("POST "+apiPrefix+"/user/login", func(w http.ResponseWriter, r *http.Request) {
		req := loginRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid login body"})
			return
		}

		token, err := authService.Login(req.Name, req.Password)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid username or password"})
			return
		}

		writeJSON(w, http.StatusOK, loginResponse{TokenID: token})
	})

	mux.Handle("POST "+apiPrefix+"/user/logout", requireAuth(authService, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromRequest(r)
		authService.Logout(token)
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})))

	mux.Handle("GET "+apiPrefix+"/admin/session", requireAuth(authService, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, ok := usernameFromContext(r.Context())
		if !ok {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		} else {
			writeJSON(w, http.StatusOK, sessionResponse{Username: username})
		}
	})))

	mux.HandleFunc("GET "+apiPrefix+"/price", func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("symbol")
		if symbol == "" {
			symbol = cfg.PriceSymbol
		}

		quote, err := priceService.Latest(r.Context(), symbol)
		if err != nil {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "price not found"})
			return
		}

		writeJSON(w, http.StatusOK, priceResponse{Data: quote})
	})

	mux.HandleFunc("GET "+apiPrefix+"/multisig", func(w http.ResponseWriter, r *http.Request) {
		result, err := multisigReader.Config(r.Context())
		if err != nil {
			logger.Error("read multisig config failed", slog.Any("error", err))
			writeJSON(w, http.StatusBadGateway, errorResponse{Error: "read multisig config failed"})
			return
		}
		writeJSON(w, http.StatusOK, dataResponse[multisig.Config]{Data: result})
	})

	mux.HandleFunc("GET "+apiPrefix+"/multisig/proposals/{txHash}", func(w http.ResponseWriter, r *http.Request) {
		result, err := multisigReader.ProposalStatus(r.Context(), r.PathValue("txHash"))
		if err != nil {
			if errors.Is(err, multisig.ErrInvalidTransactionHash) {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
				return
			}
			logger.Error("read multisig proposal status failed", slog.Any("error", err))
			writeJSON(w, http.StatusBadGateway, errorResponse{Error: "read multisig proposal status failed"})
			return
		}
		writeJSON(w, http.StatusOK, dataResponse[multisig.ProposalStatus]{Data: result})
	})

	mux.Handle("POST "+apiPrefix+"/multisig/proposals", requireAuth(authService, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := prepareMultisigProposalRequest{}
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid multisig proposal body"})
			return
		}

		multisigConfig, err := multisigReader.Config(r.Context())
		if err != nil {
			logger.Error("read multisig config failed", slog.Any("error", err))
			writeJSON(w, http.StatusBadGateway, errorResponse{Error: "read multisig config failed"})
			return
		}
		if req.ChainID != multisigConfig.ChainID {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "chain_id does not match the configured multisig"})
			return
		}

		var result multisig.PreparedProposal
		switch req.Operation.Type {
		case multisig.OperationAddOwner, multisig.OperationRemoveOwner,
			multisig.OperationReplaceOwner, multisig.OperationChangeThreshold:
			result, err = multisigTransactions.PrepareConfigChange(r.Context(), multisig.ConfigChangeParams{
				ChainID:         multisigConfig.ChainID,
				MultisigAddress: multisigConfig.ContractAddress,
				Operation:       req.Operation.Type,
				Owner:           req.Operation.Owner,
				OldOwner:        req.Operation.OldOwner,
				NewOwner:        req.Operation.NewOwner,
				Threshold:       req.Operation.Threshold,
				Nonce:           req.Nonce,
			})
		case multisig.OperationCreatePool:
			var poolTransaction chain.PreparedTransaction
			poolTransaction, err = poolTransactions.PrepareCreatePool(r.Context(), chain.CreatePoolParams{
				SettleTime: req.Operation.SettleTime, MaturityTime: req.Operation.MaturityTime,
				InterestRate: req.Operation.InterestRate, MaxLendSupply: req.Operation.MaxLendSupply,
				CollateralizationRatio: req.Operation.CollateralizationRatio,
				LendToken:              req.Operation.LendToken, CollateralToken: req.Operation.CollateralToken,
				LenderPositionToken:   req.Operation.LenderPositionToken,
				BorrowerPositionToken: req.Operation.BorrowerPositionToken,
				LiquidateRate:         req.Operation.LiquidateRate,
			})
			if err == nil {
				result, err = multisigTransactions.PrepareProposal(r.Context(), multisig.ProposalParams{
					ChainID: multisigConfig.ChainID, MultisigAddress: multisigConfig.ContractAddress,
					Operation: multisig.OperationCreatePool,
					Target:    poolTransaction.To, Value: poolTransaction.Value,
					Data: poolTransaction.Data, Nonce: req.Nonce,
				})
			}
		default:
			err = fmt.Errorf("%w: unsupported operation %q", multisig.ErrInvalidProposal, req.Operation.Type)
		}
		if err != nil {
			if errors.Is(err, multisig.ErrInvalidProposal) || errors.Is(err, chain.ErrInvalidCreatePool) {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
				return
			}
			logger.Error("prepare multisig proposal failed", slog.Any("error", err))
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "prepare multisig proposal failed"})
			return
		}

		transactionHash, err := multisigReader.TransactionHash(r.Context(), result.Proposal)
		if err != nil {
			logger.Error("read multisig proposal hash failed", slog.Any("error", err))
			writeJSON(w, http.StatusBadGateway, errorResponse{Error: "read multisig proposal hash failed"})
			return
		}
		result.Proposal.TransactionHash = transactionHash
		writeJSON(w, http.StatusOK, dataResponse[multisig.PreparedProposal]{Data: result})
	})))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "route not found"})
	})

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: requestLogger(logger, mux),
	}
}

func requireChainID(w http.ResponseWriter, r *http.Request) (string, bool) {
	chainID := r.URL.Query().Get("chainId")
	if chainID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "chainId is required"})
		return "", false
	}
	if _, err := strconv.ParseInt(chainID, 10, 64); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "chainId must be a number"})
		return "", false
	}
	return chainID, true
}

func requireAuth(authService *auth.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromRequest(r)
		username, err := authService.Authenticate(token)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid token"})
			return
		}

		ctx := contextWithUsername(r, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func tokenFromRequest(r *http.Request) string {
	authorization := strings.TrimSpace(r.Header.Get("Authorization"))
	return strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
}

type usernameContextKey struct{}

func contextWithUsername(r *http.Request, username string) context.Context {
	return context.WithValue(r.Context(), usernameContextKey{}, username)
}

func usernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameContextKey{}).(string)
	return username, ok
}

func requestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("http request", slog.String("method", r.Method), slog.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
