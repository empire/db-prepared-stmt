package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
)

// Struct to hold prepared statements
type PreparedQueries struct {
	db   *sql.DB
	stmt *sql.Stmt
	lock sync.Mutex
}

func (pq *PreparedQueries) GetByID(id string) (string, error) {
	// Prepare statements lazily
	if err := pq.prepare(); err != nil {
		return "", err
	}

	rows, err := pq.stmt.Query(id)
	if err != nil {
		return "", fmt.Errorf("can get rows from db")
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("can't close rows", "err", err)
		}
	}()

	if !rows.Next() {
		return "", fmt.Errorf("no rows exists")
	}

	var name string
	return name, rows.Scan(&name)
}

func (pq *PreparedQueries) prepare() (err error) {
	if pq.stmt != nil {
		return nil
	}

	pq.lock.Lock()
	defer pq.lock.Unlock()

	if pq.stmt != nil {
		return nil
	}

	pq.stmt, err = pq.db.Prepare("SELECT name FROM users WHERE id = ?")
	return err
}

func (pq *PreparedQueries) Close() error {
	if pq.stmt == nil {
		return nil
	}
	return pq.stmt.Close()
}

func userhandle(pq *PreparedQueries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := "1"
		name, err := pq.GetByID(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusNotFound)
		}

		fmt.Fprintf(w, "User found: ID=%s, Name=%s\n", id, name)
	}
}

var monitored bool

func printStats(db *sql.DB) {
	if monitored {
		return
	}
	monitored = true
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		for {
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			fmt.Printf("Heap memory (in bytes): %d\n", memStats.HeapAlloc)
			spew.Dump(db.Stats())
			<-ticker.C
		}
	}()
}

func main() {
	db, err := sql.Open("mysql", "app:pass@tcp(localhost:6033)/mydb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(60 * time.Second)
	db.SetConnMaxIdleTime(60 * time.Second)
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(30)

	pq := &PreparedQueries{db: db}
	if err != nil {
		log.Fatal(err)
	}
	defer pq.Close()

	printStats(db)

	http.HandleFunc("/random", userhandle(pq))

	port := "7070"
	log.Printf("Server listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
