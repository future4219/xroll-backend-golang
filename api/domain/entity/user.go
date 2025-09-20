package entity

import (
	"time"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
)

type User struct {
	ID             string
	Name           string
	Age            int
	UserType       entconst.UserType
	Bio            string
	Email          *string
	HashedPassword *string
	GofileToken    *string
	EmailVerified  bool
	IsDeleted      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
