package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type User struct {
	ID string `json:"id"`
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

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, parseAPIError(resp.StatusCode, body)
	}

	return body, nil
}
