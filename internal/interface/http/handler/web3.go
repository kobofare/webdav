package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/yeying-community/webdav/internal/domain/user"
	"github.com/yeying-community/webdav/internal/infrastructure/auth"
	"github.com/yeying-community/webdav/internal/interface/http/dto"
	"go.uber.org/zap"
)

// Web3Handler Web3 认证处理器
type Web3Handler struct {
	web3Auth *auth.Web3Authenticator
	userRepo user.Repository
	logger   *zap.Logger
}

// NewWeb3Handler 创建 Web3 处理器
func NewWeb3Handler(
	web3Auth *auth.Web3Authenticator,
	userRepo user.Repository,
	logger *zap.Logger,
) *Web3Handler {
	return &Web3Handler{
		web3Auth: web3Auth,
		userRepo: userRepo,
		logger:   logger,
	}
}

// HandleChallenge 处理挑战请求
// GET /api/auth/challenge?address=0x123...
func (h *Web3Handler) HandleChallenge(w http.ResponseWriter, r *http.Request) {
	var address string

	// 获取地址参数
	switch r.Method {
	case http.MethodGet:
		address = r.URL.Query().Get("address")

	case http.MethodPost:
		var req dto.ChallengeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}
		address = req.Address

	default:
		h.sendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET and POST methods are allowed")
		return
	}

	if address == "" {
		h.sendError(w, http.StatusBadRequest, "MISSING_ADDRESS", "Address parameter is required")
		return
	}

	// 规范化地址
	address = strings.ToLower(strings.TrimSpace(address))

	// 检查用户是否存在
	ctx := r.Context()
	u, err := h.userRepo.FindByWalletAddress(ctx, address)
	if err != nil {
		if err == user.ErrUserNotFound {
			h.logger.Info("wallet address not registered", zap.String("address", address))
			h.sendError(w, http.StatusNotFound, "USER_NOT_FOUND", "Wallet address not registered")
			return
		}

		h.logger.Error("failed to find user", zap.String("address", address), zap.Error(err))
		h.sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process request")
		return
	}

	// 创建挑战
	challenge, err := h.web3Auth.CreateChallenge(address)
	if err != nil {
		h.logger.Error("failed to create challenge", zap.String("address", address), zap.Error(err))
		h.sendError(w, http.StatusInternalServerError, "CHALLENGE_CREATION_FAILED", "Failed to create challenge")
		return
	}

	h.logger.Info("challenge created",
		zap.String("address", address),
		zap.String("username", u.Username),
		zap.String("nonce", challenge.Nonce))

	// 返回挑战
	response := dto.ChallengeResponse{
		Nonce:     challenge.Nonce,
		Message:   challenge.Message,
		ExpiresAt: challenge.ExpiresAt,
	}

	h.sendJSON(w, http.StatusOK, response)
}

// HandleVerify 处理验证请求
// POST /api/auth/verify
func (h *Web3Handler) HandleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST method is allowed")
		return
	}

	// 解析请求
	var req dto.VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		h.sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// 验证必填字段
	if req.Address == "" {
		h.sendError(w, http.StatusBadRequest, "MISSING_ADDRESS", "Address is required")
		return
	}

	if req.Signature == "" {
		h.sendError(w, http.StatusBadRequest, "MISSING_SIGNATURE", "Signature is required")
		return
	}

	// 规范化地址
	req.Address = strings.ToLower(strings.TrimSpace(req.Address))

	// 查找用户
	ctx := r.Context()
	u, err := h.userRepo.FindByWalletAddress(ctx, req.Address)
	if err != nil {
		if err == user.ErrUserNotFound {
			h.logger.Info("wallet address not registered", zap.String("address", req.Address))
			h.sendError(w, http.StatusNotFound, "USER_NOT_FOUND", "Wallet address not registered")
			return
		}

		h.logger.Error("failed to find user", zap.String("address", req.Address), zap.Error(err))
		h.sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process request")
		return
	}

	// 验证签名并生成 token
	token, err := h.web3Auth.VerifySignature(ctx, req.Address, req.Signature)
	if err != nil {
		h.logger.Warn("signature verification failed",
			zap.String("address", req.Address),
			zap.Error(err))
		h.sendError(w, http.StatusUnauthorized, "INVALID_SIGNATURE", "Signature verification failed")
		return
	}

	h.logger.Info("user authenticated via web3",
		zap.String("address", req.Address),
		zap.String("username", u.Username))

	// 构建响应
	response := dto.VerifyResponse{
		Token:     token.Value,
		ExpiresAt: token.ExpiresAt,
		User: &dto.UserInfo{
			Username:      u.Username,
			WalletAddress: u.WalletAddress,
			Permissions:   h.getPermissionStrings(u.Permissions),
		},
	}

	h.sendJSON(w, http.StatusOK, response)
}

func (h *Web3Handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	// TODO:
}

// getPermissionStrings 获取权限字符串列表
func (h *Web3Handler) getPermissionStrings(perms *user.Permissions) []string {
	var permissions []string

	if perms.Create {
		permissions = append(permissions, "create")
	}
	if perms.Read {
		permissions = append(permissions, "read")
	}
	if perms.Update {
		permissions = append(permissions, "update")
	}
	if perms.Delete {
		permissions = append(permissions, "delete")
	}

	return permissions
}

// sendJSON 发送 JSON 响应
func (h *Web3Handler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode response", zap.Error(err))
	}
}

// sendError 发送错误响应
func (h *Web3Handler) sendError(w http.ResponseWriter, status int, code, message string) {
	response := dto.NewErrorResponse(code, message)
	h.sendJSON(w, status, response)
}
