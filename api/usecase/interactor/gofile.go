package interactor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type GofileUseCase struct {
	Gofile output_port.Gofile
	Video  output_port.VideoRepository
	ulid   output_port.ULID
}

func NewGofileUseCase(gofile output_port.Gofile, videoRepo output_port.VideoRepository, ulid output_port.ULID) input_port.IGofileUseCase {
	return &GofileUseCase{
		Gofile: gofile,
		Video:  videoRepo,
		ulid:   ulid,
	}
}

func (u *GofileUseCase) ProxyVideo(gofileURL string) () {

}
