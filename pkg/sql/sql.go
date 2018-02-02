package sql

import (
	"database/sql"
	"fmt"
	"time"

	// import mysql driver(don't delete me)
	_ "github.com/go-sql-driver/mysql"
	"github.com/juju/errors"
	"github.com/ngaut/log"
)

var (
	MaxDMLRetryCount = 100
	MaxDDLRetryCount = 5
	RetryWaitTime    = 3 * time.Second
)

// ExecuteSQLs execute sqls in a transaction with retry.
func ExecuteSQLs(db *sql.DB, sqls []string, args [][]interface{}, isDDL bool) error {
	if len(sqls) == 0 {
		return nil
	}

	retryCount := MaxDMLRetryCount
	if isDDL {
		retryCount = MaxDDLRetryCount
	}

	var err error
	for i := 0; i < retryCount; i++ {
		if i > 0 {
			log.Warnf("exec sql retry %d - %v", i, sqls)
			time.Sleep(RetryWaitTime)
		}

		err = appleTxn(db, sqls, args)
		if err == nil {
			return nil
		}
	}

	return errors.Trace(err)
}

func appleTxn(db *sql.DB, sqls []string, args [][]interface{}) error {
	txn, err := db.Begin()
	if err != nil {
		log.Errorf("exec sqls[%v] begin failed %v", sqls, errors.ErrorStack(err))
		return errors.Trace(err)
	}

	for i := range sqls {
		log.Debugf("[exec][sql]%s[args]%v", sqls[i], args[i])

		_, err = txn.Exec(sqls[i], args[i]...)
		if err != nil {
			log.Warnf("[exec][sql]%s[args]%v[error]%v", sqls[i], args[i], err)
			rerr := txn.Rollback()
			if rerr != nil {
				log.Errorf("[rollback][error]%v", rerr)
			}
			return errors.Trace(err)
		}
	}

	err = txn.Commit()
	if err != nil {
		log.Errorf("exec sqls[%v] commit failed %v", sqls, errors.ErrorStack(err))
		return errors.Trace(err)
	}

	return nil
}

// OpenDB creates an instance of sql.DB.
func OpenDB(proto string, host string, port int, username string, password string) (*sql.DB, error) {
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&multiStatements=true", username, password, host, port)
	db, err := sql.Open(proto, dbDSN)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return db, nil
}
