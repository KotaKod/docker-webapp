package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getRequiredEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		log.Fatalf("Ошибка конфигурации: переменная окружения %s обязательна и не задана", key)
	}
	return value
}

func main() {
	port := getEnv("PORT", "8080")
	staticDir := getEnv("STATIC_DIR", "/app/static")

	// Получаем параметры подключения строго из окружения (.env через docker-compose)
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getRequiredEnv("DB_USER")
	dbPassword := getRequiredEnv("DB_PASSWORD")
	dbName := getRequiredEnv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Внимание: БД пока недоступна при старте: %v", err)
	}

	// 1. API: список записей
	http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM items ORDER BY id ASC")
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка запроса к базе данных: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		items := make([]Item, 0)
		for rows.Next() {
			var item Item
			if err := rows.Scan(&item.ID, &item.Name); err != nil {
				http.Error(w, fmt.Sprintf("Ошибка чтения строки: %v", err), http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})

	// 2. Статика без циклического редиректа
	fileServer := http.FileServer(http.Dir(staticDir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// При прямом запросе корня отдаем index.html без внутренних редиректов
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}

		// Для запросов вида /style.css, /app.js
		fileServer.ServeHTTP(w, r)
	})

	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}