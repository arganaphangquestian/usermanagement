package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/arganaphangquestian/usermanagement/server/model"
	"github.com/jackc/pgx/v4"
)

type (
	postgres struct {
		db *pgx.Conn
	}
)

// Register methods
func (r *postgres) Register(register model.InputUser) (*model.User, error) {
	var role string
	if register.Referral == nil || *register.Referral == "" {
		role = "ADMIN"
	} else {
		role = "USER"
	}
	var id uint64 = 0
	err := r.db.QueryRow(context.Background(), `
	INSERT INTO users (name, username, email, password, role, referral)
	VALUES($1, $2, $3, $4, $5, $6) RETURNING id`, register.Name, register.Username, register.Email, register.Password, role, register.Referral).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &model.User{
		ID:       fmt.Sprintf("%d", id),
		Username: register.Username,
		Name:     register.Name,
		Email:    register.Email,
		Referral: register.Referral,
		Role:     role,
	}, nil
}

// Users methods
func (r *postgres) Users() ([]*model.User, error) {
	rows, err := r.db.Query(context.Background(), `SELECT id, name, username, email, referral, role FROM users`)
	if err != nil {
		log.Println(err)
	}
	var users []*model.User
	var id uint64 = 0
	for rows.Next() {
		var u model.User
		err = rows.Scan(&id, &u.Name, &u.Username, &u.Email, &u.Referral, &u.Role)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		u.ID = fmt.Sprintf("%d", id)
		users = append(users, &u)
	}
	return users, nil
}

// Users methods
func (r *postgres) GetUserByUsername(username string) (*model.UserWithPassword, error) {
	var u model.UserWithPassword
	var id uint64 = 0
	err := r.db.QueryRow(context.Background(), `SELECT id, name, username, password, email, referral, role FROM users where username=$1`, username).Scan(&id, &u.Name, &u.Username, &u.Password, &u.Email, &u.Referral, &u.Role)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	u.ID = fmt.Sprintf("%d", id)
	return &u, nil
}

// New UserRepository
func New() UserRepository {
	conn, err := pgx.Connect(context.Background(), "postgres://argadev:123456@127.0.0.1:5432/userdatabase?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	// defer conn.Close(context.Background())
	return &postgres{conn}
}
