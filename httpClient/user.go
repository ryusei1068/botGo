package httpclient

type User struct {
	Id            string `json:"id"`
	UserName      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
	System        bool   `json:"system,omitempty"`
	MfaEnabled    bool   `json:"mfa_enabled,omitempty"`
	Banner        string `json:"banner,omitempty"`
	AccentColor   int64  `json:"accent_color,omitempty"`
	Locale        string `json:"locale,omitempty"`
	Varified      bool   `json:"varified,omitempty"`
	Email         string `json:"email,omitempty"`
	Flags         int64  `json:"flags,omitempty"`
	PreminumType  int64  `json:"preminum_type,omitempty"`
	PublicFlags   int64  `json:"public_flags"`
}
