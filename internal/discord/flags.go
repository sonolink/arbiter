package discord

import (
	"encoding/json"
	"strconv"
)

type flagValue interface {
	~int | ~uint64
}

func HasAnyFlags[T flagValue](flags T, check ...T) bool {
	for _, c := range check {
		if flags&c == c {
			return true
		}
	}
	return false
}

func HasAllFlags[T flagValue](flags T, check ...T) bool {
	for _, c := range check {
		if flags&c != c {
			return false
		}
	}
	return true
}

type UserFlags int

const (
	UserFlagStaff                 UserFlags = 1 << 0  // Discord Employee
	UserFlagPremiumEarlySupporter UserFlags = 1 << 9  // Early Nitro Supporter
	UserFlagVerifiedBot           UserFlags = 1 << 16 // Verified Bot
	UserFlagVerifiedDeveloper     UserFlags = 1 << 17 // Early Verified Bot Developer
	UserFlagCertifiedModerator    UserFlags = 1 << 18 // Moderator Programs Alumni
)

type GuildMemberFlags int

const (
	GuildMemberFlagDidRejoin                  GuildMemberFlags = 1 << 0  // Member has left and rejoined the guild
	GuildMemberFlagCompletedOnboarding        GuildMemberFlags = 1 << 1  // Member has completed onboarding
	GuildMemberFlagBypassesVerification       GuildMemberFlags = 1 << 2  // Member is exempt from guild verification requirements
	GuildMemberFlagStartedOnboarding          GuildMemberFlags = 1 << 3  // Member has started onboarding
	GuildMemberFlagIsGuest                    GuildMemberFlags = 1 << 4  // Member is a guest and can only access the voice channel they were invited to
	GuildMemberFlagStartedHomeActions         GuildMemberFlags = 1 << 5  // Member has started Server Guide new member actions
	GuildMemberFlagCompletedHomeActions       GuildMemberFlags = 1 << 6  // Member has completed Server Guide new member actions
	GuildMemberFlagAutomodQuarantinedUsername GuildMemberFlags = 1 << 7  // Member's username, display name, or nickname is blocked by AutoMod
	GuildMemberFlagAutomodQuarantinedGuildTag GuildMemberFlags = 1 << 10 // Member's guild tag is blocked by AutoMod
)

type Permissions uint64

