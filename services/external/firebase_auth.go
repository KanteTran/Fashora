package external

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"

	"fashora-backend/config"
	"fashora-backend/logger"
)

type firebaseClient struct {
	app *auth.Client
}

type IFirebaseClient interface {
	VerifyIdToken(ctx context.Context, idToken string) (uid string, err error)
}

var _ IFirebaseClient = &firebaseClient{}

// NewFirebaseClient init new firebase client from instance
func NewFirebaseClient(app *auth.Client) IFirebaseClient {
	return &firebaseClient{app: app}
}

// NewFirebaseClientFromConfig init new firebase client from config
func NewFirebaseClientFromConfig(ctx context.Context, cfg config.Config) (IFirebaseClient, error) {
	// Firebase initialization
	opt := option.WithCredentialsFile(cfg.FireBase.FileKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Errorf("Failed to initialize Firebase: %s", err)
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		logger.Errorf("Failed to initialize Firebase Auth: %s", err)
		return nil, err
	}

	return &firebaseClient{app: authClient}, nil
}

// VerifyIdToken return uid of firebase user
func (f *firebaseClient) VerifyIdToken(ctx context.Context, idToken string) (string, error) {
	token, err := f.app.VerifyIDToken(ctx, idToken)
	if err != nil {
		logger.Errorf("Failed to verify ID Token: %s", err)
		return "", nil
	}

	return token.UID, nil
}
