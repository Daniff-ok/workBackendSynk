package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"rndmBackend/internal/domain"
)

type Logingateway struct {
	apiEndpoint  string
	clientId     string
	clientSecret string
	redirect     string
	client       http.Client
}

func NewLogingateway() *Logingateway {
	return &Logingateway{
		apiEndpoint:  "https://discord.com/api/v10",
		clientId:     "1574494144263360583",
		clientSecret: os.Getenv("DISCORD_OAUTH_SECRET"),
		redirect:     "http://127.0.0.1",
		client: http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *Logingateway) GetTokens(code string) (domain.TokenResponce, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", r.redirect)

	reqURL := fmt.Sprintf("%s/oauth2/token", r.apiEndpoint)
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(form.Encode()))
	if err != nil {
		return domain.TokenResponce{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(r.clientId, r.clientSecret)

	resp, err := r.client.Do(req)
	if err != nil {
		return domain.TokenResponce{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return domain.TokenResponce{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	tokens := domain.TokenResponce{}
	err = json.NewDecoder(resp.Body).Decode(&tokens)
	if err != nil {
		return domain.TokenResponce{}, fmt.Errorf("Json decode error: %w", err)
	}

	return tokens, nil
}

func (r *Logingateway) RefreshAccessToken(refreshToken string) (string, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)

	reqURL := fmt.Sprintf("%s/oauth2/token", r.apiEndpoint)
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(r.clientId, r.clientSecret)

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	token := ""
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return "", fmt.Errorf("Json decode error: %w", err)
	}

	return token, nil
}
func (r *Logingateway) GetAccountInfo(token string) (domain.User, error) {
	reqURL := fmt.Sprintf("%s/users/@me", r.apiEndpoint)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return domain.User{}, fmt.Errorf("Create reques error: %w", err)
	}

	req.Header.Set("User-Agent", "DiscordBot (http://localhost, 0.1)")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := r.client.Do(req)
	if err != nil {
		return domain.User{}, fmt.Errorf("Do request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return domain.User{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	accInfo := domain.User{}
	err = json.NewDecoder(resp.Body).Decode(&accInfo)
	if err != nil {
		return domain.User{}, fmt.Errorf("Json decode error: %w", err)
	}

	return accInfo, nil
}
