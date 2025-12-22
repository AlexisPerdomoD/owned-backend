package auth

type Session struct {
	UserID  string `json:"user_id"`
	Role    int    `json:"role"`
	DateISO int    `json:"date_iso"`
}
