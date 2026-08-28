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
	var snowflake = ((n >> 22) + Epoch) / 1000
	return time.UnixMilli(snowflake), nil
}

func (c *Client) Me(ctx context.Context, accessToken string) (*User, error) {
	body, err := c.get(ctx, accessToken, "/users/@me")
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("discord: decoding user (%s): %w (body: %.200q)", user.ID, err, body)
	}

	return &user, nil
}

func (c *Client) GuildMember(ctx context.Context, accessToken, guildID string) (*Member, error) {
	body, err := c.get(
		ctx,
		accessToken,
		"/users/@me/guilds/"+url.PathEscape(guildID)+"/member",
	)
	if err != nil {
		return nil, err
	}

	var member Member
	if err := json.Unmarshal(body, &member); err != nil {
		return nil, fmt.Errorf("discord: decoding member (%s): %w (body: %.200q)", member.User.ID, err, body)
	}

	return &member, nil
}

func (c *Client) Guilds(ctx context.Context, accessToken string) ([]Guild, error) {
	body, err := c.get(ctx, accessToken, "/users/@me/guilds")
	if err != nil {
		return nil, err
	}

	var guilds []Guild
	if err := json.Unmarshal(body, &guilds); err != nil {
		return nil, fmt.Errorf("discord: decoding guilds: %w (body: %.200q)", err, body)
	}

	return guilds, nil
}

func (c *Client) Connections(ctx context.Context, accessToken string) ([]Connection, error) {
	body, err := c.get(ctx, accessToken, "/users/@me/connections")
	if err != nil {
		return nil, err
	}

	var connections []Connection
	if err := json.Unmarshal(body, &connections); err != nil {
		return nil, fmt.Errorf("discord: decoding connections: %w (body: %.200q)", err, body)
	}

	return connections, nil
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
