package model

import (
	"fmt"
	"strconv"
)

// LoginIDがInsert/Updateされる時にintに変換できるか確認する
func (u User) validateLoginID(loginID string) error {
	_, err := strconv.Atoi(loginID)
	if err != nil || len(loginID) != 6 {
		return fmt.Errorf(
			"attempt to insert/update invalid loginID, should be 6 digit numeric (%s): %w",
			loginID, err)
	}
	return nil
}
