//go:generate mockgen -source=query_repository.go -destination=../../mock/mock_query_repository.go -package=mock
package repository

type QueryRepository interface {
	Fetch(rawQuery string, args []interface{}, scanType interface{}) (interface{}, error)
}
