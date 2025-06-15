package api

import (
	"crypto/x509"
	"crypto/rsa"
	"encoding/pem"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	//"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type GitHubAppAuth struct {
	AppID          string
	InstallationID string
	PrivateKeyPath string
	tokenCache     string
	tokenExpires   time.Time
}

func (gh *GitHubAppAuth) loadPrivateKey() (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(gh.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func (gh *GitHubAppAuth) generateJWT() (string, error) {
	key, err := gh.loadPrivateKey()
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"iss": gh.AppID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 9).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

func (gh *GitHubAppAuth) GetAccessToken() (string, error) {
	if time.Now().Before(gh.tokenExpires) && gh.tokenCache != "" {
		return gh.tokenCache, nil
	}

	jwtToken, err := gh.generateJWT()
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", gh.InstallationID)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expires_at"`
	}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	parsedExp, err := time.Parse(time.RFC3339, response.ExpiresAt)
	if err != nil {
		return "", err
	}
	gh.tokenCache = response.Token
	gh.tokenExpires = parsedExp.Add(-1 * time.Minute) // refresh buffer

	return response.Token, nil
}

