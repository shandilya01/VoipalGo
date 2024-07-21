package services

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/shandilya01/VoipalGo/internal/models"
	"github.com/shandilya01/VoipalGo/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *pgx.Conn
}

func NewUserService(db *pgx.Conn) *UserService {
	return &UserService{db: db}
}

func ConvertInterfaceToString(reqBody map[string]interface{}, fields []string) (map[string]string, error) {
	stringUserObj := make(map[string]string)

	for _, field := range fields {
		value, ok := reqBody[field]
		if !ok {
			stringUserObj[field] = ""
			continue
		}
		valueString, ok := value.(string)
		if !ok {
			return nil, errors.New(field + " cannot be stringified")
		}
		stringUserObj[field] = valueString
	}

	return stringUserObj, nil
}

func (s *UserService) UserLogin(ctx context.Context, reqBody map[string]interface{}) (*models.User, error) {
	repo := repository.NewUserRepository(s.db)

	stringUserObj, err := ConvertInterfaceToString(reqBody, []string{"email", "password", "pushToken"})
	if err != nil {
		return nil, err
	}

	userObj := repo.FindUserByEmail(ctx, stringUserObj["email"])

	if userObj == nil {
		log.Print("No user found in db")
		return nil, errors.New("no user found in database for email")
	}

	if err := bcrypt.CompareHashAndPassword(userObj.Password, []byte(stringUserObj["password"])); err != nil {
		// TODO: Properly handle error
		log.Print("Password do not match")
		return nil, errors.New("password do not match")
	}

	// login success till now so update the pushToken

	tokenError := repo.UpdatePushToken(ctx, stringUserObj["email"], stringUserObj["pushToken"])

	return userObj, tokenError
}

func (s *UserService) UserSignUp(ctx context.Context, reqBody map[string]interface{}) (*models.User, error) {
	repo := repository.NewUserRepository(s.db)

	stringUserObj, err := ConvertInterfaceToString(reqBody, []string{"email", "password", "incantation", "phoneNumber", "pushToken", "name"})
	if err != nil {
		return nil, err
	}

	userObj := repo.FindUserByEmail(ctx, stringUserObj["email"])

	if userObj != nil {
		return nil, errors.New("user email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(stringUserObj["password"]), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("unable to hash the password")
	}

	err = repo.CreateUser(ctx, hashedPassword, stringUserObj)
	if err != nil {
		return nil, err
	}

	return userObj, nil
}
