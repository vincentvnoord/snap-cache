package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/vincentvnoord/snap-cache/cmd/benchmark/database/database"
	"github.com/vincentvnoord/snap-cache/internal/client"
)

func main() {
	ctx := context.Background()
	conn, err := database.ConnectDB(ctx)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	fmt.Println("Connected to DB successfully")

	database.Seed(ctx, conn)

	benchmarkSelect(conn, "SELECT * FROM users WHERE id = 12345", 10000)
	dbSums := benchmarkUserOrdersSum(conn, "SELECT u.id, SUM(o.amount) FROM users u JOIN orders o ON u.id = o.user_id GROUP BY u.id;", 100)
	fmt.Printf("Fetched %d order sums from database\n", len(dbSums))

	client, err := client.NewClient("localhost:8080")
	defer client.Close()

	benchmarkCacheSet(client, dbSums)
	cachedSums := benchmarkCacheGet(client, dbSums)

	fmt.Printf("Fetched %d order sums from cache\n", len(cachedSums))
}

type UserSum struct {
	UserID int
	Sum    int
}

func benchmarkCacheGet(client *client.Client, userSums []UserSum) []UserSum {
	var min, max, total time.Duration
	min = time.Hour
	results := make([]UserSum, 0)

	for _, sum := range userSums {
		start := time.Now()
		data, err := client.SendCommand("GET", strconv.Itoa(sum.UserID), nil)
		if err != nil {
			panic(err)
		}

		s, err := strconv.Atoi(string(data))
		if err != nil {
			panic(err)
		}

		results = append(results, UserSum{
			UserID: sum.UserID,
			Sum:    s,
		})

		elapsed := time.Since(start)

		if elapsed < min {
			min = elapsed
		}
		if elapsed > max {
			max = elapsed
		}
		total += elapsed
	}

	printTimings("CACHE: GET user order sum", min, max, total, len(userSums))
	return results
}

func benchmarkCacheSet(client *client.Client, userSums []UserSum) {
	var min, max, total time.Duration
	min = time.Hour

	for _, sum := range userSums {
		start := time.Now()
		_, err := client.SendCommand("SET", strconv.Itoa(sum.UserID), []byte(strconv.Itoa(sum.Sum)))

		if err != nil {
			panic(err)
		}
		elapsed := time.Since(start)

		if elapsed < min {
			min = elapsed
		}
		if elapsed > max {
			max = elapsed
		}
		total += elapsed
	}

	printTimings("CACHE: SET user order sum", min, max, total, len(userSums))
}

func benchmarkUserOrdersSum(conn *pgx.Conn, query string, iterations int) []UserSum {
	var min, max, total time.Duration
	min = time.Hour
	results := make([]UserSum, 0)

	for i := 0; i < iterations; i++ {
		results = results[:0] // reset results slice

		start := time.Now()
		rows, err := conn.Query(context.Background(), query)
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			var sum UserSum
			err := rows.Scan(&sum.UserID, &sum.Sum)
			if err != nil {
				panic(err)
			}
			results = append(results, sum)
		}

		elapsed := time.Since(start)

		if elapsed < min {
			min = elapsed
		}
		if elapsed > max {
			max = elapsed
		}
		total += elapsed
	}

	printTimings("DB: SELECT User order SUM", min, max, total, iterations)
	return results
}

func benchmarkSelect(conn *pgx.Conn, query string, iterations int) {
	var min, max, total time.Duration
	min = time.Hour

	for i := 0; i < iterations; i++ {
		start := time.Now()
		row := conn.QueryRow(context.Background(), query)
		var id int
		var name string
		row.Scan(&id, &name) // read one column or all columns if needed
		elapsed := time.Since(start)

		if elapsed < min {
			min = elapsed
		}
		if elapsed > max {
			max = elapsed
		}
		total += elapsed
	}

	printTimings("DB: SELECT Username + ID", min, max, total, iterations)
}

func printTimings(title string, min, max, total time.Duration, iterations int) {
	fmt.Printf("%s - min: %v, max: %v, avg: %v, total: %v\n",
		title, min, max, total/time.Duration(iterations), total)
}
