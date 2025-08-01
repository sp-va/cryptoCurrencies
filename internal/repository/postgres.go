package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	customErrors "github.com/sp-va/cryptoCurrencies/internal/custom-errors"
	"github.com/sp-va/cryptoCurrencies/internal/dto"
	"github.com/sp-va/cryptoCurrencies/internal/models"
)

type PostgresConnectionConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func (c *PostgresConnectionConfig) PostgresConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}

func (c *PostgresConnectionConfig) ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", c.PostgresConnectionString())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *PostgresConnectionConfig) RunMigrations() {
	m, err := migrate.New(
		"file://internal/repository/migrations",
		c.PostgresConnectionString(),
	)
	if err != nil {
		log.Fatalf("Migration init error: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration run error: %v", err)
	}

	log.Println("Migrations applied successfully.")
}

type CurrencyRepository interface {
	InsertCurrencyToTrack(ctx context.Context, c dto.AddCurrency) error
	InsertCoinData(c *models.Currency) error
	GetCoinsToTrack() ([]string, error)
	DeleteCurrencyFromTracking(ctx context.Context, coin string) (int64, error)
	GetCoinValue(ctx context.Context, coin string, timestamp uint32) (*models.Currency, error)
}

type PostgresCurrencyRepo struct {
	Db *sql.DB
}

func InitDB() (*PostgresCurrencyRepo, error) {
	config := &PostgresConnectionConfig{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Database: os.Getenv("POSTGRES_DB"),
	}

	config.RunMigrations()

	db, err := config.ConnectToDB()
	if err != nil {
		return nil, err
	}

	return NewPostgresCurrencyRepository(db), nil
}

func NewPostgresCurrencyRepository(db *sql.DB) *PostgresCurrencyRepo {
	return &PostgresCurrencyRepo{Db: db}
}

func (r *PostgresCurrencyRepo) InsertCurrencyToTrack(ctx context.Context, c dto.AddCurrency) error {
	_, err := r.Db.Exec("INSERT INTO track_currencies (coin) VALUES ($1)", c.Coin)

	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23505":
				return customErrors.CurrencyAlreadyInserted
			default:
				return err
			}
		}
	}

	return err
}

func (r *PostgresCurrencyRepo) InsertCoinData(c *models.Currency) error {
	_, err := r.Db.Exec("INSERT INTO currency_prices (coin, timestamp, price) VALUES ($1, $2, $3)", c.Coin, c.Timestamp, c.Price)
	return err
}

func (r *PostgresCurrencyRepo) DeleteCurrencyFromTracking(ctx context.Context, coin string) (int64, error) {
	res, err := r.Db.Exec("DELETE FROM track_currencies WHERE coin = $1", coin)

	if err != nil {
		return 0, err
	}
	rAffected, _ := res.RowsAffected()
	return rAffected, nil
}

func (r *PostgresCurrencyRepo) GetCoinsToTrack() ([]string, error) {
	rows, err := r.Db.Query("SELECT coin FROM track_currencies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coins []string
	for rows.Next() {
		var coin string
		if err := rows.Scan(&coin); err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	return coins, rows.Err()
}

func (r *PostgresCurrencyRepo) GetCoinValue(ctx context.Context, coin string, timestamp uint32) (*models.Currency, error) {
	q := "SELECT coin, timestamp, price FROM currency_prices Where coin = $1 ORDER BY ABS(timestamp-$2) LIMIT 1;"

	var currency models.Currency

	err := r.Db.QueryRow(q, coin, timestamp).Scan(&currency.Coin, &currency.Timestamp, &currency.Price)
	if err != nil {
		log.Printf("Полученное значение стоимости: %v", err)
		return nil, err
	}
	return &currency, nil
}
