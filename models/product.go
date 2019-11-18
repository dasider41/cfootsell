package models

import (
	"github.com/dasider41/cfootsell/db"
	"github.com/dasider41/cfootsell/util"
)

// Product : product model
type Product struct {
	Title     string
	Condition string
	Price     int
	Member    string
	Size      int
	Updated   string
}

// IsAvailable :
func (p Product) IsAvailable() bool {
	if len(p.Member) > 0 {
		return true
	}
	return false
}

// UpdateItemStatus :
func UpdateItemStatus(arr []int) {
	conn := db.InitDB()
	defer conn.Close()

	stmUpdateQuery, err := conn.Prepare(
		"UPDATE market SET " +
			"`isNew` = ? " +
			" WHERE id = ?")
	util.ErrCheck(err)
	defer stmUpdateQuery.Close()

	for _, id := range arr {
		_, err := stmUpdateQuery.Exec(false, id)
		util.ErrCheck(err)
	}
}

// UpdateTransaction :
func UpdateTransaction(e Product) (int64, error) {
	conn := db.InitDB()
	defer conn.Close()

	stmInsertQuery, err := conn.Prepare("INSERT IGNORE INTO market SET " +
		"`title`=?," +
		"`condition`=?," +
		"`size`=?," +
		"`price`=?," +
		"`member`=?," +
		"`updated`=?")

	if err != nil {
		return 0, err
	}

	defer stmInsertQuery.Close()

	res, err := stmInsertQuery.Exec(e.Title, e.Condition, e.Size, e.Price, e.Member, e.Updated)
	util.ErrCheck(err)
	count, err := res.RowsAffected()
	return count, nil
}
