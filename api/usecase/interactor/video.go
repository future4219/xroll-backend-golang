package interactor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type VideoUseCase struct {
	ulid output_port.ULID
}

func NewVideoUseCase(ulid output_port.ULID) input_port.IVideoUseCase {
	return &VideoUseCase{
		ulid:        ulid,
	}
}

func (u *VideoUseCase) Search(search input_port.VideoSearch) (videos []input_port.Video, err error) {
	videos, err = u.videoRepo.Search(search)
	if err != nil {
		return nil, err
	}
	return videos, nil
}