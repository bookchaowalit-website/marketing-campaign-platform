package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// InstagramClient implements SocialPublishClient for Instagram Business
type InstagramClient struct {
	facebookClient *FacebookClient
}

// NewInstagramClient creates a new Instagram Business API client
func NewInstagramClient(appID, appSecret string) *InstagramClient {
	return &InstagramClient{
		facebookClient: NewFacebookClient(appID, appSecret),
	}
}

func (i *InstagramClient) GetOAuthURL(state string) string {
	// Instagram uses the same OAuth flow as Facebook
	return i.facebookClient.GetOAuthURL(state)
}

func (i *InstagramClient) ExchangeCode(code string) (accessToken, refreshToken string, expiresIn int, err error) {
	// Same as Facebook since Instagram Business uses Facebook's OAuth
	return i.facebookClient.ExchangeCode(code)
}

func (i *InstagramClient) GetAccountInfo(accessToken string) (accountID, accountName string, err error) {
	// Get Facebook user pages, then get Instagram account for each page
	pages, err := i.facebookClient.GetUserPages(accessToken)
	if err != nil || len(pages) == 0 {
		return "", "", fmt.Errorf("no pages found")
	}

	page := pages[0]
	igAccountURL := fmt.Sprintf("%s/v18.0/%s?fields=instagram_business_account", i.facebookClient.baseURL, page.ID)
	req, _ := http.NewRequest("GET", igAccountURL, nil)
	req.Header.Set("Authorization", "Bearer "+page.Token)

	resp, err := i.facebookClient.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to get IG account: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		InstagramBusinessAccount struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"instagram_business_account"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	if result.InstagramBusinessAccount.ID == "" {
		return "", "", fmt.Errorf("no Instagram Business account connected to this page")
	}

	return result.InstagramBusinessAccount.ID, result.InstagramBusinessAccount.Username, nil
}

func (i *InstagramClient) PublishPost(accessToken, content string, mediaURLs []string) (postID string, err error) {
	// Get Instagram account ID
	igAccountID, err := i.getInstagramAccountID(accessToken)
	if err != nil {
		return "", err
	}

	if len(mediaURLs) == 0 {
		return "", fmt.Errorf("Instagram requires media content")
	}

	// Step 1: Create media container
	containerID, err := i.createMediaContainer(accessToken, igAccountID, mediaURLs[0], content)
	if err != nil {
		return "", err
	}

	// Step 2: Publish container
	return i.publishMediaContainer(accessToken, igAccountID, containerID)
}

