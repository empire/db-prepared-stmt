package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
)

// Struct to hold prepared statements
type PreparedQueries struct {
	db   *sql.DB
	stmt *sql.Stmt
}

func (pq *PreparedQueries) GetByID(id string) (string, error) {
	var i atomic.Int32
	// var lock sync.Mutex

	if pq.stmt == nil {
		// lock.Lock()
		// defer lock.Unlock()
		if pq.stmt == nil {
			var err error
			pq.stmt, err = pq.db.Prepare("SELECT name FROM users WHERE id = ?")
			if err != nil {
				return "", err
			}
			i.Add(1)
			fmt.Println("prepare the query", i.Load())
		}
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

func main() {
	// Connect to the database
	db, err := sql.Open("mysql", "app:pass@tcp(localhost:6033)/mydb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(3 * time.Second)
	db.SetMaxOpenConns(300)
	db.SetMaxIdleConns(300)
	db.SetConnMaxIdleTime(3 * time.Second)

	// Prepare SQL statements
	pq := &PreparedQueries{db: db}
	if err != nil {
		log.Fatal(err)
	}
	defer pq.Close()

	if true {
		ctx := context.Background()
		stm, err := db.PrepareContext(ctx, "do sleep(?)")
		if err != nil {
			panic(err)
		}

		stm2, err := db.PrepareContext(ctx, "do sleep(?)")
		if err != nil {
			panic(err)
		}

		fetch := func() {
			rows, err := stm.Query("1")
			if err != nil {
				panic(err)
			}
			defer rows.Close()
			for rows.Next() {
				var name string
				rows.Scan(&name)
				fmt.Printf("name='%s'\n", name)
			}
			if rows.Err() != nil {
				panic(rows.Err())
			}
		}

		fetch2 := func() {
			var name string
			rows, err := stm2.Query("1")
			if err != nil {
				panic(err)
			}
			defer rows.Close()
			for rows.Next() {
				rows.Scan(&name)
				fmt.Printf("name='%s'\n", name)
			}
			if rows.Err() != nil {
				panic(rows.Err())
			}
		}
		// for i := 0; i < 16; i++ {
		// 	// go fetch()
		// 	fetch()
		// }

		time.Sleep(10 * time.Millisecond)
		fetch()
		fetch() // We need to restart proxysql and start another mysql server
		fetch2()
		fetch()
		time.Sleep(10 * time.Millisecond)
		spew.Dump(db.Stats())
		fetch2()
		// It reports (2+1) * 10s wait duration
	}

	fmt.Println(pq.GetByID("1"))

	if true {
		go func() {
			ticker := time.NewTicker(3 * time.Second)
			for {
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)

				fmt.Printf("Heap memory (in bytes): %d\n", memStats.HeapAlloc)
				spew.Dump(db.Stats())
				<-ticker.C
			}
		}()
	}

	// HTTP handler function
	http.HandleFunc("/random", userhandle(pq))
	http.HandleFunc("/gc", func(w http.ResponseWriter, r *http.Request) {
		runtime.GC()

		fmt.Fprintf(w, "done")
	})

	// Start the HTTP server
	log.Println("Server listening on port 7070...")
	log.Fatal(http.ListenAndServe(":7070", nil))
}
