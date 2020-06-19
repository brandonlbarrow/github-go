package github

import (
	"context"
	"fmt"
	g "github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"bufio"
	"net/http"
	"os"
	"strings"
)

type Client struct {
	*g.Client
}

type OAuth2Provider struct {
	token oauth2.TokenSource
	creds Credentials
}

type HttpBasicProvider struct {
	basic g.BasicAuthTransport
	creds Credentials
}

type Provider interface {
	Auth() *Client
}

type Credentials struct {
	Username    string
	Password    string
	token       string
	Multifactor bool
}

func NewOAuth2Provider(c *Credentials) *OAuth2Provider {
	return &OAuth2Provider{
		token: oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: c.token},
		),
	}
}

func (o *OAuth2Provider) Auth() *Client {
	oauthClient := oauth2.NewClient(context.Background(), o.token)
	return &Client{g.NewClient(oauthClient)}
}

func NewHttpBasicProvider(c *Credentials) (*HttpBasicProvider, error) {
	switch mf := c.Multifactor; mf {
	case true:
		in := bufio.NewReader(os.Stdin)
		fmt.Println("enter OTP token")
		otp, err := in.ReadString('\n')
		if err != nil {
			return nil, err
		}
		fmt.Println("OTP token submitted")
		return &HttpBasicProvider{
			basic: g.BasicAuthTransport{
				Username: c.Username,
				Password: c.Password,
				OTP: strings.TrimSpace(otp),
			},
		}, nil
	}
	return &HttpBasicProvider{
		basic: g.BasicAuthTransport{
			Username: c.Username,
			Password: c.Password,
		},
	}, nil
}

func (b *HttpBasicProvider) Auth() *Client {
	basicClient := http.Client{Transport: &b.basic}
	return &Client{g.NewClient(&basicClient)}
}

func (c *Client) ListRepos() ([]*g.Repository, *g.Response, error) {
	return c.Repositories.ListAll(context.Background(), &g.RepositoryListAllOptions{})
}
