package model

// TwoFactorSetupRequest represents a request to set up two-factor authentication
type TwoFactorSetupRequest struct {
	Method      string  `json:"method" validate:"required,oneof=totp sms"`
	PhoneNumber *string `json:"phoneNumber,omitempty" validate:"omitempty,e164"`
}

// TwoFactorSetupResponse represents the response for 2FA setup
type TwoFactorSetupResponse struct {
	Secret      string   `json:"secret,omitempty"`
	QRCode      string   `json:"qrCode,omitempty"`
	BackupCodes []string `json:"backupCodes"`
}

// TwoFactorVerifyRequest represents a request to verify a 2FA code
type TwoFactorVerifyRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// TwoFactorDisableRequest represents a request to disable 2FA
type TwoFactorDisableRequest struct {
	Password string `json:"password" validate:"required"`
}

// RegenerateBackupCodesRequest represents a request to regenerate backup codes
type RegenerateBackupCodesRequest struct {
	Password string `json:"password" validate:"required"`
}

// TwoFactorStatusResponse represents the status of 2FA for a user
type TwoFactorStatusResponse struct {
	Enabled          bool   `json:"enabled"`
	Method           string `json:"method,omitempty"`
	BackupCodesCount int    `json:"backupCodesCount"`
}
