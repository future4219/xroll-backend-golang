package schema

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type UserRes struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Bio       string `json:"bio"`
	UserType  string `json:"user_type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UsersRes struct {
	List  []UserRes `json:"list"`
	Total int       `json:"total"`
}

type CreateUserReq struct {
	Name     string  `json:"name" validate:"required"`
	Age      int     `json:"age" validate:"required,gte=0"`
	UserType string  `json:"userType" validate:"required,oneof=guest user admin"`
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required,min=8"`
}

type CreateUserReqByAdmin struct {
	StudentID string `json:"studentId"`
	IdmUniv   string `json:"idmUniv"`
	IdmBus    string `json:"idmBus"`
	UserType  string `json:"userType"`
}

type UpdateUserReq struct {
	Name string `json:"name" validate:"required"`
	Bio  string `json:"bio"`
}

type UserSearchQueryReq struct {
	Query    string `query:"q"`
	UserType string `query:"user-type"`
	Skip     int    `query:"skip"`
	Limit    int    `query:"limit"`
}

func UserResFromEntity(user entity.User) UserRes {
	return UserRes{
		ID:        user.ID,
		Name:      user.Name,
		Age:       user.Age,
		Bio:       user.Bio,
		UserType:  user.UserType.String(),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}

func UserWithDetailResFromEntity(user entity.User) UserRes {
	return UserRes{
		ID:   user.ID,
		Name: user.Name,
		Age:  user.Age,
	}
}

func UsersResFromSearchResult(list []entity.User, total int) UsersRes {
	users := make([]UserRes, len(list))
	for i, user := range list {
		users[i] = UserWithDetailResFromEntity(user)
	}
	return UsersRes{
		List:  users,
		Total: total,
	}
}
