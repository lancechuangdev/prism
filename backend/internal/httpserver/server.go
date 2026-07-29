package httpserver

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lancechuangdev/prism/backend/internal/auth"
	"github.com/lancechuangdev/prism/backend/internal/chain"
	"github.com/lancechuangdev/prism/backend/internal/config"
	"github.com/lancechuangdev/prism/backend/internal/multisig"
	"github.com/lancechuangdev/prism/backend/internal/price"
	"github.com/lancechuangdev/prism/backend/internal/readiness"
	"github.com/lancechuangdev/prism/backend/internal/store"
)

const (
	httpReadHeaderTimeout = 5 * time.Second
	httpReadTimeout       = 15 * time.Second
	httpWriteTimeout      = 30 * time.Second
	httpIdleTimeout       = 60 * time.Second
	httpMaxHeaderBytes    = 1 << 20
	httpRequestDeadline   = 25 * time.Second
	loginBodyLimit        = 4 << 10
	proposalBodyLimit     = 64 << 10
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
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params"`
}

type ownerOperationParams struct {
	Owner string `json:"owner"`
}

type replaceOwnerParams struct {
	OldOwner string `json:"old_owner"`
	NewOwner string `json:"new_owner"`
}

type changeThresholdParams struct {
	Threshold uint64 `json:"threshold"`
}

