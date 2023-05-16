package googleApi

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"os"
	"path/filepath"

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

	currentDir, err := os.Getwd()
	if err != nil {
		log.WithError(err).Errorf("failed to get current directory path")
		panic(errors.Wrap(err, "failed to get current directory path"))
	}

	credFile := filepath.Join(currentDir, "credentials.json")

	privateCredBytes, err := os.ReadFile(credFile)
	if err != nil {
		log.WithError(err).Errorf("unable to read client secret file")
		panic(errors.Wrap(err, "unable to read client secret file"))
	}

	scopes := []string{admin.AdminDirectoryUserScope}

	myConfig, err := google.JWTConfigFromJSON(privateCredBytes, scopes...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	myConfig.Subject = "mykhailo.velykodnyi@centrilisedgym.online"

	// Use the client to authenticate API requests
	client := myConfig.Client(ctx)

	service, err := admin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
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
