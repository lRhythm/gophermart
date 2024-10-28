package models

import "strings"

const (
	JWTKeyExp       = "exp"
	JWTKeyUserID    = "user_id"
	JWTHeaderKey    = "Authorization"
	JWTHeaderBearer = "Bearer"
)

// JWTHeaderValue - `JWT` to `Bearer JWT`.
func JWTHeaderValue(token string) string {
	return strings.Join([]string{JWTHeaderBearer, token}, " ")
}
