package twitterstream

type Association struct {
	word     string
	url      string
	streamId string
}

// WebHookId, Key word, url
var DirectInfo = make(map[string]Association)

const (
	webhook = "https://discord.com/api/webhooks/"
)