type createPoolOperationParams struct {
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

type settlePoolOperationParams struct {
	PoolID string `json:"poolId"`
}

type repayPoolOperationParams struct {
	PoolID              string `json:"poolId"`
	MaxCollateralAmount string `json:"maxCollateralAmount"`
}

type liquidatePoolOperationParams struct {
	PoolID              string `json:"poolId"`
	MaxCollateralAmount string `json:"maxCollateralAmount"`
}

type errorResponse struct {
	Error     string `json:"error"`
	Code      string `json:"code,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func New(cfg config.Config,
	logger *slog.Logger,
	chainQueryService *chain.QueryService,
	poolTransactions chain.PoolTransactionPreparer,
	multisigTransactions multisig.ProposalPreparer,
	multisigReader multisig.ChainReader,
	localAuth *auth.LocalAuthenticator,
	authorizer auth.Authorizer,
	priceService *price.Service,
	readinessChecker readiness.Checker) *http.Server {
	mux := http.NewServeMux()
	apiPrefix := "/api/v" + strings.TrimPrefix(cfg.APIVersion, "v")

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{
			Status: "ok",
			App:    "prism-backend",
		})
	})

	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		report := readinessChecker.Check(r.Context())
		status := http.StatusOK
		if report.Status != readiness.StatusReady {
			status = http.StatusServiceUnavailable
			logRequest(logger, r, slog.LevelWarn, "readiness check failed", slog.Any("dependencies", report.Dependencies))
		}
		writeJSON(w, status, report)
	})

	mux.HandleFunc("GET "+apiPrefix+"/poolBaseInfo", func(w http.ResponseWriter, r *http.Request) {
		chainID, ok := requireChainID(w, r)
		if !ok {
			return
		}

		pools, err := chainQueryService.ListPoolBases(r.Context(), chainID)
		if err != nil {
			logRequest(logger, r, slog.LevelError, "list pool base info failed", slog.Any("error", err))
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
			logRequest(logger, r, slog.LevelError, "list pool data info failed", slog.Any("error", err))
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
			logRequest(logger, r, slog.LevelError, "list tokens failed", slog.Any("error", err))
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "list tokens failed"})
			return
		}
		writeJSON(w, http.StatusOK, tokenListResponse{Data: tokens})
	})

	if cfg.AuthMode != "cognito" {
		mux.HandleFunc("POST "+apiPrefix+"/user/login", func(w http.ResponseWriter, r *http.Request) {
			req := loginRequest{}
			if err := decodeJSONBody(w, r, &req, loginBodyLimit); err != nil {
				writeJSONBodyError(w, err, "invalid login body")
				return
			}

			token, err := localAuth.Login(r.Context(), req.Name, req.Password)
			if err != nil {
				if errors.Is(err, auth.ErrInvalidCredentials) {
					writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid username or password"})
				} else {
					logRequest(logger, r, slog.LevelError, "create login session failed", slog.Any("error", err))
					writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "login failed"})
				}
				return
			}

			writeJSON(w, http.StatusOK, loginResponse{TokenID: token})
		})

		mux.Handle("POST "+apiPrefix+"/user/logout", requireAuth(authorizer, "", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := tokenFromRequest(r)
			if err := localAuth.Logout(r.Context(), token); err != nil {
				logRequest(logger, r, slog.LevelError, "delete login session failed", slog.Any("error", err))
				writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "logout failed"})
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})))
	}

	mux.Handle("GET "+apiPrefix+"/admin/session", requireAuth(authorizer, cfg.AdminReadScope, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			logRequest(logger, r, slog.LevelError, "read multisig config failed", slog.Any("error", err))
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
			logRequest(logger, r, slog.LevelError, "read multisig proposal status failed", slog.Any("error", err))
			writeJSON(w, http.StatusBadGateway, errorResponse{Error: "read multisig proposal status failed"})
			return
		}
		writeJSON(w, http.StatusOK, dataResponse[multisig.ProposalStatus]{Data: result})
	})

	mux.Handle("POST "+apiPrefix+"/multisig/proposals", requireAuth(authorizer, cfg.ProposalWriteScope, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := prepareMultisigProposalRequest{}
		if err := decodeJSONBody(w, r, &req, proposalBodyLimit); err != nil {
			writeJSONBodyError(w, err, "invalid multisig proposal body")
			return
		}

		multisigConfig, err := multisigReader.Config(r.Context())
		if err != nil {
			logRequest(logger, r, slog.LevelError, "read multisig config failed", slog.Any("error", err))
			writeJSON(w, http.StatusBadGateway, errorResponse{Error: "read multisig config failed"})
			return
		}
		if req.ChainID != multisigConfig.ChainID {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "chain_id does not match the configured multisig"})
			return
		}

		var result multisig.PreparedProposal
		switch req.Operation.Type {
		case multisig.OperationAddOwner, multisig.OperationRemoveOwner:
			params, decodeErr := decodeOperationParams[ownerOperationParams](req.Operation.Params)
			if decodeErr != nil {
				err = decodeErr
				break
			}
			result, err = multisigTransactions.PrepareConfigChange(r.Context(), multisig.ConfigChangeParams{
				ChainID:         multisigConfig.ChainID,
				MultisigAddress: multisigConfig.ContractAddress,
				Operation:       req.Operation.Type,
				Owner:           params.Owner,
				Nonce:           req.Nonce,
			})
		case multisig.OperationReplaceOwner:
			params, decodeErr := decodeOperationParams[replaceOwnerParams](req.Operation.Params)
			if decodeErr != nil {
				err = decodeErr
				break
			}
			result, err = multisigTransactions.PrepareConfigChange(r.Context(), multisig.ConfigChangeParams{
				ChainID: multisigConfig.ChainID, MultisigAddress: multisigConfig.ContractAddress,
				Operation: req.Operation.Type, OldOwner: params.OldOwner,
				NewOwner: params.NewOwner, Nonce: req.Nonce,
			})
		case multisig.OperationChangeThreshold:
			params, decodeErr := decodeOperationParams[changeThresholdParams](req.Operation.Params)
			if decodeErr != nil {
				err = decodeErr
				break
			}
			result, err = multisigTransactions.PrepareConfigChange(r.Context(), multisig.ConfigChangeParams{
				ChainID: multisigConfig.ChainID, MultisigAddress: multisigConfig.ContractAddress,
				Operation: req.Operation.Type, Threshold: params.Threshold, Nonce: req.Nonce,
			})
		case multisig.OperationCreatePool:
			params, decodeErr := decodeOperationParams[createPoolOperationParams](req.Operation.Params)
			if decodeErr != nil {
				err = decodeErr
				break
			}
			var poolTransaction chain.PreparedTransaction
			poolTransaction, err = poolTransactions.PrepareCreatePool(r.Context(), chain.CreatePoolParams{
				SettleTime: params.SettleTime, MaturityTime: params.MaturityTime,
				InterestRate: params.InterestRate, MaxLendSupply: params.MaxLendSupply,
				CollateralizationRatio: params.CollateralizationRatio,
				LendToken:              params.LendToken, CollateralToken: params.CollateralToken,
				LenderPositionToken:   params.LenderPositionToken,
				BorrowerPositionToken: params.BorrowerPositionToken,
				LiquidateRate:         params.LiquidateRate,
			})
			if err == nil {
				result, err = multisigTransactions.PrepareProposal(r.Context(), multisig.ProposalParams{
					ChainID: multisigConfig.ChainID, MultisigAddress: multisigConfig.ContractAddress,
					Operation: multisig.OperationCreatePool,
					Target:    poolTransaction.To, Value: poolTransaction.Value,
					Data: poolTransaction.Data, Nonce: req.Nonce,
				})
			}
		case multisig.OperationSettlePool:
			params, decodeErr := decodeOperationParams[settlePoolOperationParams](req.Operation.Params)
			if decodeErr != nil {
				err = decodeErr
				break
			}
			var poolTransaction chain.PreparedTransaction
			poolTransaction, err = poolTransactions.PrepareSettlePool(r.Context(), params.PoolID)
			if err == nil {
				result, err = multisigTransactions.PrepareProposal(r.Context(), multisig.ProposalParams{
					ChainID: multisigConfig.ChainID, MultisigAddress: multisigConfig.ContractAddress,
					Operation: multisig.OperationSettlePool,
					Target:    poolTransaction.To, Value: poolTransaction.Value,
					Data: poolTransaction.Data, Nonce: req.Nonce,
				})
			}
		case multisig.OperationRepayPool:
			params, decodeErr := decodeOperationParams[repayPoolOperationParams](req.Operation.Params)
			if decodeErr != nil {
				err = decodeErr
				break
			}
			var poolTransaction chain.PreparedTransaction
			poolTransaction, err = poolTransactions.PrepareRepayPool(r.Context(), chain.RepayPoolParams{
				PoolID: params.PoolID, MaxCollateralAmount: params.MaxCollateralAmount,
			})
			if err == nil {
				result, err = multisigTransactions.PrepareProposal(r.Context(), multisig.ProposalParams{
					ChainID: multisigConfig.ChainID, MultisigAddress: multisigConfig.ContractAddress,
					Operation: multisig.OperationRepayPool,
					Target:    poolTransaction.To, Value: poolTransaction.Value,
					Data: poolTransaction.Data, Nonce: req.Nonce,
				})
			}
		case multisig.OperationLiquidatePool:
			params, decodeErr := decodeOperationParams[liquidatePoolOperationParams](req.Operation.Params)
			if decodeErr != nil {
				err = decodeErr
				break
			}
			var poolTransaction chain.PreparedTransaction
			poolTransaction, err = poolTransactions.PrepareLiquidatePool(r.Context(), chain.LiquidatePoolParams{
				PoolID: params.PoolID, MaxCollateralAmount: params.MaxCollateralAmount,
			})
			if err == nil {
				result, err = multisigTransactions.PrepareProposal(r.Context(), multisig.ProposalParams{
					ChainID: multisigConfig.ChainID, MultisigAddress: multisigConfig.ContractAddress,
					Operation: multisig.OperationLiquidatePool,
					Target:    poolTransaction.To, Value: poolTransaction.Value,
					Data: poolTransaction.Data, Nonce: req.Nonce,
				})
			}
		default:
			err = fmt.Errorf("%w: unsupported operation %q", multisig.ErrInvalidProposal, req.Operation.Type)
		}
		if err != nil {
			if errors.Is(err, multisig.ErrInvalidProposal) ||
				errors.Is(err, chain.ErrInvalidCreatePool) ||
				errors.Is(err, chain.ErrInvalidSettlePool) ||
				errors.Is(err, chain.ErrInvalidRepayPool) ||
				errors.Is(err, chain.ErrInvalidLiquidatePool) {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
				return
			}
			logRequest(logger, r, slog.LevelError, "prepare multisig proposal failed", slog.Any("error", err))
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "prepare multisig proposal failed"})
			return
		}

		transactionHash, err := multisigReader.TransactionHash(r.Context(), result.Proposal)
		if err != nil {
			logRequest(logger, r, slog.LevelError, "read multisig proposal hash failed", slog.Any("error", err))
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
		Addr:              ":" + cfg.Port,
		Handler:           withCorrelationID(requestLogger(logger, withRequestDeadline(mux))),
		ReadHeaderTimeout: httpReadHeaderTimeout,
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
		MaxHeaderBytes:    httpMaxHeaderBytes,
	}
}

func withRequestDeadline(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), httpRequestDeadline)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func decodeOperationParams[T any](raw json.RawMessage) (T, error) {
	var params T
	if len(raw) == 0 {
		return params, fmt.Errorf("%w: operation params are required", multisig.ErrInvalidProposal)
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&params); err != nil {
		return params, fmt.Errorf("%w: invalid operation params: %v", multisig.ErrInvalidProposal, err)
	}
	// Try reading another JSON value. If the result is anything other than end-of-input, reject the request.
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return params, fmt.Errorf("%w: operation params must contain one JSON object", multisig.ErrInvalidProposal)
	}
	return params, nil
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, target any, limit int64) error {
	r.Body = http.MaxBytesReader(w, r.Body, limit)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("request body must contain one JSON object")
		}
		return err
	}
	return nil
}

func writeJSONBodyError(w http.ResponseWriter, err error, invalidMessage string) {
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		writeJSON(w, http.StatusRequestEntityTooLarge, errorResponse{Error: "request body too large"})
		return
	}
	writeJSON(w, http.StatusBadRequest, errorResponse{Error: invalidMessage})
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

func requireAuth(authorizer auth.Authorizer, requiredScope string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromRequest(r)
		var scopes []string
		if requiredScope != "" {
			scopes = append(scopes, requiredScope)
		}
		identity, err := authorizer.Authorize(r.Context(), token, scopes...)
		if err != nil {
			if errors.Is(err, auth.ErrInsufficientScope) {
				writeJSON(w, http.StatusForbidden, errorResponse{Error: "insufficient scope"})
				return
			}
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid token"})
			return
		}

		ctx := contextWithUsername(r, identity.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func tokenFromRequest(r *http.Request) string {
	authorization := strings.TrimSpace(r.Header.Get("Authorization"))
	return strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
}

type usernameContextKey struct{}
type requestIDContextKey struct{}

func contextWithUsername(r *http.Request, username string) context.Context {
	return context.WithValue(r.Context(), usernameContextKey{}, username)
}

func usernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameContextKey{}).(string)
	return username, ok
}

func withCorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if !validRequestID(requestID) {
			requestID = newRequestID()
		}
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestIDContextKey{}, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validRequestID(value string) bool {
	if value == "" || len(value) > 128 {
		return false
	}
	for _, character := range value {
		if (character < 'a' || character > 'z') &&
			(character < 'A' || character > 'Z') &&
			(character < '0' || character > '9') &&
			character != '-' && character != '_' && character != '.' {
			return false
		}
	}
	return true
}

func newRequestID() string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return fmt.Sprintf("%x", value)
}

func requestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDContextKey{}).(string)
	return requestID
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (w *statusRecorder) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusRecorder) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(body)
}

func requestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		recorder := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(recorder, r)
		if recorder.status == 0 {
			recorder.status = http.StatusOK
		}
		logger.Info(
			"http request completed",
			slog.String("event", "http_request"),
			slog.String("request_id", requestIDFromContext(r.Context())),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", recorder.status),
			slog.Int64("duration_ms", time.Since(started).Milliseconds()),
		)
	})
}

func logRequest(logger *slog.Logger, r *http.Request, level slog.Level, message string, attrs ...slog.Attr) {
	args := make([]any, 0, len(attrs)+1)
	args = append(args, slog.String("request_id", requestIDFromContext(r.Context())))
	for _, attr := range attrs {
		args = append(args, attr)
	}
	logger.Log(r.Context(), level, message, args...)
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	if response, ok := data.(errorResponse); ok {
		if response.Code == "" {
			response.Code = "http_" + strconv.Itoa(statusCode)
		}
		if response.RequestID == "" {
			response.RequestID = w.Header().Get("X-Request-ID")
		}
		data = response
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
