package interactor

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/constructor"
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
	twitter         output_port.Twitter
}

func NewGofileUseCase(ulid output_port.ULID, gofileRepo output_port.GofileRepository, gofileAPIDriver output_port.GofileAPIDriver, userRepo output_port.UserRepository, clock output_port.Clock, twitter output_port.Twitter) input_port.IGofileUseCase {
	return &GofileUseCase{
		gofileRepo:      gofileRepo,
		gofileAPIDriver: gofileAPIDriver,
		userRepo:        userRepo,
		ulid:            ulid,
		clock:           clock,
		twitter:         twitter,
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
		Description:     "",    // 最初は空文字
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

func (u *GofileUseCase) Update(user entity.User, update input_port.GofileUpdate) (entity.GofileVideo, error) {
	updatingGofile, err := constructor.NewGofileUpdate(
		update.ID,
		update.Name,
		update.Description,
		update.TagIDs,
		update.IsShare,
	)
	if err != nil {
		return entity.GofileVideo{}, fmt.Errorf("failed to construct gofile update: %w", err)
	}

	video, err := u.gofileRepo.FindByID(updatingGofile.ID)
	if err != nil {
		return entity.GofileVideo{}, err
	}

	// 自分の動画のみ更新可能
	if video.UserID != user.ID {
		//　セキュリティ面を考慮して、非公開動画であればnotfoundで返す
		if video.IsShared {
			return entity.GofileVideo{}, fmt.Errorf("%w: you do not have permission to update this video", ErrKind.Unauthorized)
		} else {
			return entity.GofileVideo{}, fmt.Errorf("%w: video not found", ErrKind.NotFound)
		}
	}

	video.Name = updatingGofile.Name
	video.Description = updatingGofile.Description
	video.GofileTags = func() []entity.GofileTag {
		tags := make([]entity.GofileTag, 0, len(updatingGofile.TagIDs))
		for _, tagID := range updatingGofile.TagIDs {
			tags = append(tags, entity.GofileTag{ID: tagID})
		}
		return tags
	}()
	video.IsShared = updatingGofile.IsShare

	if err := u.gofileRepo.Update(video); err != nil {
		return entity.GofileVideo{}, err
	}

	res, err := u.gofileRepo.FindByID(video.ID)
	if err != nil {
		return entity.GofileVideo{}, err
	}

	return res, nil
}

func (u *GofileUseCase) FindByID(user entity.User, id string) (entity.GofileVideo, bool, error) {
	if id == "" {
		return entity.GofileVideo{}, false, fmt.Errorf("id is required")
	}

	video, err := u.gofileRepo.FindByID(id)
	if err != nil {
		return entity.GofileVideo{}, false, err
	}

	// 自分の動画か、共有されている動画のみ閲覧可能
	if video.UserID != user.ID && !video.IsShared {
		// セキュリティを考慮して、notfoundで返す
		return entity.GofileVideo{}, false, fmt.Errorf("%w: video not found", ErrKind.NotFound)
	}

	// videoに対してuserがLikeしているかどうか
	hasLike, err := u.gofileRepo.HasLike(user.ID, video.ID)
	if err != nil {
		return entity.GofileVideo{}, false, err
	}

	return video, hasLike, nil
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

	if err := u.gofileRepo.Update(video); err != nil {
		return err
	}

	return nil
}

func (u *GofileUseCase) Delete(user entity.User, videoID string) error {
	if videoID == "" {
		return fmt.Errorf("videoID is required")
	}

	video, err := u.gofileRepo.FindByID(videoID)
	if err != nil {
		return err
	}

	// 自分の動画のみ削除可能
	if video.UserID != user.ID {
		return fmt.Errorf("you do not have permission to delete this video")
	}

	if err := u.gofileRepo.Delete(videoID); err != nil {
		return err
	}

	return nil
}

// userがいいねした動画を返す
func (u *GofileUseCase) FindLikedVideos(user entity.User) ([]entity.GofileVideo, error) {
	if user.ID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	videos, err := u.gofileRepo.FindLikedVideos(user.ID)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (u *GofileUseCase) LikeVideo(user entity.User, videoID string) error {
	if user.ID == "" || videoID == "" {
		return fmt.Errorf("userID and videoID are required")
	}

	// 動画取得 & 閲覧可否（非公開は所有者のみ）
	video, err := u.gofileRepo.FindByID(videoID)
	if err != nil {
		return err
	}
	if video.UserID != user.ID && !video.IsShared {
		return fmt.Errorf("%w: video not found", ErrKind.NotFound)
	}

	// 既にLike済みなら冪等に成功扱い
	has, err := u.gofileRepo.HasLike(user.ID, videoID)
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	// いいね作成（ユニーク制約で二重押しを遮断）
	like := entity.GofileVideoLike{
		ID:            u.ulid.GenerateID(),
		GofileVideoID: videoID,
		UserID:        user.ID,
		CreatedAt:     u.clock.Now(),
	}
	if err := u.gofileRepo.CreateLike(like); err != nil {
		// すでに存在（ユニーク違反）＝並走二重押し → 冪等成功扱い
		if output_port.IsUniqueViolation(err) {
			return nil
		}
		return err
	}

	// LikeCountを+1（レース時も最終的に正しく近づく）
	video.LikeCount += 1
	return u.gofileRepo.Update(video)
}

func (u *GofileUseCase) UnlikeVideo(user entity.User, videoID string) error {
	if user.ID == "" || videoID == "" {
		return fmt.Errorf("userID and videoID are required")
	}

	// 動画取得 & 閲覧可否（非公開は所有者のみ）
	video, err := u.gofileRepo.FindByID(videoID)
	if err != nil {
		return err
	}
	if video.UserID != user.ID && !video.IsShared {
		return fmt.Errorf("%w: video not found", ErrKind.NotFound)
	}

	// 未Likeなら冪等に成功扱い
	has, err := u.gofileRepo.HasLike(user.ID, videoID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	// いいね削除（削除できたらだけカウントを-1）
	affected, err := u.gofileRepo.DeleteLike(user.ID, videoID)
	if err != nil {
		return err
	}
	if affected == 0 {
		// 既に消えていた等 → 冪等成功扱い
		return nil
	}

	if video.LikeCount > 0 {
		video.LikeCount -= 1
	}
	return u.gofileRepo.Update(video)
}

func (u *GofileUseCase) Search(user entity.User, query input_port.GofileSearchQuery) ([]entity.GofileVideo, error) {
	videos, err := u.gofileRepo.Search(output_port.GofileSearchQuery{
		Q:       query.Q,
		Skip:    query.Skip,
		Limit:   query.Limit,
		OrderBy: query.OrderBy,
		Order:   query.Order,
	})
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (u *GofileUseCase) CreateComment(user entity.User, input input_port.GofileVideoCommentCreate) (entity.GofileVideoComment, error) {
	if user.ID == "" {
		return entity.GofileVideoComment{}, fmt.Errorf("userID is required")
	}
	if input.VideoID == "" {
		return entity.GofileVideoComment{}, fmt.Errorf("videoID is required")
	}
	if input.Comment == "" {
		return entity.GofileVideoComment{}, fmt.Errorf("comment is required")
	}

	// 動画取得 & 閲覧可否（非公開は所有者のみ）
	video, err := u.gofileRepo.FindByID(input.VideoID)
	if err != nil {
		return entity.GofileVideoComment{}, err
	}
	if video.UserID != user.ID && !video.IsShared {
		return entity.GofileVideoComment{}, fmt.Errorf("%w: video not found", ErrKind.NotFound)
	}

	e := entity.GofileVideoComment{
		ID:            u.ulid.GenerateID(),
		GofileVideoID: input.VideoID,
		UserID:        user.ID,
		Comment:       input.Comment,
		LikeCount:     0,
		CreatedAt:     u.clock.Now(),
		UpdatedAt:     u.clock.Now(),
	}

	if err := u.gofileRepo.CreateComment(e); err != nil {
		return entity.GofileVideoComment{}, err
	}

	res, err := u.gofileRepo.FindCommentByID(e.ID)
	if err != nil {
		return entity.GofileVideoComment{}, err
	}

	return res, nil
}

func (u *GofileUseCase) CreateFromTwimgURL(user entity.User, srcURL string) (entity.GofileVideo, error) {
	if srcURL == "" {
		return entity.GofileVideo{}, fmt.Errorf("url is required")
	}

	// 1) video.twimg.com の直リンクのみ対応（tweet URL の解析は別レイヤでやる）
	uParsed, err := url.Parse(srcURL)
	if err != nil {
		return entity.GofileVideo{}, fmt.Errorf("invalid url: %w", err)
	}
	if !strings.EqualFold(uParsed.Host, "video.twimg.com") {
		return entity.GofileVideo{}, fmt.Errorf("unsupported host: %s (expect video.twimg.com)", uParsed.Host)
	}
	// ファイル名推定
	filename := path.Base(uParsed.Path)
	if filename == "" || filename == "." || filename == "/" {
		filename = "video.mp4"
	}
	if !strings.HasSuffix(strings.ToLower(filename), ".mp4") {
		filename += ".mp4"
	}

	// 2) Twitter CDN からストリーム取得
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	filename, dlResp, err := u.twitter.FetchTwimgStream(ctx, srcURL)
	if err != nil {
		return entity.GofileVideo{}, err
	}
	defer dlResp.Body.Close()

	// 3) Gofile へストリーム転送 (multipart/form-data)
	//    エンドポイントは東京。必要なら ENV で差し替え
	upData, err := u.gofileAPIDriver.Upload(ctx, filename, "", dlResp.Body)
	if err != nil {
		return entity.GofileVideo{}, fmt.Errorf("gofile upload: %w", err)
	}

	// ★ contentId 決定（idを優先、fileIdフォールバック）
	gofileID := upData.ID
	if gofileID == "" {
		gofileID = upData.FileID
	}
	if gofileID == "" {
		return entity.GofileVideo{}, fmt.Errorf("gofile upload ok but content id missing")
	}

	// 4) DirectLink を発行（必要に応じて）
	fmt.Println("-------------------------------start issuing direct link-------------------------------")
	issueRes, err := u.gofileAPIDriver.IssueDirectLink(gofileID, os.Getenv("GOFILE_API_KEY"))
	if err != nil {
		return entity.GofileVideo{}, fmt.Errorf("issue direct link: %w", err)
	}
	fmt.Println("-------------------------------\nend issuing direct link\n-------------------------------")
	if issueRes.DirectLink == "" {
		return entity.GofileVideo{}, fmt.Errorf("no directLink returned for fileId=%s", gofileID)
	}

	// 5) サムネ等メタ取得（既存どおり）
	fmt.Println("-------------------------------\nstart getting content\n-------------------------------")
	getRes, err := u.gofileAPIDriver.GetContent(gofileID, os.Getenv("GOFILE_API_KEY"))
	if err != nil {
		return entity.GofileVideo{}, fmt.Errorf("get content: %w", err)
	}
	fmt.Println("-------------------------------\nend getting content\n-------------------------------")

	// 6) DB登録（Create と揃える）
	proxiedVideoURL := os.Getenv("XROLL_API_ENDPOINT") + "/gofile/proxy?url=" + url.QueryEscape(issueRes.DirectLink)

	entityGofile := entity.GofileVideo{
		ID:              u.ulid.GenerateID(),
		Name:            filename, // デフォはファイル名。必要なら引数で渡してもよい
		GofileID:        gofileID,
		GofileDirectURL: issueRes.DirectLink,
		VideoURL:        proxiedVideoURL,
		ThumbnailURL:    getRes.Data.Thumbnail,
		Description:     "",
		IsShared:        false,
		UserID:          user.ID,
		GofileTags:      nil,
	}

	if err := u.gofileRepo.Create(entityGofile); err != nil {
		return entity.GofileVideo{}, err
	}

	// 7) レコード確認して返す
	res, err := u.gofileRepo.FindByID(entityGofile.ID)
	if err != nil {
		return entity.GofileVideo{}, err
	}

	// 8) （任意）ユーザーの GofileToken を保存したい場合はここで更新
	// 既存コードの条件は逆転してたので、空ならセットするのが自然
	if (user.GofileToken == nil || *user.GofileToken == "") && upData.ParentFolder != "" {
		// ここでは例として ParentFolder を token 的に扱うなら（実際は専用APIの token を保存するのが正）
		user.GofileToken = &upData.ParentFolder
		if err := u.userRepo.Update(user); err != nil {
			// token 保存失敗は致命でないのでログだけでもOK
			fmt.Printf("warn: failed to save user gofile token: %v\n", err)
		}
	}

	return res, nil
}
