package initdata

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/model"
)

// User returns []*model.User（例）
// このように、固定のマスタデータなどは初期投入データを関数で定義しておくと、jsonより楽な場合が多い
func User() []*model.User {
	return []*model.User{
		{
			ID:    "system-admin",
			Name: 	"システム管理者",
			Age: 0,
			UserType: "admin",
		},
		{
			ID:    "admin",
			Name: 	"管理者",
			Age: 0,
			UserType: "admin",
		},
		{
			ID:    "user",
			Name: 	"ユーザ",
			Age: 0,
			UserType: "user",
		},
		{
			ID:    "user",
			Name: 	"ユーザ",
			Age: 0,
			UserType: "user",
		},
	}
}
