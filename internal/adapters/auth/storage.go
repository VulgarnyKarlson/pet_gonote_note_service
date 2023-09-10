package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

type storageModelDB struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	IsValid  bool   `json:"is_valid"`
}

func (c *ClientImpl) store(ctx context.Context, token string, res *ValidateTokenResponse) error {
	var storageDto storageModelDB
	if res.User != nil {
		storageDto.UserID = res.User.ID()
		storageDto.UserName = res.User.UserName()
		storageDto.IsValid = true
	} else {
		storageDto.IsValid = false
	}

	val, err := json.Marshal(storageDto)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %v", err)
	}

	err = c.storage.Set(ctx, token, val, c.config.BackupStorageTime)
	if err != nil {
		return fmt.Errorf("failed to store token: %v", err)
	}

	return nil
}

func (c *ClientImpl) get(ctx context.Context, token string) (*ValidateTokenResponse, error) {
	var storageDto storageModelDB
	val, err := c.storage.Get(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}
	// string to struct
	err = json.Unmarshal([]byte(val), &storageDto)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %v", err)
	}

	var res *ValidateTokenResponse
	if storageDto.IsValid {
		res = &ValidateTokenResponse{
			User:  domain.NewUser(storageDto.UserID, storageDto.UserName),
			Valid: true,
		}
	} else {
		res = &ValidateTokenResponse{
			Valid: false,
		}
	}
	return res, nil
}
