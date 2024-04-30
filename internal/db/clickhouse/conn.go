package clickhouse

import (
	"context"
	"database/sql"
	"github.com/ClickHouse/clickhouse-go/v2"

	"github.com/WildEgor/e-shop-fiber-wrapper/internal/configs"
	"github.com/gofiber/fiber/v3/log"
	"time"
)

type ClickhouseConnection struct {
	conn *sql.DB
	cfg  *configs.ClickhouseConfig
}

func NewClickhouseConnection(
	cfg *configs.ClickhouseConfig,
) *ClickhouseConnection {
	conn := &ClickhouseConnection{
		conn: nil,
		cfg:  cfg,
	}

	conn.Connect()

	return conn
}

func (cc *ClickhouseConnection) Connect() {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{cc.cfg.DSN},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
	})

	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)
	cc.conn = conn

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cc.Ping(ctx); err != nil {
		log.Panic("Fail connect Clickhouse", err)
		return
	}

	log.Info("Connection to Clickhouse established!")

}

func (cc *ClickhouseConnection) Ping(ctx context.Context) error {
	_, err := cc.conn.ExecContext(ctx, "SELECT 1;")
	return err
}

func (cc *ClickhouseConnection) Disconnect() {
	if err := cc.conn.Close(); err != nil {
		log.Panic("Fail disconnect Clickhouse", err)
		return
	}

	log.Info("Connection to Clickhouse closed.")
}

func (cc *ClickhouseConnection) QueryWithTimeout(ctx context.Context, sql string) (*sql.Rows, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query, err := cc.conn.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	return query, nil
}
