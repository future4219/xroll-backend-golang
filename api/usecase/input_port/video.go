package input_port

import "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"

type VideoSearch struct {
	Limit  int
	Offset int
}

type IVideoUseCase interface {
	Search(search VideoSearch) (entity.Video, error)
}
