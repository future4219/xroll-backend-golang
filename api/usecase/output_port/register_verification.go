package output_port

import "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"

type RegisterVerificationRepository interface {
	UpsertInTx(tx interface{}, entity entity.RegisterVerification) error
	FindByEmail(email string) (entity.RegisterVerification, error)
	DeleteByEmailInTx(tx interface{}, email string) error
}
