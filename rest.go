package disgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/slf4go/logger"
)

type EndPoint struct {
	URL    string
	Bucket string
}

var (
	BaseUrl = "https://discordapp.com/api"

	EndPointBotGateway   = makeEndPoint("/gateway/bot")
	EndPointVoiceRegions = makeEndPoint("/voice/regions")

	EndPointChannel            = makeEndPoint("/channels/:channel_id")
	EndPointMessages           = makeEndPoint("/channels/:channel_id/messages")
	EndPointMessage            = makeEndPoint("/channels/:channel_id/messages/:message_id")
	EndPointMessageBulkDelete  = makeEndPoint("/channels/:channel_id/messages/bulk-delete")
	EndPointReactions          = makeEndPoint("/channels/:channel_id/messages/:mesasge_id/reactions")
	EndPointReaction           = makeEndPoint("/channels/:channel_id/messages/:mesasge_id/reactions/%s/:user_id")
	EndPointOwnReaction        = makeEndPoint("/channels/:channel_id/messages/:message_id/reactions/%s/@me")
	EndPointChannelPermissions = makeEndPoint("/channels/:channel_id/permissions/:overwrite_id")
	EndPointChannelInvites     = makeEndPoint("/channels/:channel_id/invites")
	EndPointChannelTyping      = makeEndPoint("/channels/:channel_id/typing")
	EndPointChannelPins        = makeEndPoint("/channels/:channel_id/pins")
	EndPointChannelPin         = makeEndPoint("/channels/:channel_id/pins/:message_id")
	EndPointChannelRecipient   = makeEndPoint("/channels/:channel_id/recipients/:user_id")

	EndPointGuilds               = makeEndPoint("/guilds")
	EndPointGuild                = makeEndPoint("/guilds/:guild_id")
	EndPointGuildChannels        = makeEndPoint("/guilds/:guild_id/channels")
	EndPointGuildMembers         = makeEndPoint("/guilds/:guild_id/members")
	EndPointGuildMember          = makeEndPoint("/guilds/:guild_id/members/:user_id")
	EndPointGuildOwnNick         = makeEndPoint("/guilds/:guild_id/members/@me/nick")
	EndPointGuildMemberRoles     = makeEndPoint("/guilds/:guild_id/members/:user_id/roles/:role_id")
	EndPointGuildBans            = makeEndPoint("/guilds/:guild_id/bans")
	EndPointGuildMemberBan       = makeEndPoint("/guilds/:guild_id/bans/:user_id")
	EndPointGuildRoles           = makeEndPoint("/guilds/:guild_id/roles")
	EndPointGuildRole            = makeEndPoint("/guilds/:guild_id/roles/:role_id")
	EndPointGuildPrune           = makeEndPoint("/guilds/:guild_id/prune")
	EndPointGuildRegions         = makeEndPoint("/guilds/:guild_id/regions")
	EndPointGuildInvites         = makeEndPoint("/guilds/:guild_id/invites")
	EndPointGuildIntegrations    = makeEndPoint("/guilds/:guild_id/integrations")
	EndPointGuildIntegration     = makeEndPoint("/guilds/:guild_id/integrations/:integration_id")
	EndPointGuildIntegrationSync = makeEndPoint("/guilds/:guild_id/integrations/:integration_id/sync")
	EndPointGuildEmbed           = makeEndPoint("/guilds/:guild_id/embed")

	EndPointOwnUser    = makeEndPoint("/users/@me")
	EndPointUser       = makeEndPoint("/users/:user_id")
	EndPointOwnGuilds  = makeEndPoint("/users/@me/guilds")
	EndPointOwnGuild   = makeEndPoint("/users/@me/guilds/:guild_id")
	EndPointDMChannels = makeEndPoint("/users/@me/channels")
)

func makeEndPoint(path string) func(ids ...Snowflake) EndPoint {
	parts := strings.Split(path, "/")

	capacity := strings.Count(path, ":")
	snowflakes := make([]Snowflake, 0, capacity)
	endPoint := BaseUrl
	endPointIDs := make([]interface{}, 0, capacity)
	bucketID := ""
	bucketIDs := make([]interface{}, 0, capacity)

	for _, part := range parts {
		if utf8.RuneCountInString(part) != 0 {
			// For every variable in the path
			if part[:1] == ":" {
				snowflakes = append(snowflakes, Snowflake(42))

				// Add a format
				endPoint += "/%s"
				endPointIDs = append(endPointIDs, &snowflakes[len(snowflakes)-1])

				// If we're dealing with one of the "major" IDs, we can consider this part of the bucket
				switch part[1:] {
				case "guild_id":
					fallthrough
				case "channel_id":
					bucketID += "/%s"
					bucketIDs = append(bucketIDs, &snowflakes[len(snowflakes)-1])
				default:
					// Otherwise, generalize it with a simple zero, an ID that *should* never occur
					bucketID += "/0"
				}
			} else {
				// And well, if we're not dealing with a variable, just append it
				extension := "/" + part
				endPoint += extension
				bucketID += extension
			}
		}
	}

	// The idea is that any part of the library can create an EndPoint with their snowflakes
	return func(ids ...Snowflake) EndPoint {
		for i, snowflake := range ids {
			snowflakes[i] = snowflake
		}

		// But only the rest api functions themselves should decide whether they need the bucket ID or the url
		return EndPoint{fmt.Sprintf(endPoint, endPointIDs...), fmt.Sprintf(bucketID, bucketIDs...)}
	}
}

func (s *Session) doHttpGet(endPoint EndPoint, target interface{}) (err error) {
	err = s.rateLimit(endPoint, func() (*http.Response, error) {
		return s.doRequest("GET", endPoint.URL, nil, target)
	})
	return
}

func (s *Session) doHttpDelete(endPoint EndPoint, target interface{}) (err error) {
	err = s.rateLimit(endPoint, func() (*http.Response, error) {
		return s.doRequest("DELETE", endPoint.URL, nil, target)
	})
	return
}

func (s *Session) doHttpPost(endPoint EndPoint, body, target interface{}) (err error) {
	jsonBody, err := json.Marshal(body)

	if err == nil {
		byteBuf := bytes.NewReader(jsonBody)
		err = s.rateLimit(endPoint, func() (*http.Response, error) {
			return s.doRequest("POST", endPoint.URL, byteBuf, target)
		})
	}

	return
}

func (s *Session) doRequest(method, url string, body io.Reader, target interface{}) (response *http.Response, err error) {
	logger.Debugf("HTTP %s %s", method, strings.Replace(url, BaseUrl, "", 1))

	var (
		req    *http.Request
		client = http.Client{
			Timeout: 10 * time.Second,
		}
	)

	if req, err = http.NewRequest(method, url, body); err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.tokenType+" "+s.token)
	req.Header.Add("User-Agent", "DiscordBot (https://github.com/ikkerens/disgo, 1.0.0)")

	if response, err = client.Do(req); err != nil {
		return
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		fallthrough
	case 201:
		if target != nil {
			body := response.Body
			defer body.Close()
			if err = json.NewDecoder(body).Decode(target); err != nil {
				return
			}
		}
	case 204:
		fallthrough
	case 304:
		return
	default:
		var bodyBuf []byte
		bodyBuf, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return
		}
		return response, fmt.Errorf("Discord replied with status code %d: %s", response.StatusCode, string(bodyBuf))
	}

	return
}
