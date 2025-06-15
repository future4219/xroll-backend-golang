package interactor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type TwitterUseCase struct {
	Twitter output_port.Twitter
	Video   output_port.VideoRepository
	ulid    output_port.ULID
}

func NewTwitterUseCase(twitter output_port.Twitter, videoRepo output_port.VideoRepository, ulid output_port.ULID) input_port.ITwitterUseCase {
	return &TwitterUseCase{
		Twitter: twitter,
		Video:   videoRepo,
		ulid:    ulid,
	}
}

func (u *TwitterUseCase) GetVideoByURL(url string) (string, error) {
	videoUrl, err := u.Twitter.GetVideoByURL(url)
	if err != nil {
		return "", err
	}

	// 投稿された動画を保存しておく
	thumbnailUrl, err := u.Twitter.GetThumbnailByURL(url)
	if err != nil {
		return "", err
	}
	create := entity.Video{
		ID:            u.ulid.GenerateID(),
		Ranking:       ,
		VideoURL:      videoUrl,
		ThumbnailURL:  thumbnailUrl,
		DownloadCount: 0,
		LikeCount:     0,
		Comments:      []entity.Comment{},
	}
	if err := u.Video.Create(create); err != nil {
		return "", err
	}

	return videoUrl, nil
}
