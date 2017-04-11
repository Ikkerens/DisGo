package disgo

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type EndPoint struct {
	URL    string
	Bucket string
}

var (
	BaseUrl = "https://discordapp.com/api"

	EndPointBotGateway = makeEndPoint("/gateway/bot")

	EndPointMessages = makeEndPoint("/channels/:channel_id/messages")
	EndPointMessage  = makeEndPoint("/channels/:channel_id/messages/:message_id")
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
