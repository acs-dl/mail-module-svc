package googleApi

import (
	"context"
	"os"

	"gitlab.com/distributed_lab/logan/v3"

	"github.com/acs-dl/mail-module-svc/internal/config"
	"github.com/acs-dl/mail-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

type GoogleClient interface {
	GetUsersFromApi(domain string) ([]data.User, error)
	GetDomainUserFromApi(email string) (*data.User, error)
	SearchByUsersFromApi(email string) ([]data.User, error)

	AddUserInDomainFromApi(firstName string, lastName string, email string) (*data.User, error)
	DeleteUserInDomainFromApi(email string) error
}

type googleApi struct {
	client *admin.Service
	log    *logan.Entry
}

func NewGoogleAsInterface(cfg config.Config, ctx context.Context) interface{} {
	log := cfg.Log()

	credentials := os.Getenv("SERVICE_ACCOUNT_CREDENTIALS")
	if credentials == "" {
		log.Errorf("failed to get service account credentials")
		panic(errors.New("failed to get service account credentials"))
	}

	scopes := []string{admin.AdminDirectoryUserScope}

	myConfig, err := google.JWTConfigFromJSON([]byte(credentials), scopes...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		panic(errors.New(fmt.Sprintf("Unable to parse client secret file to config: %v", err)))
	}

	// Use the client to authenticate API requests
	client := myConfig.Client(ctx)

	service, err := admin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
		panic(errors.New(fmt.Sprintf("Failed to create service: %v", err)))
	}

	return interface{}(&googleApi{
		client: service,
		log:    log,
	})
}

func GoogleClientInstance(ctx context.Context) GoogleClient {
	return ctx.Value("google").(GoogleClient)
}

func CtxGoogleClientInstance(entry interface{}, ctx context.Context) context.Context {
	return context.WithValue(ctx, "google", entry)
}
