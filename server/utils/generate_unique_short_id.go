package utils

import (
	"context"
	"fmt"
	"picourl-backend/db"
)

func GenerateUniqueShortId(ctx context.Context) (string, error) {
	return generateUniqueShortId(ctx, 0)
}

func generateUniqueShortId(ctx context.Context, attempts int) (string, error) {
	const maxRetries = 10

	if attempts > maxRetries {
		return "", fmt.Errorf("failed to generate unique ID after %d attempts", maxRetries)
	}

	id, err := GenerateRandomBase62String(6)
	if err != nil {
		return "", fmt.Errorf("failed to generate unique short ID: %w", err)
	}

	exists, err := db.Queries.LinkWithIdExists(ctx, id)
	if err != nil {
		return "", err
	}

	if exists {
		return generateUniqueShortId(ctx, attempts+1)
	}

	return id, nil
}
