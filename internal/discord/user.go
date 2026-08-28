package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Epoch is the discord epoch in milliseconds.
const Epoch = 1420070400000

func SnowflakeTime(id string) (time.Time, error) {
	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("discord: invalid snowflake %q: %w", id, err)
	}
	return time.UnixMilli(n>>22/1000 + Epoch), nil
}

type User struct {
	ID         string  `json:"id"`
	Flags      int     `json:"public_flags"`
	GlobalName *string `json:"global_name"`
	Username   string  `json:"username"`
	MfaEnabled bool    `json:"mfa_enabled"`
	Locale     string  `json:"locale"`
	Verified   bool    `json:"verified"`
	Email      *string `json:"email"`
}

func (u *User) CreatedAt() (time.Time, error) {
	return SnowflakeTime(u.ID)
}

func (c *Client) Me(ctx context.Context, accessToken string) (*User, error) {
	body, err := c.get(ctx, accessToken, "/users/@me")
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("discord: decoding user: %w (body: %.200q)", err, body)
	}

	return &user, nil
}

func (c *Client) GuildMember(ctx context.Context, accessToken, guildID string) (json.RawMessage, error) {
	body, err := c.get(
		ctx,
		accessToken,
		"/users/@me/guilds/"+url.PathEscape(guildID)+"/member",
	)
	if err != nil {
		return nil, err
	}

	if !json.Valid(body) {
		return nil, fmt.Errorf("discord: member response is not valid JSON (body: %.200q)", body)
	}

	return body, nil
}

func (c *Client) get(ctx context.Context, accessToken, path string) ([]byte, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, parseRateLimitError(resp, body, parseAPIError(resp.StatusCode, body))
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, parseAPIError(resp.StatusCode, body)
	}

	return body, nil
}
