package httpclient

type Webhook struct {
	Type          int    `josn:"type"`
	Id            string `json:"id"`
	Name          string `json:"name"`
	Avatar        string `json:"avatar"`
	ChannelId     string `json:"channleid"`
	GuildId       string `json:"guildId"`
	ApplicationId string `json:"applicationid"`
	Token         string `josn:"token"`
	User          User   `josn:"user"`
}
