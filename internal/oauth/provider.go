package oauth

type UserInfo struct {
	ProviderID string
	Email      string
	Name       string
}

type Provider interface {
	Name() string
	AuthCodeURL(state string) string
	Exchange(code string) (*UserInfo, error)
}
