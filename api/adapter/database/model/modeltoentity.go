package model

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type model[T any] interface {
	Entity() T
}

// ToEntities converts slice of models to slice of entities.
func ToEntities[Entity any, Model model[Entity]](models []Model) []Entity {
	ret := make([]Entity, len(models))
	for i, model := range models {
		ret[i] = model.Entity()
	}
	return ret
}

func (u User) Entity() entity.User {
	return entity.User{
		ID:   u.ID,
		Name: u.Name,
		Age:  u.Age,
	}
}
