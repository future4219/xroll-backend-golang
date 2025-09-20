package constructor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
)

func NewGofileUpdate(
	Id string,
	name string,
	description string,
	tagIDs []string,
	isShare bool,
) (input_port.GofileUpdate, error) {
	if Id == "" {
		return input_port.GofileUpdate{}, entconst.NewValidationError("id is empty")
	}

	if name == "" {
		return input_port.GofileUpdate{}, entconst.NewValidationError("name is empty")
	}

	return input_port.GofileUpdate{
		ID:          Id,
		Name:        name,
		Description: description,
		TagIDs:      tagIDs,
		IsShare:     isShare,
	}, nil
}
