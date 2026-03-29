package auth_actions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"gorm.io/gorm"

	"backend/models"
	"backend/utils/response"
	utils "backend/utils/token"
)

type MetaMaskSignInInput struct {
	WalletAddress string
	Message       string
	Signature     string
}

type MetaMaskSignInAction struct {
	userRepo    models.UserRepository
	tokenIssuer utils.TokenIssuer
}

func NewMetaMaskSignInAction(userRepo models.UserRepository, tokenIssuer utils.TokenIssuer) MetaMaskSignInAction {
	return MetaMaskSignInAction{userRepo: userRepo, tokenIssuer: tokenIssuer}
}

func (a MetaMaskSignInAction) Exec(ctx context.Context, input MetaMaskSignInInput) (*SignInOutput, error) {
	challenge, err := parseMetaMaskChallengeMessage(input.Message)
	if err != nil {
		return nil, response.NewHttpError(nil, "Invalid MetaMask challenge.", http.StatusBadRequest)
	}

	walletAddress := normalizeWalletAddress(input.WalletAddress)
	if walletAddress == "" || challenge.WalletAddress != walletAddress {
		return nil, response.NewHttpError(nil, "Wallet address does not match challenge.", http.StatusUnauthorized)
	}

	if time.Now().UTC().After(challenge.ExpiresAt) {
		return nil, response.NewHttpError(nil, "MetaMask challenge expired.", http.StatusUnauthorized)
	}

	recoveredAddress, err := recoverWalletAddress(input.Message, input.Signature)
	if err != nil {
		return nil, response.NewHttpError(nil, "Invalid MetaMask signature.", http.StatusUnauthorized)
	}

	if recoveredAddress != walletAddress {
		return nil, response.NewHttpError(nil, "Invalid MetaMask signature.", http.StatusUnauthorized)
	}

	user, err := a.findOrCreateWalletUser(ctx, walletAddress)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, _, err := a.tokenIssuer.CreateCredential(user.Id)
	if err != nil {
		return nil, response.NewInternalError(err)
	}

	return &SignInOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (a MetaMaskSignInAction) findOrCreateWalletUser(ctx context.Context, walletAddress string) (*models.User, error) {
	user, err := a.userRepo.GetUserByWalletAddress(ctx, walletAddress)
	if err == nil {
		return a.syncWalletUser(ctx, user, walletAddress)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, response.NewInternalError(err)
	}

	walletAddressCopy := walletAddress
	user = &models.User{
		DisplayEmail:  walletAddress,
		Password:      "",
		IsVerified:    true,
		AuthProvider:  "metamask",
		WalletAddress: &walletAddressCopy,
	}

	if err := a.userRepo.CreateUser(ctx, user); err != nil {
		return nil, response.NewInternalError(err)
	}

	return user, nil
}

func (a MetaMaskSignInAction) syncWalletUser(ctx context.Context, user *models.User, walletAddress string) (*models.User, error) {
	walletAddressCopy := walletAddress
	user.DisplayEmail = walletAddress
	user.IsVerified = true
	user.AuthProvider = "metamask"
	user.WalletAddress = &walletAddressCopy

	if err := a.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, response.NewInternalError(err)
	}

	return user, nil
}

func parseMetaMaskChallengeMessage(message string) (*metaMaskChallengeMessage, error) {
	var challenge metaMaskChallengeMessage
	if err := json.Unmarshal([]byte(message), &challenge); err != nil {
		return nil, err
	}

	challenge.WalletAddress = normalizeWalletAddress(challenge.WalletAddress)
	if challenge.WalletAddress == "" || challenge.Nonce == "" || challenge.ExpiresAt.IsZero() {
		return nil, fmt.Errorf("invalid challenge payload")
	}

	return &challenge, nil
}

func recoverWalletAddress(message string, signature string) (string, error) {
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return "", err
	}

	if len(sig) != 65 {
		return "", fmt.Errorf("invalid signature length")
	}

	if sig[64] >= 27 {
		sig[64] -= 27
	}

	if sig[64] != 0 && sig[64] != 1 {
		return "", fmt.Errorf("invalid signature recovery id")
	}

	pubKey, err := crypto.SigToPub(accounts.TextHash([]byte(message)), sig)
	if err != nil {
		return "", err
	}

	return normalizeWalletAddress(crypto.PubkeyToAddress(*pubKey).Hex()), nil
}

func normalizeWalletAddress(walletAddress string) string {
	return strings.ToLower(strings.TrimSpace(walletAddress))
}
