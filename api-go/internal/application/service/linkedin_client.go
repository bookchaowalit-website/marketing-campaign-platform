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

// LinkedInClient implements SocialPublishClient for LinkedIn
type LinkedInClient struct {
	clientID     string
	clientSecret string
	redirectURL  string
	baseURL      string
	httpClient   *http.Client
}

// NewLinkedInClient creates a new LinkedIn API client
func NewLinkedInClient(clientID, clientSecret, redirectURL string) *LinkedInClient {
	return &LinkedInClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		baseURL:      "https://api.linkedin.com/v2",
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (l *LinkedInClient) GetOAuthURL(state string) string {
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", l.clientID)
	params.Set("redirect_uri", l.redirectURL)
	params.Set("state", state)
	params.Set("scope", "w_member_social r_liteprofile r_organization_social w_organization_social")
	return fmt.Sprintf("https://www.linkedin.com/oauth/v2/authorization?%s", params.Encode())
}

func (l *LinkedInClient) ExchangeCode(code string) (accessToken, refreshToken string, expiresIn int, err error) {
	tokenURL := "https://www.linkedin.com/oauth/v2/accessToken"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", l.clientID)
	data.Set("client_secret", l.clientSecret)
	data.Set("redirect_uri", l.redirectURL)

	resp, err := l.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", 0, fmt.Errorf("failed to decode token response: %w", err)
	}

	if result.Error != "" {
		return "", "", 0, fmt.Errorf("linkedin oauth error: %s", result.Error)
	}

	return result.AccessToken, "", result.ExpiresIn, nil
}

func (l *LinkedInClient) GetAccountInfo(accessToken string) (accountID, accountName string, err error) {
	userURL := fmt.Sprintf("%s/me?projection=(id,firstName,lastName,profilePicture)", l.baseURL)
	req, _ := http.NewRequest("GET", userURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		ID        string `json:"id"`
		FirstName struct {
			Localized map[string]string `json:"localized"`
		} `json:"firstName"`
		LastName struct {
			Localized map[string]string `json:"localized"`
		} `json:"lastName"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	accountID = result.ID
	firstName := ""
	lastName := ""
	if len(result.FirstName.Localized) > 0 {
		for _, v := range result.FirstName.Localized {
			firstName = v
			break
		}
	}
	if len(result.LastName.Localized) > 0 {
		for _, v := range result.LastName.Localized {
			lastName = v
			break
		}
	}
	accountName = firstName + " " + lastName

	return accountID, accountName, nil
}

// GetOrganizationPages returns organizations the user can manage
func (l *LinkedInClient) GetOrganizationPages(accessToken string) ([]struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}, error) {
	orgsURL := fmt.Sprintf("%s/organizationalEntityAcls?q=roleAssignee&role=ADMINISTRATOR&projection=(elements*(*,organizationalTarget~(localizedName,vanityName)))", l.baseURL)
	req, _ := http.NewRequest("GET", orgsURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Elements []struct {
			OrganizationalTarget string `json:"organizationalTarget"`
			OrgInfo              struct {
				LocalizedName string `json:"localizedName"`
			} `json:"organizationalTarget~"`
		} `json:"elements"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&resp.Body); err != nil {
		return nil, err
	}

	var orgs []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	for _, el := range result.Elements {
		orgID := el.OrganizationalTarget
		if idx := lastIndexOf(orgID, ":"); idx >= 0 && idx+1 < len(orgID) {
			orgID = orgID[idx+1:]
		}
		name := el.OrgInfo.LocalizedName
		orgs = append(orgs, struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{ID: orgID, Name: name})
	}

	return orgs, nil
}

