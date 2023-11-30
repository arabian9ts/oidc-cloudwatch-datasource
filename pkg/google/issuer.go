package google

import (
	"context"
	"fmt"

	"google.golang.org/api/idtoken"
)

type TokenIssuer struct {
	credFilePath string
	audience     string
}

type Config struct {
	CredsFilePath string
	Audience      string
}

func NewTokenIssuer(cfg *Config) *TokenIssuer {
	return &TokenIssuer{
		credFilePath: cfg.CredsFilePath,
		audience:     cfg.Audience,
	}
}

func (o *TokenIssuer) IssueAccessToken(ctx context.Context) (string, error) {
	opts := make([]idtoken.ClientOption, 0, 1)
	if o.credFilePath != "" {
		opts = append(opts, idtoken.WithCredentialsFile(o.credFilePath))
	}
	ts, err := idtoken.NewTokenSource(ctx, o.audience, opts...)
	if err != nil {
		return "", fmt.Errorf("failed to create token source: %w", err)
	}

	t, err := ts.Token()
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	return t.AccessToken, nil
}
