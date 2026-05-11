package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleProvider struct {
	config *oauth2.Config
}

func NewGoogleProvider(clientID, clientSecret, redirectURL string) Provider {
	return &googleProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (g *googleProvider) Name() string { return "google" }

func (g *googleProvider) AuthCodeURL(state string) string {

	return g.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (g *googleProvider) Exchange(code string) (*UserInfo, error) {
	token, err := g.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("google: code exchange failed: %w", err)
	}

	client := g.config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("google: userinfo request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google: userinfo returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("google: reading userinfo body: %w", err)
	}

	var raw struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("google: parsing userinfo: %w", err)
	}

	return &UserInfo{
		ProviderID: raw.Sub,
		Email:      raw.Email,
		Name:       raw.Name,
	}, nil
}