const (
	PermissionCreateInstantInvite              Permissions = 1 << 0  // Allows creation of instant invites
	PermissionKickMembers                      Permissions = 1 << 1  // Allows kicking members
	PermissionBanMembers                       Permissions = 1 << 2  // Allows banning members
	PermissionAdministrator                    Permissions = 1 << 3  // Allows all permissions and bypasses channel permission overwrites
	PermissionManageChannels                   Permissions = 1 << 4  // Allows management and editing of channels
	PermissionManageGuild                      Permissions = 1 << 5  // Allows management and editing of the guild
	PermissionAddReactions                     Permissions = 1 << 6  // Allows for adding new reactions to messages
	PermissionViewAuditLog                     Permissions = 1 << 7  // Allows for viewing of audit logs
	PermissionPrioritySpeaker                  Permissions = 1 << 8  // Allows for using priority speaker in a voice channel
	PermissionStream                           Permissions = 1 << 9  // Allows the user to go live
	PermissionViewChannel                      Permissions = 1 << 10 // Allows guild members to view a channel
	PermissionSendMessages                     Permissions = 1 << 11 // Allows for sending messages in a channel and creating threads in a forum
	PermissionSendTTSMessages                  Permissions = 1 << 12 // Allows for sending of /tts messages
	PermissionManageMessages                   Permissions = 1 << 13 // Allows for deletion of other users messages
	PermissionEmbedLinks                       Permissions = 1 << 14 // Links sent by users with this permission will be auto-embedded
	PermissionAttachFiles                      Permissions = 1 << 15 // Allows for uploading images and files
	PermissionReadMessageHistory               Permissions = 1 << 16 // Allows for reading of message history
	PermissionMentionEveryone                  Permissions = 1 << 17 // Allows for using the @everyone and @here tags
	PermissionUseExternalEmojis                Permissions = 1 << 18 // Allows the usage of custom emojis from other servers
	PermissionViewGuildInsights                Permissions = 1 << 19 // Allows for viewing guild insights
	PermissionConnect                          Permissions = 1 << 20 // Allows for joining of a voice channel
	PermissionSpeak                            Permissions = 1 << 21 // Allows for speaking in a voice channel
	PermissionMuteMembers                      Permissions = 1 << 22 // Allows for muting members in a voice channel
	PermissionDeafenMembers                    Permissions = 1 << 23 // Allows for deafening of members in a voice channel
	PermissionMoveMembers                      Permissions = 1 << 24 // Allows for moving of members between voice channels
	PermissionUseVAD                           Permissions = 1 << 25 // Allows for using voice-activity-detection in a voice channel
	PermissionChangeNickname                   Permissions = 1 << 26 // Allows for modification of own nickname
	PermissionManageNicknames                  Permissions = 1 << 27 // Allows for modification of other users nicknames
	PermissionManageRoles                      Permissions = 1 << 28 // Allows management and editing of roles
	PermissionManageWebhooks                   Permissions = 1 << 29 // Allows management and editing of webhooks
	PermissionManageGuildExpressions           Permissions = 1 << 30 // Allows for editing and deleting emojis, stickers, and soundboard sounds created by all users
	PermissionUseApplicationCommands           Permissions = 1 << 31 // Allows members to use application commands, including slash commands and context menu commands
	PermissionRequestToSpeak                   Permissions = 1 << 32 // Allows for requesting to speak in stage channels
	PermissionManageEvents                     Permissions = 1 << 33 // Allows for editing and deleting scheduled events created by all users
	PermissionManageThreads                    Permissions = 1 << 34 // Allows for deleting and archiving threads, and viewing all private threads
	PermissionCreatePublicThreads              Permissions = 1 << 35 // Allows for creating public and announcement threads
	PermissionCreatePrivateThreads             Permissions = 1 << 36 // Allows for creating private threads
	PermissionUseExternalStickers              Permissions = 1 << 37 // Allows the usage of custom stickers from other servers
	PermissionSendMessagesInThreads            Permissions = 1 << 38 // Allows for sending messages in threads
	PermissionUseEmbeddedActivities            Permissions = 1 << 39 // Allows for using Activities
	PermissionModerateMembers                  Permissions = 1 << 40 // Allows for timing out users
	PermissionViewCreatorMonetizationAnalytics Permissions = 1 << 41 // Allows for viewing role subscription insights
	PermissionUseSoundboard                    Permissions = 1 << 42 // Allows for using soundboard in a voice channel
	PermissionCreateGuildExpressions           Permissions = 1 << 43 // Allows for creating emojis, stickers, and soundboard sounds
	PermissionCreateEvents                     Permissions = 1 << 44 // Allows for creating scheduled events
	PermissionUseExternalSounds                Permissions = 1 << 45 // Allows the usage of custom soundboard sounds from other servers
	PermissionSendVoiceMessages                Permissions = 1 << 46 // Allows sending voice messages
	PermissionSetVoiceChannelStatus            Permissions = 1 << 48 // Allows setting voice channel status
	PermissionSendPolls                        Permissions = 1 << 49 // Allows sending polls
	PermissionUseExternalApps                  Permissions = 1 << 50 // Allows user-installed apps to send public responses
	PermissionPinMessages                      Permissions = 1 << 51 // Allows pinning and unpinning messages
	PermissionBypassSlowmode                   Permissions = 1 << 52 // Allows bypassing slowmode restrictions
)

func (p *Permissions) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}

	*p = Permissions(n)
	return nil
}

func (p Permissions) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(p), 10))
}
