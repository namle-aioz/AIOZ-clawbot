package auth_actions

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"backend/actions/mail_actions"
	"backend/models"
	"backend/templates"
	"backend/utils/otp"
	"backend/utils/response"
	utils "backend/utils/token"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type VerifyOTPInput struct {
	Email    string
	Password string
	OTPCode  string
}

type SignUpInput struct {
	Email string
}

type SignUpOutput struct {
	AccessToken  string
	RefreshToken string
	User         *models.User
}

type SignUpAction struct {
	userRepo    models.UserRepository
	otpRepo     models.OTPCodeRepository
	mailSender  mail_actions.MailSender
	tokenIssuer utils.TokenIssuer
	bcrypt_cost int
}

var OTPExpireTTL = 5 * time.Minute

func NewSignUpAction(userRepo models.UserRepository, otpRepo models.OTPCodeRepository, mailSender mail_actions.MailSender, tokenIssuer utils.TokenIssuer, bcrypt_cost int) SignUpAction {
	return SignUpAction{userRepo: userRepo, otpRepo: otpRepo, mailSender: mailSender, tokenIssuer: tokenIssuer, bcrypt_cost: bcrypt_cost}
}

func (sa SignUpAction) SignUp(ctx context.Context, input SignUpInput) error {
	user, err := sa.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return response.NewInternalError(err)
	} else if user != nil && user.IsVerified && user.Status {
		return response.NewHttpError(nil, response.ErrEmailAlreadyInUse, http.StatusBadRequest)
	}

	otpCode, err := otp.GenerateOTP(6)
	if err != nil {
		return response.NewInternalError(err)
	}

	sa.sendOTPEmail(context.Background(), input.Email, otpCode)

	now := time.Now().UTC()
	hashedOTP, err := bcrypt.GenerateFromPassword([]byte(otpCode), sa.bcrypt_cost)
	if err != nil {
		return response.NewInternalError(err)
	}

	err = sa.otpRepo.CreateOrUpdateOTPCode(ctx, &models.OTPCode{
		Email:       input.Email,
		Code:        hashedOTP,
		IssuedAt:    now,
		ExpiredAt:   now.Add(OTPExpireTTL),
		Attempts:    1,
		MaxAttempts: 3,
	})
	if err != nil {
		return response.NewInternalError(err)
	}

	return nil
}

func (sa SignUpAction) VerifyOTP(ctx context.Context, input VerifyOTPInput) (*SignUpOutput, error) {
	otpCode, err := sa.otpRepo.GetOTPCodeByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, response.NewInternalError(err)
	} else if otpCode == nil {
		return nil, response.NewHttpError(nil, response.ErrInvalidOrExpiredOTPCode, http.StatusBadRequest)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(otpCode.Code), []byte(input.OTPCode)); err != nil {
		return nil, response.NewHttpError(nil, response.ErrInvalidOrExpiredOTPCode, http.StatusBadRequest)
	}

	if err := bcrypt.CompareHashAndPassword(otpCode.Code, []byte(input.OTPCode)); err != nil || time.Now().UTC().After(otpCode.ExpiredAt) {
		return nil, response.NewHttpError(nil, response.ErrInvalidOrExpiredOTPCode, http.StatusBadRequest)
	}

	user, err := sa.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, response.NewInternalError(err)
	} else if user != nil && user.IsVerified && user.Status {
		return nil, response.NewHttpError(nil, response.ErrEmailAlreadyInUse, http.StatusBadRequest)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), sa.bcrypt_cost)
	if err != nil {
		return nil, response.NewInternalError(err)
	}

	if user == nil {
		user = &models.User{
			Email:      input.Email,
			IsVerified: true,
			Status:     true,
			AuthMethod: "email",
			Password:   hashedPassword,
		}
		if err := sa.userRepo.CreateUser(ctx, user); err != nil {
			return nil, response.NewInternalError(err)
		}
	} else {
		user.IsVerified = true
		user.Status = true
		user.AuthMethod = "email"
		user.Password = hashedPassword
		if err := sa.userRepo.UpdateUser(ctx, user); err != nil {
			return nil, response.NewInternalError(err)
		}
	}

	accessToken, refreshToken, _, err := sa.tokenIssuer.CreateCredential(user.Id)
	if err != nil {
		return nil, response.NewInternalError(err)
	}

	return &SignUpOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (sa SignUpAction) ResendOTPcode(ctx context.Context, input SignUpInput) error {
	user, err := sa.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return response.NewInternalError(err)
	} else if user != nil && user.IsVerified && user.Status {
		return response.NewHttpError(nil, response.ErrEmailAlreadyInUse, http.StatusBadRequest)
	}

	prevOTPCode, err := sa.otpRepo.GetOTPCodeByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return response.NewInternalError(err)
	} else if prevOTPCode != nil && time.Now().UTC().Before(prevOTPCode.ExpiredAt) {
		return response.NewHttpError(nil, response.ErrPreviousOTPCodeStillValid, http.StatusBadRequest)
	}

	prevAttempts := 0
	if prevOTPCode != nil {
		prevAttempts = prevOTPCode.Attempts
		if prevOTPCode.MaxAttempts > 0 && prevAttempts >= prevOTPCode.MaxAttempts {
			return response.NewHttpError(nil, response.ErrTooManyOTPRequests, http.StatusTooManyRequests)
		}
	}

	otpCode, err := otp.GenerateOTP(6)
	if err != nil {
		return response.NewInternalError(err)
	}

	sa.sendOTPEmail(context.Background(), input.Email, otpCode)

	hashedOTP, err := bcrypt.GenerateFromPassword([]byte(otpCode), sa.bcrypt_cost)
	if err != nil {
		return response.NewInternalError(err)
	}

	now := time.Now().UTC()
	err = sa.otpRepo.CreateOrUpdateOTPCode(ctx, &models.OTPCode{
		Email:       input.Email,
		Code:        hashedOTP,
		IssuedAt:    now,
		ExpiredAt:   now.Add(OTPExpireTTL),
		Attempts:    prevAttempts + 1,
		MaxAttempts: 3,
	})
	if err != nil {
		return response.NewInternalError(err)
	}

	return nil
}

func (sa SignUpAction) sendOTPEmail(ctx context.Context, email string, otp string) {
	go func() {
		template := templates.BuildOTPEmail(email, otp, OTPExpireTTL)

		mailInput := mail_actions.SendMailInput{
			To:      []string{email},
			Subject: "Sign up",
			Html:    template,
		}

		if err := sa.mailSender.SendOTPMail(ctx, mailInput); err != nil {
			slog.Error("send OTP mail failed",
				slog.String("email", email),
				slog.Any("err", err),
			)
		}
	}()
}
