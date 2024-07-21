package repository

import (
	"context"
	"errors"
	"log"

	"github.com/shandilya01/VoipalGo/internal/models"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// 	id          int    `json:"id"`
// 	name        string `json:"name"`
// 	email       string `json:"email"`
// 	phoneNumber string `json:"phoneNumber"`
// 	spell       string `json:"spell"`
// 	active      bool   `json:"active"`
// 	pushToken   string `json:"pushToken"`

func (r *UserRepository) CreateUser(ctx context.Context, password []byte, userObj map[string]string) error {
	query := `insert into users (name, email, password, phoneNumber, incantation, pushToken) values ($1, $2, $3, $4, $5, $6)`
	log.Print("creating user")
	_, err := r.db.Exec(ctx, query, userObj["name"], userObj["email"], password, userObj["phoneNumber"], userObj["incantation"], userObj["pushToken"])
	log.Print("user creation error", err)
	return err
}

func (r *UserRepository) FindUserById(ctx context.Context, userId int) *models.User {
	query := `select * from users where id = $1`
	var userObj *models.User
	err := r.db.QueryRow(ctx, query, userId).Scan(userObj)
	if err != nil {
		log.Print("find by user id error", err)
		return nil
	}
	return userObj
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) *models.User {
	query := `select * from users where email = $1`
	rows, _ := r.db.Query(ctx, query, email)
	log.Print("rows", rows)
	userObj, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.User])
	log.Print("rows", userObj)
	if err != nil {
		log.Print("find by emailerrror", err)
		return nil
	}
	return userObj
}

func (r *UserRepository) UpdatePushToken(ctx context.Context, email string, token string) error {
	query := `update users set pushToken=$1 where email = $2`
	_, err := r.db.Query(ctx, query, token, email)
	if err != nil {
		return errors.New("could not update token, user might not recieve notifications/calls")
	}
	return nil
}
