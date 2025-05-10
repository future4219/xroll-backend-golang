package schema

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type UserRes struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	
}

type UsersRes struct {
	List  []UserRes `json:"list"`
	Total int       `json:"total"`
}

type CreateUserReq struct {
	StudentID string `json:"student_id"`
	IdmUniv   string `json:"idm_univ"`
	IdmBus    string `json:"idm_bus"`
}

type CreateUserReqByAdmin struct {
	StudentID string `json:"studentId"`
	IdmUniv   string `json:"idmUniv"`
	IdmBus    string `json:"idmBus"`
	UserType  string `json:"userType"`
}

type UpdateUserReq struct {
	StudentID string `json:"studentId"`
	IdmUniv   string `json:"idmUniv"`
	IdmBus    string `json:"idmBus"`
	UserType  string `json:"userType"`
}

type UserSearchQueryReq struct {
	Query    string `query:"q"`
	UserType string `query:"user-type"`
	Skip     int    `query:"skip"`
	Limit    int    `query:"limit"`
}

func UserResFromEntity(user entity.User) UserRes {
	return UserRes{
		ID:    user.ID,
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
