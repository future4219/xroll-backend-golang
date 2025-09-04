package constructor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

func NewUserCreate(
	id string,
	name string,
	age int,
	userType string,
	Email *string,
	HashedPassword *string,
	GofileToken *string,
	EmailVerified bool,
	IsDeleted bool,
) (entity.User, error) {

	// usertypeはgue
	if id == "" {
		return entity.User{}, entconst.NewValidationError("id is empty")
	}
	if name == "" {
		return entity.User{}, entconst.NewValidationError("name is empty")
	}
	if age < 0 {
		return entity.User{}, entconst.NewValidationError("age is empty")
	}
	if userType != entconst.GuestUser.String() &&
		userType != entconst.MemberUser.String() &&
		userType != entconst.SystemAdmin.String() {
		return entity.User{}, entconst.NewValidationError("user type is invalid")
	}
	if Email != nil && *Email == "" {
		return entity.User{}, entconst.NewValidationError("email is empty")
	}
	if HashedPassword != nil && *HashedPassword == "" {
		return entity.User{}, entconst.NewValidationError("hashed password is empty")
	}
	if GofileToken != nil && *GofileToken == "" {
		return entity.User{}, entconst.NewValidationError("gofile token is empty")
	}

	return entity.User{
		ID:             id,
		Name:           name,
		Age:            age,
		UserType:       entconst.UserType(userType),
		Email:          Email,
		HashedPassword: HashedPassword,
		GofileToken:    GofileToken,
		EmailVerified:  EmailVerified,
		IsDeleted:      IsDeleted,
	}, nil
}

func NewUserUpdate(
	userId string,
	studentId string,
	idmUniv string,
	idmBus string,
) (entity.User, error) {

	if userId == "" {
		return entity.User{}, entconst.NewValidationError("user id is empty")
	}

	if studentId == "" {
		return entity.User{}, entconst.NewValidationError("student id is empty")
	}

	if idmUniv == "" {
		return entity.User{}, entconst.NewValidationError("idm univ is empty")
	}

	if idmBus == "" {
		return entity.User{}, entconst.NewValidationError("idm bus is empty")
	}

	return entity.User{
		ID:   userId,
		Name: "",
		Age:  0,
	}, nil
}