func (i *InstagramClient) getInstagramAccountID(accessToken string) (string, error) {
	pages, err := i.facebookClient.GetUserPages(accessToken)
	if err != nil || len(pages) == 0 {
		return "", fmt.Errorf("no pages found")
	}

	page := pages[0]
	igAccountURL := fmt.Sprintf("%s/v18.0/%s?fields=instagram_business_account.id", i.facebookClient.baseURL, page.ID)
	req, _ := http.NewRequest("GET", igAccountURL, nil)
	req.Header.Set("Authorization", "Bearer "+page.Token)

	resp, err := i.facebookClient.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		InstagramBusinessAccount struct {
			ID string `json:"id"`
		} `json:"instagram_business_account"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.InstagramBusinessAccount.ID == "" {
		return "", fmt.Errorf("no Instagram Business account found")
	}

	return result.InstagramBusinessAccount.ID, nil
}

func (i *InstagramClient) createMediaContainer(accessToken, igAccountID, mediaURL, caption string) (string, error) {
	containerURL := fmt.Sprintf("%s/v18.0/%s/media", i.facebookClient.baseURL, igAccountID)
	data := map[string]string{
		"image_url":  mediaURL,
		"caption":    caption,
	}

	jsonData, _ := json.Marshal(data)
	resp, err := i.facebookClient.httpClient.Post(containerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("create container failed: %s", string(body))
	}

	var result struct {
		ID string `json:"id"`
	}
	json.Unmarshal(body, &result)
	return result.ID, nil
}

func (i *InstagramClient) publishMediaContainer(accessToken, igAccountID, containerID string) (string, error) {
	publishURL := fmt.Sprintf("%s/v18.0/%s/media_publish", i.facebookClient.baseURL, igAccountID)
	data := map[string]string{
		"creation_id": containerID,
	}

	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", publishURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := i.facebookClient.httpClient.Do(req)
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

func (i *InstagramClient) SchedulePost(accessToken, content string, mediaURLs []string, scheduledAt string) (postID string, err error) {
	return "", fmt.Errorf("Instagram does not support scheduled publishing via API")
}

func (i *InstagramClient) DeletePost(accessToken, postID string) error {
	deleteURL := fmt.Sprintf("%s/v18.0/%s", i.facebookClient.baseURL, postID)
	req, _ := http.NewRequest("DELETE", deleteURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := i.facebookClient.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (i *InstagramClient) GetAnalytics(accessToken, postID string) (*struct {
	PostID        string
	Impressions   int
	Reach         int
	Likes         int
	Comments      int
	Saves         int
	Shares        int
	EngagementRate float64
}, error) {
	// Get post insights
	insightsURL := fmt.Sprintf("%s/v18.0/%s/insights?metric=impressions,reach,likes,comments,saves,shares,engagement", i.facebookClient.baseURL, postID)
	req, _ := http.NewRequest("GET", insightsURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := i.facebookClient.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get insights: %w", err)
	}
	defer resp.Body.Close()

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
		PostID        string
		Impressions   int
		Reach         int
		Likes         int
		Comments      int
		Saves         int
		Shares        int
		EngagementRate float64
	}{PostID: postID}

	for _, insight := range insightsResult.Data {
		if len(insight.Values) > 0 {
			switch insight.Name {
			case "impressions":
				analytics.Impressions = insight.Values[0].Value
			case "reach":
				analytics.Reach = insight.Values[0].Value
			case "likes":
				analytics.Likes = insight.Values[0].Value
			case "comments":
				analytics.Comments = insight.Values[0].Value
			case "saves":
				analytics.Saves = insight.Values[0].Value
			case "shares":
				analytics.Shares = insight.Values[0].Value
			case "engagement":
				analytics.Likes = insight.Values[0].Value // Simplified
			}
		}
	}

	if analytics.Impressions > 0 {
		analytics.EngagementRate = float64(analytics.Likes+analytics.Comments+analytics.Saves+analytics.Shares) / float64(analytics.Impressions) * 100
	}

	return analytics, nil
}

func (i *InstagramClient) RefreshToken(refreshToken string) (newAccessToken string, newExpiresIn int, err error) {
	return i.facebookClient.RefreshLongLivedToken(refreshToken)
}

// Instagram Reels support
func (i *InstagramClient) PublishReel(accessToken, videoURL, caption string, thumbnailURL string) (postID string, err error) {
	igAccountID, err := i.getInstagramAccountID(accessToken)
	if err != nil {
		return "", err
	}

	// Create reel container
	containerURL := fmt.Sprintf("%s/v18.0/%s/media", i.facebookClient.baseURL, igAccountID)
	data := map[string]string{
		"video_url":    videoURL,
		"caption":      caption,
		"media_type":   "REELS",
	}
	if thumbnailURL != "" {
		data["thumb_offset"] = "0"
	}

	jsonData, _ := json.Marshal(data)
	resp, err := i.facebookClient.httpClient.Post(containerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create reel: %w", err)
	}
	defer resp.Body.Close()

	var containerResult struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&containerResult); err != nil {
		return "", err
	}

	// Publish reel
	publishURL := fmt.Sprintf("%s/v18.0/%s/media_publish", i.facebookClient.baseURL, igAccountID)
	publishData := map[string]string{
		"creation_id": containerResult.ID,
	}
	jsonPublishData, _ := json.Marshal(publishData)
	req, _ := http.NewRequest("POST", publishURL, bytes.NewBuffer(jsonPublishData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	pubResp, err := i.facebookClient.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer pubResp.Close()

	var publishResult struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(pubResp.Body).Decode(&publishResult); err != nil {
		return "", err
	}

	return publishResult.ID, nil
}

// Instagram Stories support
func (i *InstagramClient) PublishStory(accessToken, mediaURL string, isVideo bool) (postID string, err error) {
	igAccountID, err := i.getInstagramAccountID(accessToken)
	if err != nil {
		return "", err
	}

	containerURL := fmt.Sprintf("%s/v18.0/%s/media", i.facebookClient.baseURL, igAccountID)
	data := map[string]string{
		"is_shared_to_feed": "false",
	}
	
	if isVideo {
		data["video_url"] = mediaURL
		data["media_type"] = "STORIES"
	} else {
		data["image_url"] = mediaURL
	}

	jsonData, _ := json.Marshal(data)
	resp, err := i.facebookClient.httpClient.Post(containerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create story: %w", err)
	}
	defer resp.Body.Close()

	var containerResult struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&containerResult); err != nil {
		return "", err
	}

	// Publish story
	publishURL := fmt.Sprintf("%s/v18.0/%s/media_publish", i.facebookClient.baseURL, igAccountID)
	publishData := map[string]string{
		"creation_id": containerResult.ID,
	}
	jsonPublishData, _ := json.Marshal(publishData)
	req, _ := http.NewRequest("POST", publishURL, bytes.NewBuffer(jsonPublishData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	pubResp, err := i.facebookClient.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer pubResp.Close()

	var publishResult struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(pubResp.Body).Decode(&publishResult); err != nil {
		return "", err
	}

	return publishResult.ID, nil
}
