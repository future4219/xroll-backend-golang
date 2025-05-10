package constructor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

func NewUserCreate(
	userId string,
	studentId string,
	idmUniv string,
	idmBus string,
	userType string,
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
		ID:    userId,
		Name: "",
		Age: 0,
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
		ID:    userId,
		Name: "",
		Age: 0,
	}, nil
}
