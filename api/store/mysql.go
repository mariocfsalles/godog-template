package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Order struct {
	ID        string    `json:"id"`
	Customer  string    `json:"customer"`
	Status    string    `json:"status"`
	Items     []string  `json:"items"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MustMySQL abre conexão com MySQL, espera o DB ficar pronto e
// opcionalmente derruba/cria a tabela orders (reset=true).
func MustMySQL(dsn string, reset bool) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}

	// espera o MySQL ficar pronto
	deadline := time.Now().Add(60 * time.Second)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			log.Fatalf("mysql ping timeout: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	if reset {
		if _, err := db.Exec(`DROP TABLE IF EXISTS orders`); err != nil {
			log.Fatalf("drop table orders: %v", err)
		}

		ddl := `
		CREATE TABLE orders (
			id         CHAR(26)     PRIMARY KEY,
			customer   VARCHAR(255) NOT NULL,
			status     VARCHAR(32)  NOT NULL,
			items_json JSON         NOT NULL,
			created_at DATETIME(6)  NOT NULL,
			updated_at DATETIME(6)  NOT NULL,
			KEY idx_status (status),
			KEY idx_created (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`
		if _, err := db.Exec(ddl); err != nil {
			log.Fatalf("apply ddl: %v", err)
		}
	}

	return db
}

func ScanOrder(row *sql.Row) (*Order, error) {
	var (
		o         Order
		itemsJSON []byte
	)
	err := row.Scan(&o.ID, &o.Customer, &o.Status, &itemsJSON, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(itemsJSON, &o.Items)
	return &o, nil
}

func ScanOrders(rows *sql.Rows) ([]Order, error) {
	var out []Order
	for rows.Next() {
		var (
			o         Order
			itemsJSON []byte
		)
		if err := rows.Scan(&o.ID, &o.Customer, &o.Status, &itemsJSON, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(itemsJSON, &o.Items)
		out = append(out, o)
	}
	return out, rows.Err()
}

// Helpers para distinguir erro de "não encontrado" se quiser reutilizar
var ErrNotFound = errors.New("not found")

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, ErrNotFound)
}
