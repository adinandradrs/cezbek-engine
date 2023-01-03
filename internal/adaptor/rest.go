package adaptor

import (
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"time"
)

type Rest struct {
	BackoffInterval   time.Duration
	MaxJitterInterval time.Duration
	Timeout           time.Duration
}

func (r *Rest) Client() *httpclient.Client {
	retrier := heimdall.NewRetrier(heimdall.NewConstantBackoff(r.BackoffInterval, r.MaxJitterInterval))
	return httpclient.NewClient(httpclient.WithHTTPTimeout(r.Timeout), httpclient.WithRetrier(retrier), httpclient.WithRetryCount(2))
}
