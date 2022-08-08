package persistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"
	"github.com/pkg/errors"
)

type AthenaQueryAdapter struct {
	database         string
	s3OutputLocation string
	*athena.Client
}

func NewAthenaQueryAdapter(ctx context.Context, region, database, s3OutputLocation string) (repository.AthenaQueryRepository, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure aws")
	}

	return &AthenaQueryAdapter{
		database,
		s3OutputLocation,
		athena.NewFromConfig(cfg),
	}, nil
}

func (aa *AthenaQueryAdapter) Fetch(ctx context.Context, rawQuery string, args []interface{}) ([][]string, error) {
	query := fmt.Sprintf(rawQuery, args...)

	startQueryExecutionInput := &athena.StartQueryExecutionInput{
		QueryExecutionContext: &types.QueryExecutionContext{
			Database: aws.String(aa.database),
		},
		ResultConfiguration: &types.ResultConfiguration{
			OutputLocation: aws.String(aa.s3OutputLocation),
		},
		QueryString: aws.String(query),
	}

	queryExecutionOutput, err := aa.StartQueryExecution(ctx, startQueryExecutionInput)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to start query: %+v, %+v", &rawQuery, args)
	}

	id := queryExecutionOutput.QueryExecutionId
	getQueryExecutionInput := &athena.GetQueryExecutionInput{
		QueryExecutionId: id,
	}

	err = aa.waitGetQueryExecution(ctx, getQueryExecutionInput)
	if err != nil {
		log.Fatal(err)
	}

	getQueryResultsInput := &athena.GetQueryResultsInput{
		QueryExecutionId: id,
	}

	output, err := aa.GetQueryResults(ctx, getQueryResultsInput)

	length := len(output.ResultSet.Rows)
	result := make([][]string, 0, length-1)

	for _, row := range output.ResultSet.Rows[1:length] {
		var array []string
		for _, val := range row.Data {
			array = append(array, *val.VarCharValue)
		}

		result = append(result, array)
	}

	return result, nil
}

func (aa *AthenaQueryAdapter) waitGetQueryExecution(ctx context.Context, input *athena.GetQueryExecutionInput) error {
	for {
		gqeOutput, err := aa.GetQueryExecution(ctx, input)
		if err != nil {
			return errors.Wrapf(err, "failed to get query execution")
		}

		switch *&gqeOutput.QueryExecution.Status.State {
		case types.QueryExecutionStateSucceeded:
			return nil
		case types.QueryExecutionStateFailed, types.QueryExecutionStateCancelled:
			return errors.New("QueryExecutionStateFailed or QueryExecutionStateCancelled")
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
