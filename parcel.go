package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	result, err := s.db.Exec(
		"INSERT INTO parcel (client, status, address, created_at) VALUES ($1, $2, $3, $4)",
		p.Client, p.Status, p.Address, p.CreatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении посылки: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID: %w", err)
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	p := Parcel{}
	err := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = $1", number).Scan(
		&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt,
	)

	if err != nil {
		return p, fmt.Errorf("ошибка получения посылки: %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	var res []Parcel
	rows, err := s.db.Query(
		"SELECT number, client, status, address, created_at FROM parcel WHERE client = $1",
		client,
	)
	if err != nil {
		return res, fmt.Errorf("ошибка запроса посылок клиента: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return res, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {

	_, err := s.db.Exec("UPDATE parcel SET status = $1 WHERE number = $2", status, number)
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса: %w", err)
	}
	return nil

}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec(
		"UPDATE parcel SET address = $1 WHERE number = $2 AND status = $3",
		address, number, ParcelStatusRegistered,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления адреса: %w", err)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec(
		"DELETE FROM parcel WHERE number = $1 AND status = $2",
		number, ParcelStatusRegistered,
	)
	if err != nil {
		return fmt.Errorf("ошибка удаления посылки: %w", err)
	}

	return nil
}
