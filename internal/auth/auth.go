package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hydrocode-de/gorun/internal/db"
)

func ValidateApiKey(key string, ctx context.Context, DB *db.Queries) error {
	hash := sha256.Sum256([]byte(key))
	hashedKey := hex.EncodeToString(hash[:])

	status, err := DB.ValidateApiKey(ctx, hashedKey)
	if err != nil {
		return err
	}

	if status == "valid" {
		err = DB.UpdateApiKeyLastUsed(ctx, hashedKey)
		if err != nil {
			return err
		}
	}

	switch status {
	case "valid":
		return nil
	case "expired":
		return fmt.Errorf("the API key has expired")
	case "invalid":
		return fmt.Errorf("the API key is invalid")
	}

	return fmt.Errorf("unknown API key status: %s", status)
}

func CreateNewApiKey(ctx context.Context, DB *db.Queries, validFor time.Duration) (string, error) {
	// Generate a random key with only letters
	const letters = "abcdefghiyklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	key := make([]byte, 32)
	for i := range key {
		// Generate a random index into the letters string
		idx := make([]byte, 1)
		rand.Read(idx)
		key[i] = letters[int(idx[0])%len(letters)]
	}

	hash := sha256.Sum256(key)
	hashedKey := hex.EncodeToString(hash[:])

	var validUntil sql.NullTime
	if validFor == 0 {
		validUntil = sql.NullTime{
			Valid: false,
			Time:  time.Time{},
		}
	} else {
		validUntil = sql.NullTime{
			Valid: true,
			Time:  time.Now().Add(validFor),
		}
	}

	_, err := DB.CreateApiKey(ctx, db.CreateApiKeyParams{
		Key:        hashedKey,
		ValidUntil: validUntil,
	})
	if err != nil {
		return "", err
	}

	return string(key), nil
}
