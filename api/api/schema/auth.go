package schema

const TokenType = "Bearer"

type LoginReq struct {
	LoginID  string `json:"loginId"`
	Password string `json:"password"`
}

type LoginResUser struct {
	ID  string `json:"Id"`
	LoginID string `json:"loginId"`
}

type LoginRes struct {
	AccessToken string       `json:"accessToken"`
	TokenType   string       `json:"tokenType"`
	User        LoginResUser `json:"user"`
}

type ResetPasswordReq struct {
	Email string `json:"email"`
}

type UpdatePasswordReq struct {
	Password string `json:"password"`
}

type UpdateEmailReq struct {
	Email string `json:"email"`
}
