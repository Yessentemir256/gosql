package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// адрес подключения
	// протокол://логин:пароль@хост:порт/бд
	dsn := "postgres://app:pass@localhost:5434/db"
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
	ctx := context.Background()
	_, err = db.ExecContext(ctx, `
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

	name := "Victor"
	phone := "+992000000002"

	result, err := db.ExecContext(ctx, `
	  INSERT INTO customers(name, phone) VALUES ($1, $2) ON CONFLICT (phone) DO UPDATE SET name = excluded.name;
	  `, name, phone)
	if err != nil {
		log.Print(err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("Количество строк изменено: %d", rowsAffected)

	customer := &Customer{}
	err = db.QueryRowContext(ctx, `
	  Select id, name, phone, active, created FROM customers WHERE id = 1
	  `).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("%#v", customer)

	id := 1
	newPhone := "+992000000099"
	err = db.QueryRowContext(ctx, `
	  UPDATE customers SET phone = $2 WHERE id = $1 RETURNING id, name, phone, active, created
	  `, id, newPhone).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	if errors.Is(err, sql.ErrNoRows) {
		log.Print("No rows")
		return
	}

	if err != nil {
		log.Print(err)
		return
	}
}

type Customer struct {
	ID      int64
	Name    string
	Phone   string
	Active  bool
	Created time.Time
}
