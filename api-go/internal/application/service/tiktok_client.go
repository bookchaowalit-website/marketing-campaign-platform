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

// TikTokClient implements SocialPublishClient for TikTok
type TikTokClient struct {
	clientID     string
	clientSecret string
	redirectURL  string
	baseURL      string
	httpClient   *http.Client
}

// NewTikTokClient creates a new TikTok API client
func NewTikTokClient(clientID, clientSecret, redirectURL string) *TikTokClient {
	return &TikTokClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		baseURL:      "https://open.tiktokapis.com/v2",
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (t *TikTokClient) GetOAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_key", t.clientID)
	params.Set("redirect_uri", t.redirectURL)
	params.Set("state", state)
	params.Set("response_type", "code")
	params.Set("scope", "user.info.basic,video.publish,video.upload")
	return fmt.Sprintf("https://www.tiktok.com/v2/auth/authorize/?%s", params.Encode())
}

func (t *TikTokClient) ExchangeCode(code string) (accessToken, refreshToken string, expiresIn int, err error) {
	tokenURL := fmt.Sprintf("%s/oauth/token/", t.baseURL)
	data := url.Values{}
	data.Set("client_key", t.clientID)
	data.Set("client_secret", t.clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", t.redirectURL)

	resp, err := t.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		OpenID       string `json:"open_id"`
		Error        struct {
			Code        string `json:"code"`
			Message     string `json:"message"`
			LogID       string `json:"log_id"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", 0, fmt.Errorf("failed to decode token response: %w", err)
	}

	if result.Error.Code != "" {
		return "", "", 0, fmt.Errorf("tiktok error: %s - %s", result.Error.Code, result.Error.Message)
	}

	return result.AccessToken, result.RefreshToken, result.ExpiresIn, nil
}

func (t *TikTokClient) GetAccountInfo(accessToken string) (accountID, accountName string, err error) {
	userURL := fmt.Sprintf("%s/user/info/", t.baseURL)
	req, _ := http.NewRequest("GET", userURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			User struct {
				OpenID      string `json:"open_id"`
				UnionID     string `json:"union_id"`
				DisplayName string `json:"display_name"`
				Username    string `json:"username"`
			} `json:"user"`
		} `json:"data"`
		Error *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	if result.Error != nil {
		return "", "", fmt.Errorf("tiktok error: %s - %s", result.Error.Code, result.Error.Message)
	}

	accountID = result.Data.User.OpenID
	accountName = result.Data.User.DisplayName
	if accountName == "" {
		accountName = result.Data.User.Username
	}

	return accountID, accountName, nil
}

func (t *TikTokClient) PublishPost(accessToken, content string, mediaURLs []string) (postID string, err error) {
	// TikTok requires video content for publishing
	// Step 1: Initialize video upload
	initURL := fmt.Sprintf("%s/post/publish/video/init/", t.baseURL)

	postInfo := map[string]interface{}{
		"title":       content,
		"privacy_level": "MUTUAL_FOLLOW_FRIENDS",
		"disable_comment": false,
		"disable_duet": true,
		"disable_stitch": true,
	}
	if len(mediaURLs) > 0 {
		postInfo["video_cover"] = map[string]string{
			"uri": mediaURLs[0],
		}
	}

	jsonData, _ := json.Marshal(postInfo)
	req, _ := http.NewRequest("POST", initURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to initialize video upload: %w", err)
	}
	defer resp.Body.Close()

	var initResult struct {
		Data struct {
			PublishID string `json:"publish_id"`
		} `json:"data"`
		Error *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&initResult); err != nil {
		return "", err
	}

	if initResult.Error != nil {
		return "", fmt.Errorf("tiktok init error: %s - %s", initResult.Error.Code, initResult.Error.Message)
	}

	// Step 2: Upload video (would require actual video upload in production)
	// For now, return the publish_id as the post identifier
	return initResult.Data.PublishID, nil
}

func (t *TikTokClient) SchedulePost(accessToken, content string, mediaURLs []string, scheduledAt string) (postID string, err error) {
	// TikTok supports scheduled posts via schedule_time parameter
	initURL := fmt.Sprintf("%s/post/publish/video/init/", t.baseURL)

	scheduleTime, err := time.Parse(time.RFC3339, scheduledAt)
	if err != nil {
		return "", fmt.Errorf("invalid scheduled time: %w", err)
	}

	postInfo := map[string]interface{}{
		"title":         content,
		"privacy_level": "MUTUAL_FOLLOW_FRIENDS",
		"schedule_time": scheduleTime.Unix(),
	}

	jsonData, _ := json.Marshal(postInfo)
	req, _ := http.NewRequest("POST", initURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to schedule post: %w", err)
	}
	defer resp.Body.Close()

	var initResult struct {
		Data struct {
			PublishID string `json:"publish_id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&initResult); err != nil {
		return "", err
	}

	return initResult.Data.PublishID, nil
}

func (t *TikTokClient) DeletePost(accessToken, postID string) error {
	// TikTok API doesn't support post deletion via API
	return fmt.Errorf("tiktok does not support post deletion via API")
}

func (t *TikTokClient) GetAnalytics(accessToken, postID string) (*domain.SocialPostAnalytics, error) {
	// Get video analytics
	analyticsURL := fmt.Sprintf("%s/video/analytic/", t.baseURL)
	params := url.Values{}
	params.Set("video_id", postID)
	params.Set("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	params.Set("end_date", time.Now().Format("2006-01-02"))

	req, _ := http.NewRequest("GET", analyticsURL+"?"+params.Encode(), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data struct {
			VideoID   string  `json:"video_id"`
			Views     int     `json:"total_play_count"`
			Likes     int     `json:"total_like_count"`
			Comments  int     `json:"total_comment_count"`
			Shares    int     `json:"total_share_count"`
			Reach     int     `json:"total_unique_viewers"`
			WatchTime float64 `json:"total_watch_time_in_seconds"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	analytics := &domain.SocialPostAnalytics{
		PostID:      postID,
		Platform:    "tiktok",
		Impressions: result.Data.Views,
		Reach:       result.Data.Reach,
		Likes:       result.Data.Likes,
		Comments:    result.Data.Comments,
		Shares:      result.Data.Shares,
		Clicks:      0,
		FetchedAt:   time.Now(),
	}

	// Calculate engagement rate
	if analytics.Impressions > 0 {
		analytics.EngagementRate = float64(analytics.Likes+analytics.Comments+analytics.Shares) / float64(analytics.Impressions) * 100
	}

	return analytics, nil
}

func (t *TikTokClient) RefreshToken(refreshToken string) (newAccessToken string, newExpiresIn int, err error) {
	tokenURL := fmt.Sprintf("%s/oauth/token/", t.baseURL)
	data := url.Values{}
	data.Set("client_key", t.clientID)
	data.Set("client_secret", t.clientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	resp, err := t.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return "", 0, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Error        struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", 0, err
	}

	if result.Error.Code != "" {
		return "", 0, fmt.Errorf("tiktok error: %s - %s", result.Error.Code, result.Error.Message)
	}

	return result.AccessToken, result.ExpiresIn, nil
}
