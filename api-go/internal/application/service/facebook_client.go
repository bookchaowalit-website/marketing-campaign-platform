package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// FacebookClient implements SocialPublishClient for Facebook Pages
type FacebookClient struct {
	appID     string
	appSecret string
	appVersion string
	baseURL   string
	httpClient *http.Client
}

// NewFacebookClient creates a new Facebook Graph API client
func NewFacebookClient(appID, appSecret string) *FacebookClient {
	return &FacebookClient{
		appID:      appID,
		appSecret:  appSecret,
		appVersion: "v18.0",
		baseURL:    "https://graph.facebook.com",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (f *FacebookClient) GetOAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", f.appID)
	params.Set("redirect_uri", "http://localhost:3000/settings/social/callback?platform=facebook")
	params.Set("state", state)
	params.Set("scope", "pages_manage_posts,pages_manage_metadata,pages_read_engagement,pages_read_user_content,instagram_basic,instagram_content_publish")
	return fmt.Sprintf("https://www.facebook.com/%s/dialog/oauth?%s", f.appVersion, params.Encode())
}

func (f *FacebookClient) ExchangeCode(code string) (accessToken, refreshToken string, expiresIn int, err error) {
	tokenURL := fmt.Sprintf("%s/%s/oauth/access_token", f.baseURL, f.appVersion)
	params := url.Values{}
	params.Set("client_id", f.appID)
	params.Set("client_secret", f.appSecret)
	params.Set("redirect_uri", "http://localhost:3000/settings/social/callback?platform=facebook")
	params.Set("code", code)

	resp, err := f.httpClient.Get(tokenURL + "?" + params.Encode())
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", 0, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Get long-lived token (60 days)
	longLivedToken, expiresIn, err := f.getLongLivedToken(result.AccessToken)
	if err != nil {
		return result.AccessToken, "", result.ExpiresIn, nil // Return short-lived if exchange fails
	}

	return longLivedToken, "", expiresIn, nil
}

func (f *FacebookClient) getLongLivedToken(shortToken string) (string, int, error) {
	tokenURL := fmt.Sprintf("%s/oauth/access_token", f.baseURL)
	params := url.Values{}
	params.Set("grant_type", "fb_exchange_token")
	params.Set("client_id", f.appID)
	params.Set("client_secret", f.appSecret)
	params.Set("fb_exchange_token", shortToken)

	resp, err := f.httpClient.Get(tokenURL + "?" + params.Encode())
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", 0, err
	}

	return result.AccessToken, result.ExpiresIn, nil
}

func (f *FacebookClient) GetAccountInfo(accessToken string) (accountID, accountName string, err error) {
	// Get user info
	userURL := fmt.Sprintf("%s/%s/me?fields=id,name", f.baseURL, f.appVersion)
	req, _ := http.NewRequest("GET", userURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userResult struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userResult); err != nil {
		return "", "", err
	}

	return userResult.ID, userResult.Name, nil
}

