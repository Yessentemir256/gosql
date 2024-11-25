package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// адрес подключения
	// протокол://логин:пароль@хост:порт/бд
	dsn := "postgres://app:pass@localhost:5432/db"
	// получение указателя на структуру для работы с БД
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return
	}
	// закрытие структуры
	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()
	// TODO: запросы
}
