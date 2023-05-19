package ally

import (
	"context"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/oauth1"
	"golang.org/x/time/rate"
)

var RateLimitMap = map[rateLimitType]*rate.Limiter{
	tradesRL:   rate.NewLimiter(rate.Every(60*time.Second), 40),
	quotesRL:   rate.NewLimiter(rate.Every(60*time.Second), 60),
	authedRL:   rate.NewLimiter(rate.Every(60*time.Second), 180),
	infiniteRL: rate.NewLimiter(rate.Every(60*time.Second), math.MaxInt),
}

func NewClient(rl *rate.Limiter) *RLClient {
	creds := credentials{
		client_key:            os.Getenv("ALLY_CLIENT_KEY"),
		client_secret:         os.Getenv("ALLY_CLIENT_SECRET"),
		resource_owner_key:    os.Getenv("ALLY_RESOURCE_OWNER_KEY"),
		resource_owner_secret: os.Getenv("ALLY_RESOURCE_OWNER_SECRET"),
	}

	config := oauth1.NewConfig(creds.client_key, creds.client_secret)
	token := oauth1.NewToken(creds.resource_owner_key, creds.resource_owner_secret)

	c := &RLClient{
		client:      config.Client(oauth1.NoContext, token),
		Ratelimiter: rl,
	}
	return c
}

func (c *RLClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := c.Ratelimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type credentials struct {
	client_key            string
	client_secret         string
	resource_owner_key    string
	resource_owner_secret string
}

type RLClient struct {
	client      *http.Client
	Ratelimiter *rate.Limiter
}

type rateLimitType string

const (
	tradesRL   rateLimitType = "trades"
	quotesRL   rateLimitType = "quotes"
	authedRL   rateLimitType = "authed"
	infiniteRL rateLimitType = "infinite"
)
