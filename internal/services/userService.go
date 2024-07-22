package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shandilya01/VoipalGo/internal/models"
	"github.com/shandilya01/VoipalGo/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const expoPushURL = "https://exp.host/--/api/v2/push/send" //"https://api.expo.dev/v2/push/send"

type ExpoPushMessage struct {
	To        string `json:"to"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	ChannelID string `json:"channelId,omitempty"`
}

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
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

func (s *UserService) GetUserContactsById(ctx context.Context, id string) ([]*models.Contact, error) {
	repo := repository.NewUserRepository(s.db)

	var contactsArr []*models.Contact
	if id != "" {
		contactsArr = repo.GetContactsById(ctx, id)
	}

	return contactsArr, nil
}

func (s *UserService) CallPushNotification(ctx context.Context, userId string, peerId string) error {
	repo := repository.NewUserRepository(s.db)

	peerToken := repo.GetPushToken(ctx, peerId)

	if len(*peerToken) == 0 {
		return errors.New("could not get push token")
	}

	title := "Incoming Call"
	body := "Incoming Call From User Id : " + userId
	// data := "incoming"
	return SendPushNotfication(*peerToken, title, body, "voipalCall")
	// return SendPushNotfication(peerToken, &title, &body, &data)
}

func SendPushNotfication(token, title, body, channelId string) error {
	message := []ExpoPushMessage{
		{
			To:        token,
			Title:     title,
			Body:      body,
			ChannelID: channelId,
		},
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return errors.New("failed to marshal message")
	}

	req, err := http.NewRequest("POST", expoPushURL, bytes.NewBuffer(payload))
	if err != nil {
		return errors.New("failed to create request")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Fetch-Mode", "cors")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("failed to send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("status code :" + strconv.Itoa(resp.StatusCode))
	}

	return nil
}
