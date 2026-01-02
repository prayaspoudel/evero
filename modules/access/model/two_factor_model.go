package model
package model

// TwoFactorSetupResponse contains QR code and backup codes for 2FA setup
type TwoFactorSetupResponse struct {









































}	Code   string `json:"code" validate:"required"`	UserID string `json:"-" validate:"required"`type UseBackupCodeRequest struct {// UseBackupCodeRequest represents a request to use a backup code}	UserID string `json:"-" validate:"required"`type RegenerateBackupCodesRequest struct {// RegenerateBackupCodesRequest represents a request to regenerate backup codes}	PhoneNumber *string `json:"phoneNumber,omitempty"`	Method      string  `json:"method" validate:"required,oneof=totp sms"`	UserID      string  `json:"-" validate:"required"`type Setup2FARequest struct {// Setup2FARequest represents a request to setup 2FA}	Password string `json:"password" validate:"required"`	UserID   string `json:"-" validate:"required"`type Disable2FARequest struct {// Disable2FARequest represents a request to disable 2FA}	Code   string `json:"code" validate:"required"`	UserID string `json:"-" validate:"required"`type Verify2FARequest struct {// Verify2FARequest represents a 2FA verification request}	Code   string `json:"code" validate:"required"`	Method string `json:"method" validate:"required,oneof=totp sms"`	UserID string `json:"-" validate:"required"`type Enable2FARequest struct {// Enable2FARequest represents a request to enable 2FA}	BackupCodes []string `json:"backupCodes"`	QRCodeURL   string   `json:"qrCodeUrl"`	Secret      string   `json:"secret"`