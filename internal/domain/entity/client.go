package entity

type Client struct {
	ID        string
	Name      string
	CreatedAt string
}

type ClientCredential struct {
	ID        string
	ClientID  string
	Key       string
	Name      string
	IsActive  bool
	CreatedAt string
}
