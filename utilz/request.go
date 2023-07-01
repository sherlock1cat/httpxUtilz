package utilz

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// RequestClientConfig Including configuration options for the requesting client.
type RequestClientConfig struct {
	ProxyURL        string
	UseHTTPS        bool
	FollowRedirects bool
	MaxRedirects    int
	Method          string
	RandomUserAgent bool
	Headers         map[string]string
	FollowSameHost  bool
	Timeout         time.Duration
}

// NewRequestClient Create a new request client
func NewRequestClient(config RequestClientConfig) *http.Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	} else {
		config.Timeout = config.Timeout * time.Second
	}
	if config.FollowRedirects && config.MaxRedirects == 0 {
		config.MaxRedirects = 10
	}
	if config.Method == "" {
		config.Method = "GET"
	}
	if config.RandomUserAgent && len(config.Headers["User-Agent"]) == 0 {
		config.Headers["User-Agent"] = getRandomUserAgent()
	}
	if config.FollowSameHost && config.MaxRedirects == 0 {
		config.FollowSameHost = false
	}

	transport := &http.Transport{
		Proxy:           getProxy(config.ProxyURL),
		TLSClientConfig: getTLSConfig(config.UseHTTPS),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !config.FollowRedirects {
				return http.ErrUseLastResponse
			}
			if !config.FollowSameHost && len(via) > 0 {
				if req.URL.Host != via[len(via)-1].URL.Host {
					return http.ErrUseLastResponse
				}
			}
			if len(via) >= config.MaxRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	return client
}

// getProxy
func getProxy(proxyURL string) func(*http.Request) (*url.URL, error) {
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			panic(err)
		}
		return http.ProxyURL(proxy)
	}
	return nil
}

// getTLSConfig Return TLS configuration based on the flag indicating the use of HTTPS.
func getTLSConfig(useHTTPS bool) *tls.Config {
	if useHTTPS {
		return &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	return nil
}

// getRandomUserAgent Return a randomly generated User-Agent string
func getRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:86.0) Gecko/20100101 Firefox/86.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/86.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/89.0.774.54 Safari/537.36",
	}
	rand.Seed(time.Now().UnixNano())
	return userAgents[rand.Intn(len(userAgents))]
}

func parseUrl(targetUrl string) (string, error) {
	Url, err := url.Parse(targetUrl)
	if err != nil {
		return "", err
	}
	return Url.String(), nil
}
