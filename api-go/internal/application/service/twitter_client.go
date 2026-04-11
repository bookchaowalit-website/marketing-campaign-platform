package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// TwitterClient implements SocialPublishClient for Twitter/X
type TwitterClient struct {
	apiKey       string
	apiSecret    string
	bearerToken  string
	redirectURL  string
	baseURL      string
	uploadURL    string
	httpClient   *http.Client
}

// NewTwitterClient creates a new Twitter/X API client
func NewTwitterClient(apiKey, apiSecret, bearerToken, redirectURL string) *TwitterClient {
	return &TwitterClient{
		apiKey:      apiKey,
		apiSecret:   apiSecret,
		bearerToken: bearerToken,
		redirectURL: redirectURL,
		baseURL:     "https://api.twitter.com/2",
		uploadURL:   "https://upload.twitter.com/1.1",
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (t *TwitterClient) GetOAuthURL(state string) string {
	// Twitter/X OAuth 2.0 PKCE flow
	return fmt.Sprintf("https://twitter.com/i/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=tweet.read+tweet.write+users.read&state=%s",
		t.apiKey, url.QueryEscape(t.redirectURL), state)
}

func (t *TwitterClient) ExchangeCode(code string) (accessToken, refreshToken string, expiresIn int, err error) {
	tokenURL := "https://api.twitter.com/2/oauth2/token"

	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", t.apiKey)
	data.Set("redirect_uri", t.redirectURL)
	data.Set("code_verifier", "challenge") // In production, use actual PKCE code verifier

	req, _ := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(t.apiKey, t.apiSecret)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
		TokenType    string `json:"token_type"`
		Error        string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", 0, fmt.Errorf("failed to decode token response: %w", err)
	}

	if result.Error != "" {
		return "", "", 0, fmt.Errorf("twitter oauth error: %s", result.Error)
	}

	return result.AccessToken, result.RefreshToken, result.ExpiresIn, nil
}

func (t *TwitterClient) GetAccountInfo(accessToken string) (accountID, accountName string, err error) {
	userURL := fmt.Sprintf("%s/users/me?user.fields=name,username,profile_image_url", t.baseURL)
	req, _ := http.NewRequest("GET", userURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	return result.Data.ID, result.Data.Name, nil
}

func (t *TwitterClient) PublishPost(accessToken, content string, mediaURLs []string) (postID string, err error) {
	// Handle media uploads first
	var mediaIDs []string
	for _, mediaURL := range mediaURLs {
		mediaID, err := t.uploadMedia(accessToken, mediaURL)
		if err != nil {
			return "", fmt.Errorf("failed to upload media: %w", err)
		}
		mediaIDs = append(mediaIDs, mediaID)
	}

	// Create tweet
	tweetURL := fmt.Sprintf("%s/tweets", t.baseURL)

	tweetData := map[string]interface{}{
		"text": content,
	}

	if len(mediaIDs) > 0 {
		tweetData["media"] = map[string]interface{}{
			"media_ids": mediaIDs,
		}
	}

	jsonData, _ := json.Marshal(tweetData)
	req, _ := http.NewRequest("POST", tweetURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to publish tweet: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("publish failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Data.ID, nil
}

func (t *TwitterClient) SchedulePost(accessToken, content string, mediaURLs []string, scheduledAt string) (postID string, err error) {
	// Twitter/X API v2 doesn't natively support scheduling
	// Return the publish info and note that scheduling should be handled by our scheduler
	return t.PublishPost(accessToken, content, mediaURLs)
}

func (t *TwitterClient) uploadMedia(accessToken, mediaURL string) (string, error) {
	// For simplicity, this downloads the media and uploads to Twitter
	// In production, use chunked upload for large files

	// Download media
	resp, err := http.Get(mediaURL)
	if err != nil {
		return "", fmt.Errorf("failed to download media: %w", err)
	}
	defer resp.Body.Close()

	mediaData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read media: %w", err)
	}

	// Upload to Twitter
	uploadURL := fmt.Sprintf("%s/media/upload.json", t.uploadURL)
	req, _ := http.NewRequest("POST", uploadURL, bytes.NewBuffer(mediaData))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err = t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload media: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		MediaID     string `json:"media_id_string"`
		MediaIDInt  int64  `json:"media_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.MediaID != "" {
		return result.MediaID, nil
	}
	return fmt.Sprintf("%d", result.MediaIDInt), nil
}

func (t *TwitterClient) DeletePost(accessToken, postID string) error {
	deleteURL := fmt.Sprintf("%s/tweets/%s", t.baseURL, postID)
	req, _ := http.NewRequest("DELETE", deleteURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete tweet: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (t *TwitterClient) GetAnalytics(accessToken, postID string) (*domain.SocialPostAnalytics, error) {
	// Get tweet metrics using recent search or tweet lookup
	metricsURL := fmt.Sprintf("%s/tweets?ids=%s&tweet.fields=public_metrics,non_public_metrics,organic_metrics", t.baseURL, postID)
	req, _ := http.NewRequest("GET", metricsURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweet metrics: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
			PublicMetrics struct {
				Retweets    int `json:"retweet_count"`
				Replies     int `json:"reply_count"`
				Likes       int `json:"like_count"`
				Quotes      int `json:"quote_count"`
				Impressions int `json:"impression_count"`
			} `json:"public_metrics"`
			OrganicMetrics *struct {
				Impressions   int `json:"impression_count"`
				URLLinkClicks int `json:"url_link_clicks"`
				UserProfileClicks int `json:"user_profile_clicks"`
			} `json:"organic_metrics"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no metrics found for tweet %s", postID)
	}

	tweet := result.Data[0]
	analytics := &domain.SocialPostAnalytics{
		PostID:      postID,
		Platform:    "twitter",
		Impressions: tweet.PublicMetrics.Impressions,
		Reach:       tweet.PublicMetrics.Impressions,
		Likes:       tweet.PublicMetrics.Likes,
		Comments:    tweet.PublicMetrics.Replies,
		Shares:      tweet.PublicMetrics.Retweets + tweet.PublicMetrics.Quotes,
		Clicks:      0,
		FetchedAt:   time.Now(),
	}

	if tweet.OrganicMetrics != nil {
		analytics.Clicks = tweet.OrganicMetrics.URLLinkClicks
		analytics.Impressions = tweet.OrganicMetrics.Impressions
	}

	// Calculate engagement rate
	if analytics.Impressions > 0 {
		totalEngagement := analytics.Likes + analytics.Comments + analytics.Shares + analytics.Clicks
		analytics.EngagementRate = float64(totalEngagement) / float64(analytics.Impressions) * 100
	}

	return analytics, nil
}

func (t *TwitterClient) RefreshToken(refreshToken string) (newAccessToken string, newExpiresIn int, err error) {
	tokenURL := "https://api.twitter.com/2/oauth2/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", t.apiKey)

	req, _ := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(t.apiKey, t.apiSecret)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Error        string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", 0, err
	}

	if result.Error != "" {
		return "", 0, fmt.Errorf("twitter refresh error: %s", result.Error)
	}

	return result.AccessToken, result.ExpiresIn, nil
}
