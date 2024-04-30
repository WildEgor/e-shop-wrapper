package repositories

import (
	"context"
	"github.com/WildEgor/e-shop-fiber-wrapper/internal/db/clickhouse"
)

type RecordsRepository struct {
	client *clickhouse.ClickhouseConnection
}

func NewRecordsRepository(
	client *clickhouse.ClickhouseConnection,
) *RecordsRepository {
	return &RecordsRepository{
		client,
	}
}

func (rr *RecordsRepository) GetRecords(ctx context.Context, sql string) ([]map[string]interface{}, error) {
	rows, err := rr.client.QueryWithTimeout(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var dict []map[string]interface{}

	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))

		for i, _ := range vals {
			ptrs[i] = &vals[i]
		}

		err := rows.Scan(ptrs...)
		if err != nil {
			return nil, err
		}

		data := make(map[string]interface{})

		for i, val := range vals {
			data[cols[i]] = val
		}

		dict = append(dict, data)
	}

	return dict, nil
}
