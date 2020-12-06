package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

)

var ErrNotFound = errors.New("item not found")

var ErrInternal = errors.New("internal error")

type Service struct {
	pool  *pgxpool.Pool
}
func NewService(pool  *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE id = $1
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil

}
func (s *Service) All(ctx context.Context) (items []*Customer, err error) {

	rows, err:= s.pool.Query(ctx, `
		SELECT * FROM customers
	`)

	for rows.Next(){
		item := &Customer{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)
		if err != nil {
			log.Print(err)
		}

		items = append(items, item)
	}
	return items, nil
}
func (s *Service) AllActive(ctx context.Context) (c []*Customer, err  error) {

	items := c
	rows, err := s.pool.Query(ctx, `
			SELECT id, name, phone, active, created FROM customers WHERE active
		`)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, ErrNotFound
		}
		items = append(items, item)
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return items, nil
}

// //Save method
// // func (s *Service) Save(ctx context.Context, item *Customer) (c *Customer,err error) {
// // 	items := &Customer{}
// // 	if item.ID == 0 {
// // 		err = s.db.QueryRowContext(ctx, `
// // 		INSERT INTO customers(name,phone) VALUES($1,$2) RETURNING *
// // 		`,item.Name, item.Phone).Scan(
// // 			&items.ID,
// // 			&items.Name,
// // 			&items.Phone,
// // 			&items.Active,
// // 			&items.Created)

// // 	} else {
// // 		err = s.db.QueryRowContext(ctx, `
// // 		UPDATE customers SET name=$1, phone=$2 WHERE id=$3 RETURNING *
// // 		`, item.Name, item.Phone, item.ID).Scan(
// // 			&items.ID,
// // 			&items.Name,
// // 			&items.Phone,
// // 			&items.Active,
// // 			&items.Created)
// // 	}

// // 	if err != nil {
// // 		log.Print(err)
// // 		return nil, ErrInternal
// // 	}

// // 	return items,nil

// // }

// //Save method
func (s *Service) Save(ctx context.Context, customer *Customer) (c *Customer, err error) {

	item := &Customer{}

	if customer.ID == 0 {
		sqlStatement := `insert into customers(name, phone) values($1, $2) returning *`
		err = s.pool.QueryRow(ctx, sqlStatement, customer.Name, customer.Phone).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)
	} else {
		sqlStatement := `update customers set name=$1, phone=$2 where id=$3 returning *`
		err = s.pool.QueryRow(ctx, sqlStatement, customer.Name, customer.Phone, customer.ID).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

func (s *Service) RemoveById(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.pool.QueryRow(ctx, `
	DELETE FROM customers WHERE id=$1 RETURNING id,name,phone,active,created 
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil

}

func (s *Service) BlockByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.pool.QueryRow(ctx, `
		UPDATE customers SET active = false WHERE id = $1 RETURNING id, name, phone, active, created
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil

}
func (s *Service) UnBlockByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.pool.QueryRow(ctx, `
		UPDATE customers SET active = true WHERE id = $1 RETURNING id, name, phone, active, created
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil

}
