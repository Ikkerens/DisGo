package disgo

var (
	BaseUrl = "https://discordapp.com/api"

	EndPointGateway    = BaseUrl + "/gateway"
	EndPointBotGateway = EndPointGateway + "/bot"

	EndPointChannels = BaseUrl + "/channels"
	EndPointMessages = func(cID Snowflake) string { return EndPointChannels + "/" + cID.String() + "/messages" }
	EndPointMessage  = func(cID, mID Snowflake) string { return EndPointMessages(cID) + "/" + mID.String() }
)
