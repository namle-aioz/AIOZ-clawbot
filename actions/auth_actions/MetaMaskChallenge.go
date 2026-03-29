package auth_actions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	"backend/utils/response"
)

const metaMaskChallengeTTL = 5 * time.Minute

type MetaMaskChallengeInput struct {
	WalletAddress string
}

type MetaMaskChallengeOutput struct {
	Message   string
	ExpiresAt string
}

type MetaMaskChallengeAction struct{}

type metaMaskChallengeMessage struct {
	Domain        string    `json:"domain"`
	WalletAddress string    `json:"wallet_address"`
	Nonce         string    `json:"nonce"`
	IssuedAt      time.Time `json:"issued_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

func NewMetaMaskChallengeAction() MetaMaskChallengeAction {
	return MetaMaskChallengeAction{}
}

func (a MetaMaskChallengeAction) Exec(ctx context.Context, input MetaMaskChallengeInput) (*MetaMaskChallengeOutput, error) {
	_ = ctx

	walletAddress := strings.ToLower(strings.TrimSpace(input.WalletAddress))
	if walletAddress == "" || !common.IsHexAddress(walletAddress) {
		return nil, response.NewHttpError(nil, "Invalid wallet address.", http.StatusBadRequest)
	}

	now := time.Now().UTC()
	challenge := metaMaskChallengeMessage{
		Domain:        "backend",
		WalletAddress: walletAddress,
		Nonce:         uuid.NewString(),
		IssuedAt:      now,
		ExpiresAt:     now.Add(metaMaskChallengeTTL),
	}

	payload, err := json.Marshal(challenge)
	if err != nil {
		return nil, response.NewInternalError(fmt.Errorf("marshal metamask challenge: %w", err))
	}

	return &MetaMaskChallengeOutput{
		Message:   string(payload),
		ExpiresAt: challenge.ExpiresAt.Format(time.RFC3339),
	}, nil
}
