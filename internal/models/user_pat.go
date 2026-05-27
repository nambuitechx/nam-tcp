package models

type UserPATModel struct {
	ID string `json:"id"`
	UserID string `json:"user_id"`
	HashToken string `json:"hash_token"`
	CreatedAt int `json:"created_at"`
	ExpiresAt int `json:"expires_at"`
	RevokedAt int `json:"revoked_at"`
	TargetHost string `json:"target_host"`
	TargetPort string `json:"target_port"`
}

type CreateUserPATPayload struct {
	UserID string `json:"user_id"`
	TargetHost string `json:"target_host"`
	TargetPort string `json:"target_port"`
	TTLInHour int `json:"ttl_in_hour"`
}

