package db

import (
	"fmt"

	"github.com/IceCodeNew/telesend/internal/app/config"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var (
	engine           *xorm.Engine
	engineNotInitErr = fmt.Errorf("ERROR: [Internal] Engine not initialized")
)

func initEngine(dbPath string) error {
	var err error
	engine, err = xorm.NewEngine("sqlite", dbPath)
	if err != nil {
		return err
	}

	if config.TSConfig.Verbose {
		engine.Logger().SetLevel(log.LOG_DEBUG)
	}
	return nil
}

func Close() {
	_ = engine.Close()
}

func CreateTable[T any](t *T) error {
	if err := initEngine(config.TSConfig.DBPath); err != nil {
		return err
	}

	if exist, err := engine.IsTableExist(t); err != nil {
		return err
	} else if !exist {
		return engine.CreateTables(t)
	}
	return nil
}

func SetElem[T any](elem *T) error {
	if engine == nil {
		return engineNotInitErr
	}

	_, err := engine.Transaction(
		func(session *xorm.Session) (interface{}, error) {
			if _, err := session.Insert(elem); err != nil {
				return nil, err
			}
			return nil, nil
		})
	if err != nil {
		return err
	}
	return nil
}

func GetElemByID[T any](id string) (*T, error) {
	if engine == nil {
		return nil, engineNotInitErr
	}

	var elem *T
	if found, err := engine.ID(id).Get(elem); err != nil {
		return nil, err
	} else if !found {
		elem = nil
		return nil, fmt.Errorf("ERROR: [Internal] Element %s: not found", id)
	}
	return elem, nil
}
