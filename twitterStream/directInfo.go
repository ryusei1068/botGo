package twitterstream

type Association struct {
	word string
	url  string
}

type Keys struct {
	streamId  string
	webhookId string
	tag       string
}

// WebHookId, Key word, url
var DirectInfo = make(map[Keys]Association)

const (
	webhook = "https://discord.com/api/webhooks/"
)
