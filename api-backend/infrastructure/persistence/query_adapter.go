package persistence

import (
	"encoding/json"

	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type QueryAdapter struct {
	conn *gorm.DB
}

func NewQueryAdapter(conn *gorm.DB) repository.QueryRepository {
	return &QueryAdapter{
		conn,
	}
}

func (qa *QueryAdapter) Fetch(rawQuery string, args []interface{}, scanType interface{}) (interface{}, error) {
	_, err := qa.conn.Raw(rawQuery, args...).Rows()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch query: %+v, %+v", &rawQuery, args)
	}

	qa.conn.Raw(rawQuery, args...).Scan(&scanType)

	return scanType, nil
}

func deepcopyJson(src interface{}, dst interface{}) (err error) {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, dst)
	if err != nil {
		return err
	}
	return nil
}
