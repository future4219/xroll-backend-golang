package main

import (
	"fmt"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/initdata"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/log"
)

func main() {
	logger, err := log.NewLogger()
	if err != nil {
		return
	}

	db, err := database.NewMySQLDB(logger, true)
	if err != nil {
		fmt.Println("error:", err)
	}

	// マイグレーション
	err = database.Migrate(db)
	if err != nil {
		fmt.Println("error:", err)
	}

	// 開発用テストデータ作成
	err = initdata.CreateStgSeedData(db)
	if err != nil {
		fmt.Println("error:", err)
	}
}
