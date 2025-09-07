package schema

import "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"

const TokenType = "Bearer"

type BootRes struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	UserID      string `json:"userId"`
}
type LoginReq struct {
	LoginID  string `json:"loginId"`
	Password string `json:"password"`
}

type LoginResUser struct {
	ID      string `json:"Id"`
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

type CreateByMeReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyEmailReq struct {
	Email              string `json:"email"`
	AuthenticationCode string `json:"authenticationCode"`
}

type VerifyEmailRes struct {
	AccessToken string             `json:"accessToken"`
	TokenType   string             `json:"tokenType"`
	User        VerifyEmailUserRes `json:"user"`
}

type VerifyEmailUserRes struct {
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	UserType string `json:"userType"`
}

func VerifyEmailResFromEntity(user entity.User, token string) VerifyEmailRes {
	return VerifyEmailRes{
		AccessToken: token,
		TokenType:   TokenType,
		User:        VerifyEmailResUserFromEntity(user),
	}
}

func VerifyEmailResUserFromEntity(user entity.User) VerifyEmailUserRes {
	return VerifyEmailUserRes{
		UserID:   user.ID,
		Email:    *user.Email,
		UserType: string(user.UserType),
	}
}
