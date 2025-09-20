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
	Id string,
	name string,
	bio string,
) (entity.User, error) {

	if Id == "" {
		return entity.User{}, entconst.NewValidationError("id is empty")
	}
	if name == "" {
		return entity.User{}, entconst.NewValidationError("name is empty")
	}
	if len(name) > 100 {
		return entity.User{}, entconst.NewValidationError("name is too long")
	}
	if len(bio) > 500 {
		return entity.User{}, entconst.NewValidationError("bio is too long")
	}

	return entity.User{
		ID:   Id,
		Name: name,
		Bio:  bio,
	}, nil
}
