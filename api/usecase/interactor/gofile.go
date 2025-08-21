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

func (u *GofileUseCase) Create(gofileCreate input_port.GofileCreate) (entity.GofileVideo, error) {
	var User entity.User
	// TODO トランザクションにする
	if gofileCreate.UserID == nil {
		User = entity.User{
			ID:          u.ulid.GenerateID(),
			Name:        "guest100000",
			UserType:    "guest",
			GofileToken: gofileCreate.GofileToken,
		}
		if err := u.userRepo.Create(User); err != nil {
			return entity.GofileVideo{}, err
		}
	} else {
		var err error
		User, err = u.userRepo.FindByID(*gofileCreate.UserID)
		if err != nil {
			return entity.GofileVideo{}, err
		}
	}

	// 先にDireftlinkを発行させないといけない
	fmt.Println("-------------------------------start issuing direct link-------------------------------")
	IssueDirectLinkRes, err := u.gofileAPIDriver.IssueDirectLink(
		gofileCreate.GofileID,
		os.Getenv("GOFILE_API_KEY"),
	)
	if err != nil {
		return entity.GofileVideo{}, fmt.Errorf("issue direct link: %w", err)
	}
	fmt.Println("-------------------------------end issuing direct link-------------------------------")

	fmt.Printf("-------------------------------")
	fmt.Printf("Issued Direct Link: %+v\n", IssueDirectLinkRes)
	fmt.Printf("-------------------------------")

	fmt.Println("-------------------------------start getting content-------------------------------")
	gofileGetContentRes, err := u.gofileAPIDriver.GetContent(
		gofileCreate.GofileID,
		os.Getenv("GOFILE_API_KEY"),
	)
	if err != nil {
		return entity.GofileVideo{}, err
	}
	fmt.Println("-------------------------------end getting content-------------------------------")

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
		UserID:          User.ID,
		User:            User,
		GofileTags:      gofileTags,
	}

	if err := u.gofileRepo.Create(entityGofile); err != nil {
		return entity.GofileVideo{}, err
	}

	res, err := u.gofileRepo.FindByID(entityGofile.ID)
	if err != nil {
		return entity.GofileVideo{}, err
	}

	return res, nil
}

func (u *GofileUseCase) FindByID(id string) (entity.GofileVideo, error) {
	if id == "" {
		return entity.GofileVideo{}, fmt.Errorf("id is required")
	}

	video, err := u.gofileRepo.FindByID(id)
	if err != nil {
		return entity.GofileVideo{}, err
	}

	return video, nil
}


// userIDが持っているVideoを返す
func (u *GofileUseCase) FindByUserID(userID string) ([]entity.GofileVideo, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	videos, err := u.gofileRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return videos, nil
}
