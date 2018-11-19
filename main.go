package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

const (
	DbHost = "db"
	DbUser = "postgres-dev"
	DbPassword = "mysecretpassword"
	DbName = "dev"
	Migration = `CREATE TABLE IF NOT EXISTS bulletins (
id serial PRIMARY KEY,
author text NOT NULL,
content text NOT NULL,
created_at timestamp with time zone DEFAULT current_timestamp)
`
)

type Bulletin struct {
	Author string `json:"author" binding:"required"`
	Content string `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

func GetBulletins() ([]Bulletin, error) {
	const q = `SELECT author, content, created_at FROM bulletins ORDER BY created_at DESC LIMIT 100`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	results := make([]Bulletin, 0)

	for rows.Next() {
		var author string
		var content string
		var createAt time.Time
		err = rows.Scan(&author, &content, &createAt)
		if err != nil {
			return nil, err
		}
		results = append(results, Bulletin{author, content, createAt})
	}

	return results, nil
}

func AddBulletin(bulletin Bulletin) error {
	const q = `INSERT INTO bulletins(author, content, created_at) VALUES ($1, $2, $3)`
	_, err := db.Exec(q, bulletin.Author, bulletin.Content, bulletin.CreatedAt)
	return err
}

func main() {
	var err error

	r := gin.Default()
	r.GET("/board", func(context *gin.Context) {
		results, err := GetBulletins()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"status": "internal error: " + err.Error()})
			return
		}
		context.JSON(http.StatusOK, results)
	})

	r.POST("/board", func(context *gin.Context) {
		var b Bulletin

		if context.Bind(&b) == nil {
			b.CreatedAt = time.Now()
			if err := AddBulletin(b); err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"status": "internal error: " + err.Error()})
				return
			}
			context.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", DbHost, DbUser, DbPassword, DbName)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.Query(Migration)
	if err != nil {
		log.Println("failed to run migrations", err.Error())
		return
	}

	log.Println("running..")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
