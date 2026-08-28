package discord

import (
	"time"
)

type Guild struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	IsOwner     bool        `json:"owner"`
	Permissions Permissions `json:"permissions"`
	Features    []string    `json:"features"`
}

func (g *Guild) CreatedAt() (time.Time, error) {
	return SnowflakeTime(g.ID)
}

type Member struct {
	User         User             `json:"user"`
	Nick         *string          `json:"nick"`
	RolesIds     []string         `json:"roles"`
	JoinedAt     *time.Time       `json:"joined_at"`
	Pending      bool             `json:"pending"`
	PremiumSince *time.Time       `json:"premium_since"`
	Flags        GuildMemberFlags `json:"flags"`
}

func (m *Member) CreatedAt() (time.Time, error) {
	return SnowflakeTime(m.User.ID)
}

type User struct {
	ID         string    `json:"id"`
	Flags      UserFlags `json:"public_flags"`
	GlobalName *string   `json:"global_name"`
	Username   string    `json:"username"`
	MfaEnabled bool      `json:"mfa_enabled"`
	Locale     string    `json:"locale"`
	Verified   bool      `json:"verified"`
	Email      *string   `json:"email"`
}

func (u *User) CreatedAt() (time.Time, error) {
	return SnowflakeTime(u.ID)
}

type ConnectionType string

const (
	ConnectionTypeAmazonMusic ConnectionType = "amazon-music"
	ConnectionTypeBungie      ConnectionType = "bungie"
	ConnectionTypeBluesky     ConnectionType = "bluesky"
	ConnectionTypeCrunchyroll ConnectionType = "crunchyroll"
	ConnectionTypeDomain      ConnectionType = "domain"
	ConnectionTypeEbay        ConnectionType = "ebay"
	ConnectionTypeEpicGames   ConnectionType = "epicgames"
	ConnectionTypeFacebook    ConnectionType = "facebook"
	ConnectionTypeGitHub      ConnectionType = "github"
	ConnectionTypeMastodon    ConnectionType = "mastodon"
	ConnectionTypePayPal      ConnectionType = "paypal"
	ConnectionTypePlayStation ConnectionType = "playstation"
	ConnectionTypeReddit      ConnectionType = "reddit"
	ConnectionTypeRoblox      ConnectionType = "roblox"
	ConnectionTypeSpotify     ConnectionType = "spotify"
	ConnectionTypeSteam       ConnectionType = "steam"
	ConnectionTypeTikTok      ConnectionType = "tiktok"
	ConnectionTypeTwitch      ConnectionType = "twitch"
	ConnectionTypeTwitter     ConnectionType = "twitter"
	ConnectionTypeXbox        ConnectionType = "xbox"
	ConnectionTypeYouTube     ConnectionType = "youtube"
)

