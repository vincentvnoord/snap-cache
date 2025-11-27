package database

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

func ConnectDB(ctx context.Context) (*pgx.Conn, error) {
	dbURL := "postgres://bench:bench@localhost:5432/benchdb?sslmode=disable"

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	err = resetDB(ctx, conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func resetDB(ctx context.Context, conn *pgx.Conn) error {
	path, _ := filepath.Abs("migrate.sql")
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, string(sqlBytes))
	return err
}

func Seed(ctx context.Context, conn *pgx.Conn) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		panic(err)
	}

	for userID := 1; userID <= 100_000; userID++ {
		// Insert a user
		_, err := tx.Exec(ctx, "INSERT INTO users(name,email) VALUES($1,$2)", strconv.Itoa(userID)+"name", strconv.Itoa(userID)+"@mail.com")
		if err != nil {
			panic(err)
		}

		// Insert some orders for this user
		for i := 0; i < rand.Intn(5); i++ { // 0-4 orders per user
			_, err := tx.Exec(ctx, "INSERT INTO orders(user_id,amount,created_at) VALUES($1,$2,NOW())",
				userID, rand.Intn(1000))
			if err != nil {
				panic(err)
			}
		}
		fmt.Printf("Inserted user %d\r", userID)
	}

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}
}
