//go:generate mockgen -source=athena_query_repository.go -destination=../../mock/mock_athena_query_repository.go -package=mock
package repository

import "context"

type AthenaQueryRepository interface {
	Fetch(ctx context.Context, rawQuery string, args []interface{}) ([][]string, error)
}
