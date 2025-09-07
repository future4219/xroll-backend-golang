package interactor

import (
	"fmt"
	"net/url"
	"os"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type GofileUseCase struct {
	gofileRepo      output_port.GofileRepository
	gofileAPIDriver output_port.GofileAPIDriver
	userRepo        output_port.UserRepository
	ulid            output_port.ULID
	clock           output_port.Clock
}

func NewGofileUseCase(ulid output_port.ULID, gofileRepo output_port.GofileRepository, gofileAPIDriver output_port.GofileAPIDriver, userRepo output_port.UserRepository, clock output_port.Clock) input_port.IGofileUseCase {
	return &GofileUseCase{
		gofileRepo:      gofileRepo,
		gofileAPIDriver: gofileAPIDriver,
		userRepo:        userRepo,
		ulid:            ulid,
		clock:           clock,
	}
}

func (u *GofileUseCase) Create(user entity.User, gofileCreate input_port.GofileCreate) (entity.GofileVideo, error) {
	// 先にDireftlinkを発行させないといけない
	fmt.Println("-------------------------------start issuing direct link-------------------------------")
	IssueDirectLinkRes, err := u.gofileAPIDriver.IssueDirectLink(
		gofileCreate.GofileID,
		os.Getenv("GOFILE_API_KEY"),
	)
	if err != nil {
		return entity.GofileVideo{}, fmt.Errorf("issue direct link: %w", err)
	}
	fmt.Println("-------------------------------\nend issuing direct link\n-------------------------------")

	fmt.Printf("-------------------------------")
	fmt.Printf("Issued Direct Link: %+v\n", IssueDirectLinkRes)
	fmt.Printf("-------------------------------")

	fmt.Println("-------------------------------\nstart getting content\n-------------------------------")
	gofileGetContentRes, err := u.gofileAPIDriver.GetContent(
		gofileCreate.GofileID,
		os.Getenv("GOFILE_API_KEY"),
	)
	if err != nil {
		return entity.GofileVideo{}, err
	}
	fmt.Println("-------------------------------\nend getting content\n-------------------------------")

	var gofileDirectLink string
	gofileDirectLink = IssueDirectLinkRes.DirectLink
	if gofileDirectLink == "" {
		return entity.GofileVideo{}, fmt.Errorf("no directLink found in Gofile response")
	}

	var videoURL string
	videoURL = os.Getenv("XROLL_API_ENDPOINT") + "/gofile/proxy?url=" + url.QueryEscape(gofileDirectLink)

	var gofileTags []entity.GofileTag
	for _, tagID := range gofileCreate.TagIDs {
		gofileTags = append(gofileTags, entity.GofileTag{
			ID: tagID,
		})
	}

	entityGofile := entity.GofileVideo{
		ID:              u.ulid.GenerateID(),
		Name:            gofileCreate.Name,
		GofileID:        gofileCreate.GofileID,
		GofileDirectURL: gofileDirectLink,
		VideoURL:        videoURL,
		ThumbnailURL:    gofileGetContentRes.Data.Thumbnail,
		IsShared:        false, // 初期状態では公開されていない
		UserID:          user.ID,
		GofileTags:      gofileTags,
	}

	if err := u.gofileRepo.Create(entityGofile); err != nil {
		return entity.GofileVideo{}, err
	}

	res, err := u.gofileRepo.FindByID(entityGofile.ID)
	if err != nil {
		return entity.GofileVideo{}, err
	}

	// GofileTokenがUserに登録されていなければ、登録する
	if user.GofileToken != nil && *user.GofileToken != "" {
		user.GofileToken = gofileCreate.GofileToken
		if err := u.userRepo.Update(user); err != nil {
			return entity.GofileVideo{}, err
		}
	}

	return res, nil
}

func (u *GofileUseCase) FindByID(user entity.User, id string) (entity.GofileVideo, error) {
	if id == "" {
		return entity.GofileVideo{}, fmt.Errorf("id is required")
	}

	video, err := u.gofileRepo.FindByID(id)
	if err != nil {
		return entity.GofileVideo{}, err
	}

	// 自分の動画か、共有されている動画のみ閲覧可能
	if video.UserID != user.ID && !video.IsShared {
		// セキュリティを考慮して、notfoundで返す
		return entity.GofileVideo{}, fmt.Errorf("%w: video not found", ErrKind.NotFound)
	}

	return video, nil
}

// userIDが持っているVideoを返す
func (u *GofileUseCase) FindByUserID(user entity.User) ([]entity.GofileVideo, error) {
	if user.ID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	videos, err := u.gofileRepo.FindByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (u *GofileUseCase) FindByUserIDShared(user entity.User, targetUserID string) ([]entity.GofileVideo, error) {
	if targetUserID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	videos, err := u.gofileRepo.FindByUserIDShared(targetUserID)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (u *GofileUseCase) UpdateIsShareVideo(user entity.User, videoID string, isShare bool) error {
	if videoID == "" {
		return fmt.Errorf("videoID is required")
	}

	video, err := u.gofileRepo.FindByID(videoID)
	if err != nil {
		return err
	}

	// 自分の動画のみ更新可能
	if video.UserID != user.ID {
		return fmt.Errorf("you do not have permission to update this video")
	}

	video.IsShared = isShare
	fmt.Printf("-------------------------------\nupdate video isShare: %v\n-------------------------------\n", isShare)
	fmt.Printf("video: %+v\n", video)
	if err := u.gofileRepo.Update(video); err != nil {
		return err
	}

	return nil
}
