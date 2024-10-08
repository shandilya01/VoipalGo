package repository

import (
	"context"
	"errors"
	"log"

	"github.com/shandilya01/VoipalGo/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
The pgxpool.Pool manages a pool of database connections, which allows multiple goroutines to query the database concurrently.
The pool ensures that connections are reused efficiently and that there is a limit on the number of concurrent connections to the database.
*/

type UserRepository struct {
	db *pgxpool.Pool
}

type WordState struct {
	I int
	J int
	K int
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func updateWordState(state *WordState, count int) error {
	if state.K != count {
		state.K += 1
		return nil
	} else {
		state.K = 1
		if state.J != count {
			state.J += 1
			return nil
		} else {
			state.J = 1
			if state.I != count {
				state.I += 1
				return nil
			} else {
				return errors.New("words exhausted")
			}
		}
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, password []byte, userObj map[string]string) error {
	query := `select i,j,k from word_state`
	rows, _ := r.db.Query(ctx, query)
	wordState, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[WordState])
	if err != nil {
		return err
	}
	log.Print("i", wordState.I)

	query = `select count(*) from words`
	rows, _ = r.db.Query(ctx, query)
	wordCount := 0
	if rows.Next() {
		rows.Scan(&wordCount)
		log.Print("word count", wordCount)
	}

	query = `select word from words where id = $1`
	rows, _ = r.db.Query(ctx, query, wordState.I)
	word1 := ""
	if rows.Next() {
		rows.Scan(&word1)
	}
	rows, _ = r.db.Query(ctx, query, wordState.J)
	word2 := ""
	if rows.Next() {
		rows.Scan(&word2)
	}
	rows, _ = r.db.Query(ctx, query, wordState.K)
	word3 := ""
	if rows.Next() {
		rows.Scan(&word3)
	}

	err = updateWordState(wordState, wordCount)
	if err != nil {
		return err
	}

	query = `truncate word_state`
	_, err = r.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	query = `insert into word_state values (1,$1,$2,$3)`
	_, err = r.db.Exec(ctx, query, wordState.I, wordState.J, wordState.K)
	if err != nil {
		return err
	}

	userObj["voipId"] = word1 + "." + word2 + "." + word3

	query = `insert into users (name, email, password, phoneNumber, voipId, pushToken) values ($1, $2, $3, $4, $5, $6)`
	log.Print("creating user")
	_, err = r.db.Exec(ctx, query, userObj["name"], userObj["email"], password, userObj["phoneNumber"], userObj["voipId"], userObj["pushToken"])
	log.Print("user creation error", err)
	return err
}

func (r *UserRepository) FindUserById(ctx context.Context, userId string) *models.User {
	query := `select * from users where id = $1`
	rows, _ := r.db.Query(ctx, query, userId)
	// defer rows.Close()
	userObj, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.User])

	if err != nil {
		log.Print("find by userId err", err)
		return nil
	}
	return userObj
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) *models.User {
	query := `select * from users where email = $1`
	rows, _ := r.db.Query(ctx, query, email)
	// defer rows.Close()
	userObj, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.User])

	if err != nil {
		log.Print("find by emailerrror", err)
		return nil
	}
	return userObj
}

func (r *UserRepository) UpdatePushToken(ctx context.Context, email string, token string) error {
	query := `update users set pushToken=$1 where email = $2`
	_, err := r.db.Exec(ctx, query, token, email)
	if err != nil {
		return errors.New("could not update token, user might not recieve notifications/calls")
	}
	return nil
}

func (r *UserRepository) GetContactsById(ctx context.Context, id string) []*models.Contact {
	query := `select id, name, email, phonenumber, voipId from users where id != $1`
	rows, _ := r.db.Query(ctx, query, id)
	// defer rows.Close()
	contactRows, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Contact])
	if err != nil {
		log.Print("Could not get contacts for id ", id, err)
	}
	log.Print(contactRows)
	return contactRows
}

func (r *UserRepository) GetPushToken(ctx context.Context, id string) *string {
	query := `select pushtoken from users where id = $1`
	var token string
	err := r.db.QueryRow(ctx, query, id).Scan(&token)
	// defer rows.Close()
	// token, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[string])
	if err != nil {
		log.Print("Could not get token for id ", id, err)
	}
	log.Print(token)
	return &token
}

func (r *UserRepository) GetWordList(ctx context.Context) []string {
	query := `select word from words`
	rows, _ := r.db.Query(ctx, query)
	// defer rows.Close()

	wordList, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		log.Print("Could not get word list ", err)
	}
	return wordList
}

func (r *UserRepository) GetUserByVoipId(ctx context.Context, voipId string) *models.Contact {
	query := `select id,name,email,phoneNumber,voipId from users where voipId = $1`
	rows, _ := r.db.Query(ctx, query, voipId)
	// defer rows.Close()
	userObj, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Contact])

	if err != nil {
		log.Print("find by voipId err", err)
		return nil
	}
	return userObj
}
