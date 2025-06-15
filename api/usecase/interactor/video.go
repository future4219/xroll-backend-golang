package interactor

import (
	"time"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type VideoUseCase struct {
	videoRepo output_port.VideoRepository
	ulid      output_port.ULID
	clock     output_port.Clock
}

func NewVideoUseCase(ulid output_port.ULID, videoRepo output_port.VideoRepository, clock output_port.Clock) input_port.IVideoUseCase {
	return &VideoUseCase{
		videoRepo: videoRepo,
		ulid:      ulid,
		clock:     clock,
	}
}

func (u *VideoUseCase) Search(search input_port.VideoSearch) (videos []entity.Video, err error) {
	now := u.clock.Now()

	videos, err = u.videoRepo.Search(
		output_port.VideoSearch{
			Limit:  search.Limit,
			Offset: search.Offset,
			// 現在から３日前まで
			Start:   now.Add(-3 * 24 * time.Hour),
			End:     now,
			OrderBy: output_port.VideoSearchOrderByRanking,
			Order:   output_port.VideoSearchOrderAsc,
		},
	)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (u *VideoUseCase) Create(video entity.Video) (entity.Video, error) {
	// IDを生成する
	video.ID = u.ulid.GenerateID()
	video.CreatedAt = u.clock.Now()
	// 動画を保存する
	if err := u.videoRepo.Create(video); err != nil {
		return entity.Video{}, err
	}

	return video, nil
}

func (u *VideoUseCase) CreateBulk(videos []entity.Video) (err error) {
	// IDを生成する
	for i := range videos {
		videos[i].ID = u.ulid.GenerateID()
		videos[i].CreatedAt = u.clock.Now()
	}
	// 動画を保存する
	if err := u.videoRepo.CreateBulk(videos); err != nil {
		return err
	}

	return nil
}

func (u *VideoUseCase) FindByID(id string) (entity.Video, error) {
	video, err := u.videoRepo.FindByID(id)
	if err != nil {
		return entity.Video{}, err
	}
	return video, nil
}

func (u *VideoUseCase) FindByIDs(ids []string) ([]entity.Video, error) {
	videos, err := u.videoRepo.FindByIDs(ids)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (u *VideoUseCase) Like(videoID string) error {
	video, err := u.videoRepo.FindByID(videoID)
	if err != nil {
		return err
	}

	video.LikeCount++
	if err := u.videoRepo.Update(video); err != nil {
		return err
	}

	return nil
}

func (u *VideoUseCase) Comment(videoID string, comment string) error {
	_, err := u.videoRepo.FindByID(videoID)
	if err != nil {
		return err
	}

	if err := u.videoRepo.CreateComment(videoID, entity.Comment{
		ID:        u.ulid.GenerateID(),
		Comment:   comment,
		CreatedAt: u.clock.Now(),
	}); err != nil {
		return err
	}
	return nil
}
