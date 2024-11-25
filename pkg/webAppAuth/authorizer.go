package webAppAuth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TelegramWebAppUser struct {
	Id              int64  `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Username        string `json:"username"`
	LanguageCode    string `json:"language_code"`
	IsPremium       bool   `json:"is_premium"`
	AllowsWriteToPm bool   `json:"allows_write_to_pm"`
	PhotoUrl        string `json:"photo_url"`
}

type Authorizer struct {
	botToken string
}

func NewAuthorizer(botToken string) *Authorizer {
	return &Authorizer{botToken: botToken}
}

// AuthorizeUser validates the data received from Telegram WebApp and return authorized user data
func (a *Authorizer) AuthorizeUser(initData string) (TelegramWebAppUser, error) {
	if initData == "" {
		return TelegramWebAppUser{}, fmt.Errorf("initData is empty")
	}

	// Parse the initData as a query string
	parsedData, err := url.ParseQuery(initData)
	if err != nil {
		return TelegramWebAppUser{}, fmt.Errorf("error parsing initData: %w", err)
	}

	// Extract the received hash
	receivedHash := parsedData.Get("hash")
	if receivedHash == "" {
		return TelegramWebAppUser{}, fmt.Errorf("hash parameter is missing")
	}
	delete(parsedData, "hash") // Remove the hash parameter for data-check-string

	// Build the data-check-string
	var dataCheckParts []string
	for key, values := range parsedData {
		// Telegram guarantees single values per key
		dataCheckParts = append(dataCheckParts, fmt.Sprintf("%s=%s", key, values[0]))
	}

	// Sort alphabetically
	sort.Strings(dataCheckParts)
	dataCheckString := strings.Join(dataCheckParts, "\n")

	// Compute the secret key
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(a.botToken))
	secretKey := h.Sum(nil)

	// Compute the HMAC-SHA-256 signature of the data-check-string
	h = hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	computedHash := hex.EncodeToString(h.Sum(nil))

	// Compare the computed hash with the received hash
	if computedHash != receivedHash {
		return TelegramWebAppUser{}, fmt.Errorf("invalid hash")
	}

	// Optionally validate the auth_date to ensure it's not outdated
	authDate := parsedData.Get("auth_date")
	if authDate == "" {
		return TelegramWebAppUser{}, fmt.Errorf("auth_date parameter is missing")
	}

	authDateInt, err := strconv.ParseInt(authDate, 10, 64)
	if err != nil {
		return TelegramWebAppUser{}, fmt.Errorf("invalid auth_date: %w", err)
	}

	// Check if the auth_date is too old (e.g., older than 24 hours)
	if time.Now().Unix()-authDateInt > 24*60*60 {
		return TelegramWebAppUser{}, fmt.Errorf("auth_date is too old")
	}

	// Parsing user data
	var user TelegramWebAppUser
	if err := json.Unmarshal([]byte(parsedData.Get("user")), &user); err != nil {
		return TelegramWebAppUser{}, fmt.Errorf("error unmarshalling user data: %w", err)
	}

	return user, nil
}
