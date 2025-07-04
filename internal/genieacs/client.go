package genieacs

import (
	"github.com/go-resty/resty/v2"
	"github.com/praction-networks/acs-proxy/internal/config"
)

type Client struct {
	http *resty.Client
}

func NewClient(cfg *config.GeniacsConfig) *Client {
	return &Client{
		http: resty.New().SetBaseURL(cfg.NBI_URL),
	}
}
