package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"time"
)

func main() {
	// адрес подключения
	// протокол://логин:пароль@хост:порт/бд
	dsn := "postgres://app:pass@localhost:5434/db"
	// получение указателя на структуру для работы с БД
	connectCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
	pool, err := pgxpool.Connect(connectCtx, dsn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return
	}
	// закрытие структуры
	defer pool.Close()

	log.Print(pool.Stat().TotalConns()) // 1 подключение
	// контекст для запросов
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Release()
	log.Print(pool.Stat().TotalConns()) // 1 подключение

	conn, err = pool.Acquire(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	conn.Release()
	log.Print(pool.Stat().TotalConns()) // 2, несмотря на то, что закрыли

	// новое подключение не создастся, будет переисользовано предыдущее
	conn, err = pool.Acquire(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	conn.Release()
	log.Print(pool.Stat().TotalConns()) // 2, несмотря на то, что закрыли

	// новое подключение не создастся, будет переиспользовано предыдущее
	conn, err = pool.Acquire(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	conn.Release()
	log.Print(pool.Stat().TotalConns) // 2, несмотря на то, что закрыли

	_, err = pool.Exec(ctx, `
	  CREATE TABLE IF NOT EXISTS customers (
		id      BIGSERIAL PRIMARY KEY,
		name    TEXT  NOT NULL,
		phone   TEXT NOT NULL UNIQUE,
		active  BOOLEAN NOT NULL DEFAULT TRUE,
		created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		log.Print(err)
		return
	}
}
