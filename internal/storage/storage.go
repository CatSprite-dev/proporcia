package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/CatSprite-dev/proporcia/internal/domain"
	"github.com/shopspring/decimal"
	_ "modernc.org/sqlite"
)

type Storage struct {
	conn   *sql.DB
	logger *slog.Logger
}

func NewStorage(path string, logger *slog.Logger) (*Storage, error) {
	val := url.Values{}
	val.Add("_journal_mode", "WAL")
	val.Add("_busy_timeout", "5000")
	val.Add("_foreign_keys", "ON")

	dsn := path + "?" + val.Encode()

	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	conn.SetMaxOpenConns(1)

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	return &Storage{conn: conn, logger: logger}, nil
}

func (s *Storage) Init(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS targets (
		id         INTEGER PRIMARY KEY,
		name       TEXT NOT NULL,
		ticker     TEXT NOT NULL UNIQUE,
		weight     TEXT NOT NULL,
		uid        TEXT NOT NULL UNIQUE,
		class_code TEXT NOT NULL,
		lot        INTEGER NOT NULL,
		type       TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS orders (
		order_id               TEXT PRIMARY KEY,
		order_request_id       TEXT NOT NULL,
		class_code             TEXT NOT NULL,
		ticker                 TEXT NOT NULL,
		figi                   TEXT NOT NULL,
		instrument_uid         TEXT NOT NULL,
		lots_requested         TEXT NOT NULL,
		lots_executed          TEXT NOT NULL,
		initial_order_price    TEXT NOT NULL,
		initial_commission     TEXT NOT NULL,
		total_order_amount     TEXT NOT NULL,
		executed_order_price   TEXT NOT NULL,
		executed_commission    TEXT NOT NULL,
		initial_security_price TEXT NOT NULL,
		currency               TEXT NOT NULL,
		message                TEXT NOT NULL,
		server_time            TIMESTAMP NOT NULL,
		tracking_id            TEXT NOT NULL,
		created_at             TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := s.conn.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("init tables: %w", err)
	}

	s.logger.Info("storage initialized", "tables", []string{"targets", "orders"})
	return nil
}

func (s *Storage) UpsertTarget(ctx context.Context, target domain.Target) error {
	query := `INSERT INTO targets (name, ticker, weight, uid, class_code, lot, type)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(ticker) DO UPDATE SET
		weight = excluded.weight,
		updated_at = CURRENT_TIMESTAMP;`

	if _, err := s.conn.ExecContext(ctx, query, target.Name, target.Ticker, target.Weight.String(), target.InstrumentUID, target.ClassCode, target.Lot, target.Type); err != nil {
		return fmt.Errorf("upsert target: %w", err)
	}
	s.logger.Debug("target upserted", "target", target.Ticker)
	return nil
}

func (s *Storage) GetTargets(ctx context.Context) ([]domain.Target, error) {
	query := `SELECT id, name, ticker, weight, uid, class_code, lot, type FROM targets;`
	rows, err := s.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get targets: %w", err)
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		var weightStr string
		if err := rows.Scan(&target.ID, &target.Name, &target.Ticker, &weightStr, &target.InstrumentUID, &target.ClassCode, &target.Lot, &target.Type); err != nil {
			return nil, fmt.Errorf("scan target: %w", err)
		}
		target.Weight, err = decimal.NewFromString(weightStr)
		if err != nil {
			return nil, fmt.Errorf("parse weight: %w", err)
		}
		targets = append(targets, target)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate targets: %w", err)
	}

	s.logger.Debug("targets retrieved", "count", len(targets))
	return targets, nil
}

func (s *Storage) DeleteTarget(ctx context.Context, ticker string) error {
	query := `DELETE FROM targets WHERE ticker = ?;`

	if _, err := s.conn.ExecContext(ctx, query, ticker); err != nil {
		return fmt.Errorf("delete target: %w", err)
	}

	s.logger.Info("target deleted", "ticker", ticker)
	return nil
}

func (s *Storage) SaveOrders(ctx context.Context, orders []domain.PostOrderResponse) error {
	query := `INSERT INTO orders (
        order_id,
        order_request_id,
        class_code,
        ticker,
        figi,
        instrument_uid,
        lots_requested,
        lots_executed,
        initial_order_price,
        initial_commission,
        total_order_amount,
        executed_order_price,
        executed_commission,
        initial_security_price,
        currency,
        message,
        server_time,
        tracking_id
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    ON CONFLICT(order_id) DO UPDATE SET
        lots_executed = excluded.lots_executed,
        executed_order_price = excluded.executed_order_price,
        executed_commission = excluded.executed_commission,
        total_order_amount = excluded.total_order_amount,
        message = excluded.message;`

	for _, order := range orders {
		currency := order.InitialOrderPrice.Currency
		if currency == "" {
			currency = order.InitialSecurityPrice.Currency
		}

		_, err := s.conn.ExecContext(ctx, query,
			order.OrderID,
			order.OrderRequestID,
			order.ClassCode,
			order.Ticker,
			order.Figi,
			order.InstrumentUID,
			order.LotsRequested,
			order.LotsExecuted,
			order.InitialOrderPrice.Amount.String(),
			order.InitialCommission.Amount.String(),
			order.TotalOrderAmount.Amount.String(),
			order.ExecutedOrderPrice.Amount.String(),
			order.ExecutedCommission.Amount.String(),
			order.InitialSecurityPrice.Amount.String(),
			currency,
			order.Message,
			order.ResponseMetadata.ServerTime,
			order.ResponseMetadata.TrackingID,
		)
		if err != nil {
			return fmt.Errorf("save order %s: %w", order.OrderID, err)
		}
	}

	s.logger.Info("orders saved", "count", len(orders))
	return nil
}

func (s *Storage) Close() error {
	return s.conn.Close()
}
