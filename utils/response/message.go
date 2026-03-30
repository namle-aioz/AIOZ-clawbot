package response

const (
	ErrCodeHeaderNotExit      = "Header Authorization not exist"
	ErrCodeAuthTokenInvalid   = "Auth token invalid"
	ErrUserNotExit            = "User not exist"
	ErrInvalidUserAndPassword = "Invalid user and password"

	ErrInvalidRequestBody                 = "Invalid request body."
	ErrRecordNotFound                     = "Record not found."
	ErrAccountNotVerified                 = "Your account is not verified."
	ErrInvalidOrExpiredOTPCode            = "Invalid or expired OTP code."
	ErrEmailAlreadyInUse                  = "Email already in use."
	ErrPreviousOTPCodeStillValid          = "Previous OTP code is still valid. Please wait before requesting a new one."
	ErrTooManyOTPRequests                 = "Too many OTP requests. Please try again later."
	ErrInvalidWalletAddress               = "Invalid wallet address."
	ErrInvalidMetaMaskChallenge           = "Invalid MetaMask challenge."
	ErrWalletAddressDoesNotMatchChallenge = "Wallet address does not match challenge."
	ErrMetaMaskChallengeExpired           = "MetaMask challenge expired."
	ErrInvalidMetaMaskSignature           = "Invalid MetaMask signature."
	ErrOTPSessionNotFound                 = "OTP session not found. Please request a new code."
)
