package db

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type DBManager struct {
	connURL string
	conn    *pgx.Conn
	Queries *Queries
	logger  *logrus.Logger
}

func NewDBManager(connUrl string, logger *logrus.Logger) *DBManager {
	return &DBManager{
		connURL: connUrl,
		logger:  logger,
	}
}

func (m *DBManager) Connect(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, m.connURL)
	if err != nil {
		return fmt.Errorf("DB connection failed: %w", err)
	}

	m.conn = conn
	m.Queries = New(conn)
	return nil
}

func (m *DBManager) ConnectWithRetry(ctx context.Context, maxAttempts int) error {
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if attempt > 1 {
			wait := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
			m.logger.Infof("Attempt %d-%d: reconnecting to DB in %v", attempt, maxAttempts, wait)
			time.Sleep(wait)
		}

		if err = m.Connect(ctx); err == nil {
			m.logger.Info("Successfully connected to database")
			return nil
		}
		m.logger.Warnf("Connection attempt failed: %v", err)
	}
	return fmt.Errorf("after %d attempts: %w", maxAttempts, err)
}

func (m *DBManager) StartHealthCheck(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := m.Ping(ctx); err != nil {
					m.logger.Warnf("DB health check failed: %v", err)
					if err := m.ConnectWithRetry(ctx, 3); err != nil {
						m.logger.Errorf("Reconnection failed: %v", err)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (m *DBManager) Ping(ctx context.Context) error {
	if m.conn == nil {
		return fmt.Errorf("no active connection")
	}
	return m.conn.Ping(ctx)
}

func (m *DBManager) Close() error {
	if m.conn == nil {
		return nil
	}
	return m.conn.Close(context.Background())
}
