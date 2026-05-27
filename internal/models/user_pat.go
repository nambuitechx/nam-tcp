package models

type UserPATModel struct {
	ID string `json:"id"`
	UserID string `json:"user_id"`
	TargetID string `json:"target_id"`
	HashToken string `json:"hash_token"`
	CreatedAt int `json:"created_at"`
	ExpiresAt int `json:"expires_at"`
	RevokedAt int `json:"revoked_at"`
}

type CreateUserPATPayload struct {
	UserID string `json:"user_id"`
	TargetID string `json:"target_id"`
	TTLInHour int `json:"ttl_in_hour"`
}
