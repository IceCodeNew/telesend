package lotusdb

import (
	"github.com/lotusdblabs/lotusdb/v2"
)

func InitDB(dbPath string) (*lotusdb.DB, error) {
	options := lotusdb.DefaultOptions
	options.DirPath = dbPath

	db, err := lotusdb.Open(options)
	if err != nil {
		return nil, err
	}
	return db, nil
}