func (l *LinkedInClient) PublishPost(accessToken, content string, mediaURLs []string) (postID string, err error) {
	// LinkedIn uses the person URN for personal posts
	// First get the user's URN
	userURN, err := l.getUserURN(accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to get user URN: %w", err)
	}

	// Build post content
	postData := map[string]interface{}{
		"author":    userURN,
		"lifecycleState": "PUBLISHED",
		"specificContent": map[string]interface{}{
			"com.linkedin.ugc.ShareContent": map[string]interface{}{
				"shareCommentary": map[string]interface{}{
					"text": content,
				},
				"shareMediaCategory": "NONE",
			},
		},
		"visibility": map[string]interface{}{
			"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC",
		},
	}

	// Handle media if provided
	if len(mediaURLs) > 0 {
		mediaID, err := l.uploadImage(accessToken, mediaURLs[0])
		if err == nil && mediaID != "" {
			postData["specificContent"] = map[string]interface{}{
				"com.linkedin.ugc.ShareContent": map[string]interface{}{
					"shareCommentary": map[string]interface{}{
						"text": content,
					},
					"shareMediaCategory": "IMAGE",
					"media": []map[string]interface{}{
						{
							"status": "READY",
							"description": map[string]interface{}{
								"text": content,
							},
							"originalUrl": mediaURLs[0],
							"media":       mediaID,
						},
					},
				},
			}
		}
	}

	jsonData, _ := json.Marshal(postData)
	postURL := fmt.Sprintf("%s/ugcPosts", l.baseURL)
	req, _ := http.NewRequest("POST", postURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to publish post: %w", err)
	}
	defer resp.Body.Close()

	// LinkedIn returns the URN in the Location header
	location := resp.Header.Get("X-RestLi-Id")
	if location == "" {
		location = resp.Header.Get("Location")
	}

	if location != "" {
		// Extract post ID from URN (urn:li:share:123456789 -> 123456789)
		if idx := lastIndexOf(location, ":"); idx >= 0 && idx+1 < len(location) {
			return location[idx+1:], nil
		}
		return location, nil
	}

	// Read body for error
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("publish failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return location, nil
}

func (l *LinkedInClient) SchedulePost(accessToken, content string, mediaURLs []string, scheduledAt string) (postID string, err error) {
	// LinkedIn supports scheduled posts via the publishAt field
	userURN, err := l.getUserURN(accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to get user URN: %w", err)
	}

	scheduleTime, err := time.Parse(time.RFC3339, scheduledAt)
	if err != nil {
		return "", fmt.Errorf("invalid scheduled time: %w", err)
	}

	postData := map[string]interface{}{
		"author":    userURN,
		"lifecycleState": "PUBLISHED",
		"specificContent": map[string]interface{}{
			"com.linkedin.ugc.ShareContent": map[string]interface{}{
				"shareCommentary": map[string]interface{}{
					"text": content,
				},
				"shareMediaCategory": "NONE",
			},
		},
		"visibility": map[string]interface{}{
			"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC",
		},
		"distribution": map[string]interface{}{
			"feedDistribution": "MAIN_FEED",
			"publishAt": scheduleTime.UnixMilli(),
		},
	}

	jsonData, _ := json.Marshal(postData)
	postURL := fmt.Sprintf("%s/posts", l.baseURL)
	req, _ := http.NewRequest("POST", postURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to schedule post: %w", err)
	}
	defer resp.Body.Close()

	location := resp.Header.Get("X-RestLi-Id")
	if location == "" {
		location = resp.Header.Get("Location")
	}

	if location != "" {
		if idx := lastIndexOf(location, ":"); idx >= 0 && idx+1 < len(location) {
			return location[idx+1:], nil
		}
		return location, nil
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("schedule failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return location, nil
}

func (l *LinkedInClient) uploadImage(accessToken, imageURL string) (string, error) {
	// Register upload
	registerURL := fmt.Sprintf("%s/images?action=registerUpload", l.baseURL)

	registerData := map[string]interface{}{
		"registerUploadRequest": map[string]interface{}{
			"recipes": []string{"urn:li:digitalmediaRecipe:feedshare-image"},
			"owner":   "urn:li:person:me",
			"serviceRelationships": []map[string]interface{}{
				{
					"relationshipType": "OWNER",
					"identifier":       "urn:li:userGeneratedContent",
				},
			},
		},
	}

	jsonData, _ := json.Marshal(registerData)
	req, _ := http.NewRequest("POST", registerURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to register upload: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Value struct {
			UploadMechanism map[string]interface{} `json:"uploadMechanism"`
			MediaArtifact   string                 `json:"mediaArtifact"`
			Image           string                 `json:"image"`
		} `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// In production, upload the actual image data to the upload URL
	// For now, return the image URN
	return result.Value.Image, nil
}

func (l *LinkedInClient) getUserURN(accessToken string) (string, error) {
	userURL := fmt.Sprintf("%s/me", l.baseURL)
	req, _ := http.NewRequest("GET", userURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return fmt.Sprintf("urn:li:person:%s", result.ID), nil
}

func (l *LinkedInClient) DeletePost(accessToken, postID string) error {
	deleteURL := fmt.Sprintf("%s/ugcPosts/%s", l.baseURL, postID)
	req, _ := http.NewRequest("DELETE", deleteURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (l *LinkedInClient) GetAnalytics(accessToken, postID string) (*domain.SocialPostAnalytics, error) {
	// Get post statistics
	shareStatURL := fmt.Sprintf("%s/socialActions/%s", l.baseURL, postID)
	req, _ := http.NewRequest("GET", shareStatURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get post stats: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		ID       string `json:"id"`
		Likes    struct {
			Paging struct {
				Total int `json:"total"`
			} `json:"paging"`
		} `json:"likes"`
		Comments struct {
			Paging struct {
				Total int `json:"total"`
			} `json:"paging"`
		} `json:"comments"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Get impressions from post lookup
	postURL := fmt.Sprintf("%s/ugcPosts/%s", l.baseURL, postID)
	req, _ = http.NewRequest("GET", postURL+"?projection=(id,statistics)", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp2, err := l.httpClient.Do(req)
	if err != nil {
		// If we can't get impressions, still return what we have
		analytics := &domain.SocialPostAnalytics{
			PostID:      postID,
			Platform:    "linkedin",
			Impressions: 0,
			Reach:       0,
			Likes:       result.Likes.Paging.Total,
			Comments:    result.Comments.Paging.Total,
			Shares:      0,
			Clicks:      0,
			FetchedAt:   time.Now(),
		}
		return analytics, nil
	}
	defer resp2.Body.Close()

	var postResult struct {
		Statistics struct {
			Impressions struct {
				Count int `json:"count"`
			} `json:"impressionCount"`
			Clicks struct {
				Count int `json:"count"`
			} `json:"clickCount"`
			Shares int `json:"shareCount"`
		} `json:"statistics"`
	}
	json.NewDecoder(resp2.Body).Decode(&postResult)

	analytics := &domain.SocialPostAnalytics{
		PostID:      postID,
		Platform:    "linkedin",
		Impressions: postResult.Statistics.Impressions.Count,
		Reach:       postResult.Statistics.Impressions.Count,
		Likes:       result.Likes.Paging.Total,
		Comments:    result.Comments.Paging.Total,
		Shares:      postResult.Statistics.Shares,
		Clicks:      postResult.Statistics.Clicks.Count,
		FetchedAt:   time.Now(),
	}

	if analytics.Impressions > 0 {
		totalEngagement := analytics.Likes + analytics.Comments + analytics.Shares + analytics.Clicks
		analytics.EngagementRate = float64(totalEngagement) / float64(analytics.Impressions) * 100
	}

	return analytics, nil
}

func (l *LinkedInClient) RefreshToken(refreshToken string) (newAccessToken string, newExpiresIn int, err error) {
	tokenURL := "https://www.linkedin.com/oauth/v2/accessToken"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", l.clientID)
	data.Set("client_secret", l.clientSecret)

	resp, err := l.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return "", 0, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", 0, err
	}

	if result.Error != "" {
		return "", 0, fmt.Errorf("linkedin refresh error: %s", result.Error)
	}

	return result.AccessToken, result.ExpiresIn, nil
}

// lastIndexOf returns the index of the last occurrence of substr in s
func lastIndexOf(s string, substr string) int {
	for i := len(s) - len(substr); i >= 0; i-- {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
