package backend

import (
	"context"
	"net/http"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func setupSecretAPI(swPlugin *swagger.Plugin) {
	api := &SecretAPI{}
	tag := swPlugin.AddTag(&swagger.Tag{Name: "Secret", Description: "Secret API list"})
	var hInfo *swagger.HandlerInfo

	hInfo = swagger.NewHandlerInfo(api.Post)
	ucon.Handle(http.MethodPost, "/api/1/secret", hInfo)
	hInfo.Description, hInfo.Tags = "post to secret", []string{tag.Name}

	hInfo = swagger.NewHandlerInfo(api.Get)
	ucon.Handle(http.MethodGet, "/api/1/secret/{key}", hInfo)
	hInfo.Description, hInfo.Tags = "get from secret", []string{tag.Name}
}

// Secret is Datastore Entity
type Secret struct {
	Value string `datastore:",noindex"`
}

// SecretAPI is API to register and acquire Secret
type SecretAPI struct{}

// SecretAPIPostRequest is SecretAPI Post Request
type SecretAPIPostRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Post is Secret registration handler
func (api *SecretAPI) Post(ctx context.Context, form *SecretAPIPostRequest) error {
	u := user.Current(ctx)
	if u == nil {
		return &HTTPError{Code: http.StatusForbidden, Message: "You do not have permission."}
	}
	log.Infof(ctx, "%+v", u)

	ds, err := FromContext(ctx)
	if err != nil {
		return err
	}
	k := ds.NameKey("Secret", form.Key, nil)
	s := &Secret{
		Value: form.Value,
	}
	_, err = ds.Put(ctx, k, s)
	if err != nil {
		return err
	}

	return nil
}

// SecretAPIGetRequest is SecretAPI Get Request
type SecretAPIGetRequest struct {
	Key string `json:"key" swagger:",in=query"`
}

// SecretAPIGetResponse is SecretAPI Get Response
type SecretAPIGetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Get is Secret registration handler
func (api *SecretAPI) Get(ctx context.Context, form *SecretAPIGetRequest) (*SecretAPIGetResponse, error) {
	u := user.Current(ctx)
	if u == nil {
		return nil, &HTTPError{Code: http.StatusForbidden, Message: "You do not have permission."}
	}
	log.Infof(ctx, "%+v", u)

	ds, err := FromContext(ctx)
	if err != nil {
		return nil, err
	}
	k := ds.NameKey("Secret", form.Key, nil)
	s := &Secret{}
	if err := ds.Get(ctx, k, s); err != nil {
		return nil, err
	}

	return &SecretAPIGetResponse{
		Key:   form.Key,
		Value: s.Value,
	}, nil
}
