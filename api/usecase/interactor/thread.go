package interactor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type ThreadUseCase struct {
	threadRepo output_port.ThreadRepository
	ulid       output_port.ULID
	clock      output_port.Clock
}

func NewThreadUseCase(ulid output_port.ULID, threadRepo output_port.ThreadRepository, clock output_port.Clock) input_port.IThreadUseCase {
	return &ThreadUseCase{
		threadRepo: threadRepo,
		ulid:       ulid,
		clock:      clock,
	}
}

func (u *ThreadUseCase) Search(search input_port.ThreadSearch) (threads []entity.Thread, err error) {
	threads, err = u.threadRepo.Search(
		output_port.ThreadSearch{
			Limit:   search.Limit,
			Offset:  search.Offset,
			OrderBy: output_port.ThreadSearchOrderByRanking,
			Order:   output_port.ThreadSearchOrderAsc,
		},
	)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (u *ThreadUseCase) Create(thread entity.Thread) (entity.Thread, error) {
	// IDを生成する
	thread.ID = u.ulid.GenerateID()
	thread.CreatedAt = u.clock.Now()
	// Threadを保存する
	if err := u.threadRepo.Create(thread); err != nil {
		return entity.Thread{}, err
	}

	return thread, nil
}

func (u *ThreadUseCase) FindByID(id string) (entity.Thread, error) {
	thread, err := u.threadRepo.FindByID(id)
	if err != nil {
		return entity.Thread{}, err
	}
	return thread, nil
}

func (u *ThreadUseCase) FindByIDs(ids []string) ([]entity.Thread, error) {
	threads, err := u.threadRepo.FindByIDs(ids)
	if err != nil {
		return nil, err
	}
	return threads, nil
}

func (u *ThreadUseCase) Like(threadID string) error {
	thread, err := u.threadRepo.FindByID(threadID)
	if err != nil {
		return err
	}

	thread.LikeCount++
	if err := u.threadRepo.Update(thread); err != nil {
		return err
	}

	return nil
}

func (u *ThreadUseCase) Comment(threadID string, comment string) error {
	_, err := u.threadRepo.FindByID(threadID)
	if err != nil {
		return err
	}
	
	// ULIDを生成し、最初の８文字をとる
	threaderID := u.ulid.GenerateID()[:8]

	if err := u.threadRepo.CreateComment(threadID, entity.ThreadComment{
		ID:         u.ulid.GenerateID(),
		Comment:    comment,
		ThreadID:   threadID,
		ThreaderID: threaderID,
		CreatedAt:  u.clock.Now(),
	}); err != nil {
		return err
	}
	return nil
}
