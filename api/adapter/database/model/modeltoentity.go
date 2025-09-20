package model

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type model[T any] interface {
	Entity() T
}

// ToEntities converts slice of models to slice of entities.
func ToEntities[Entity any, Model model[Entity]](models []Model) []Entity {
	ret := make([]Entity, len(models))
	for i, model := range models {
		ret[i] = model.Entity()
	}
	return ret
}

func (u User) Entity() entity.User {
	return entity.User{
		ID:             u.ID,
		Name:           u.Name,
		Age:            u.Age,
		UserType:       entconst.UserType(u.UserType),
		Email:          u.Email,
		Bio:            u.Bio,
		HashedPassword: u.HashedPassword,
		GofileToken:    u.GofileToken,
		EmailVerified:  u.EmailVerified,
		IsDeleted:      u.IsDeleted,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

func (v Video) Entity() entity.Video {
	return entity.Video{
		ID:            v.ID,
		Ranking:       v.Ranking,
		ThumbnailURL:  v.ThumbnailURL,
		VideoURL:      v.VideoURL,
		DownloadCount: v.DownloadCount,
		LikeCount:     v.LikeCount,
		TweetURL:      v.TweetURL,
		Comments:      ToEntities(v.Comments),
		CreatedAt:     v.CreatedAt,
	}
}

func (c VideoComment) Entity() entity.Comment {
	return entity.Comment{
		ID:        c.ID,
		Comment:   c.Comment,
		LikeCount: c.LikeCount,
		CreatedAt: c.CreatedAt,
	}
}

func (r RegisterVerification) Entity() entity.RegisterVerification {
	return entity.RegisterVerification{
		Email:                    r.Email,
		ExpiresAt:                r.ExpiresAt,
		HashedPassword:           r.HashedPassword,
		HashedAuthenticationCode: r.HashedAuthenticationCode,
	}
}

func (g GofileVideo) Entity() entity.GofileVideo {
	return entity.GofileVideo{
		ID:                  g.ID,
		Name:                g.Name,
		GofileID:            g.GofileID,
		GofileDirectURL:     g.GofileDirectURL,
		VideoURL:            g.VideoURL,
		ThumbnailURL:        g.ThumbnailURL,
		PlayCount:           g.PlayCount,
		Description:         g.Description,
		LikeCount:           g.LikeCount,
		IsShared:            g.IsShared,
		UserID:              g.UserID,
		User:                g.User.Entity(),
		GofileTags:          ToEntities(g.GofileTags),
		CreatedAt:           g.CreatedAt,
		UpdatedAt:           g.UpdatedAt,
		GofileVideoComments: ToEntities(g.GofileVideoComments),
		IsDeleted:           g.IsDeleted,
	}
}

func (gt GofileTag) Entity() entity.GofileTag {
	return entity.GofileTag{
		ID:   gt.ID,
		Name: gt.Name,
	}
}
func (gvc GofileVideoComment) Entity() entity.GofileVideoComment {
	return entity.GofileVideoComment{
		ID:            gvc.ID,
		GofileVideoID: gvc.GofileVideoID,
		UserID:        gvc.UserID,
		User:          gvc.User.Entity(),
		Comment:       gvc.Comment,
		LikeCount:     gvc.LikeCount,
		CreatedAt:     gvc.CreatedAt,
		UpdatedAt:     gvc.UpdatedAt,
	}
}