func (f *FacebookClient) GetUserPages(accessToken string) ([]struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Token  string `json:"access_token"`
}, error) {
	pagesURL := fmt.Sprintf("%s/%s/me/accounts", f.baseURL, f.appVersion)
	req, _ := http.NewRequest("GET", pagesURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get pages: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Token string `json:"access_token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (f *FacebookClient) PublishPost(accessToken, content string, mediaURLs []string) (postID string, err error) {
	return f.publishToPage(accessToken, "", content, mediaURLs)
}

func (f *FacebookClient) publishToPage(accessToken, pageID, content string, mediaURLs []string) (string, error) {
	pageIDToUse := pageID
	if pageIDToUse == "" {
		// Get first page if not specified
		pages, err := f.GetUserPages(accessToken)
		if err != nil || len(pages) == 0 {
			return "", fmt.Errorf("no pages found")
		}
		pageIDToUse = pages[0].ID
		accessToken = pages[0].Token
	}

	data := map[string]string{
		"message": content,
	}

	// Handle media (photo)
	if len(mediaURLs) > 0 {
		// For URLs, we need to use URL parameter
		data["url"] = mediaURLs[0]
		delete(data, "message") // Facebook API uses 'message' only for text-only posts
		data["caption"] = content
	}

	jsonData, _ := json.Marshal(data)
	postURL := fmt.Sprintf("%s/%s/%s/feed", f.baseURL, f.appVersion, pageIDToUse)
	if len(mediaURLs) > 0 {
		postURL = fmt.Sprintf("%s/%s/%s/photos", f.baseURL, f.appVersion, pageIDToUse)
	}

	resp, err := f.httpClient.Post(postURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to publish: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("publish failed: %s", string(body))
	}

	var result struct {
		ID string `json:"id"`
	}
	json.Unmarshal(body, &result)
	return result.ID, nil
}

func (f *FacebookClient) SchedulePost(accessToken, content string, mediaURLs []string, scheduledAt string) (postID string, err error) {
	// Facebook supports scheduled posts via 'scheduled_publish_time' parameter
	pages, err := f.GetUserPages(accessToken)
	if err != nil || len(pages) == 0 {
		return "", fmt.Errorf("no pages found")
	}

	page := pages[0]
	t, err := time.Parse(time.RFC3339, scheduledAt)
	if err != nil {
		return "", fmt.Errorf("invalid scheduled time: %w", err)
	}

	data := map[string]interface{}{
		"message":              content,
		"scheduled_publish_time": t.Unix(),
		"published":            false,
	}

	jsonData, _ := json.Marshal(data)
	postURL := fmt.Sprintf("%s/%s/%s/feed", f.baseURL, f.appVersion, page.ID)

	resp, err := f.httpClient.Post(postURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to schedule: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("schedule failed: %s", string(body))
	}

	var result struct {
		ID string `json:"id"`
	}
	json.Unmarshal(body, &result)
	return result.ID, nil
}

func (f *FacebookClient) DeletePost(accessToken, postID string) error {
	deleteURL := fmt.Sprintf("%s/%s/%s", f.baseURL, f.appVersion, postID)
	req, _ := http.NewRequest("DELETE", deleteURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (f *FacebookClient) GetAnalytics(accessToken, postID string) (*struct {
	PostID       string
	Impressions  int
	Reach        int
	Likes        int
	Comments     int
	Shares       int
	Clicks       int
	EngagementRate float64
}, error) {
	// Get post insights
	insightsURL := fmt.Sprintf("%s/%s/%s/insights?metric=post_impressions,post_reach,post_impressions_unique,post_engagements", f.baseURL, f.appVersion, postID)
	req, _ := http.NewRequest("GET", insightsURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get insights: %w", err)
	}
	defer resp.Body.Close()

	// Parse insights data (simplified - real implementation would be more complex)
	var insightsResult struct {
		Data []struct {
			Name   string `json:"name"`
			Values []struct {
				Value int `json:"value"`
			} `json:"values"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&insightsResult); err != nil {
		return nil, err
	}

	analytics := &struct {
		PostID       string
		Impressions  int
		Reach        int
		Likes        int
		Comments     int
		Shares       int
		Clicks       int
		EngagementRate float64
	}{PostID: postID}

	// Map insights to analytics
	for _, insight := range insightsResult.Data {
		if len(insight.Values) > 0 {
			switch insight.Name {
			case "post_impressions":
				analytics.Impressions = insight.Values[0].Value
			case "post_reach":
				analytics.Reach = insight.Values[0].Value
			case "post_engagements":
				analytics.Likes = insight.Values[0].Value // Simplified
			}
		}
	}

	// Calculate engagement rate
	if analytics.Impressions > 0 {
		analytics.EngagementRate = float64(analytics.Likes) / float64(analytics.Impressions) * 100
	}

	return analytics, nil
}

func (f *FacebookClient) RefreshToken(refreshToken string) (newAccessToken string, newExpiresIn int, err error) {
	// Facebook doesn't use refresh tokens - they use long-lived tokens that expire in 60 days
	// To extend a long-lived token, call the same endpoint with the current token
	return "", 0, fmt.Errorf("facebook uses long-lived tokens, use GetLongLivedToken instead")
}

// RefreshLongLivedToken extends a long-lived Facebook token (60 days)
func (f *FacebookClient) RefreshLongLivedToken(currentToken string) (newToken string, expiresIn int, err error) {
	tokenURL := fmt.Sprintf("%s/oauth/access_token", f.baseURL)
	params := url.Values{}
	params.Set("grant_type", "fb_exchange_token")
	params.Set("client_id", f.appID)
	params.Set("client_secret", f.appSecret)
	params.Set("fb_exchange_token", currentToken)

	resp, err := f.httpClient.Get(tokenURL + "?" + params.Encode())
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", 0, err
	}

	return result.AccessToken, result.ExpiresIn, nil
}
