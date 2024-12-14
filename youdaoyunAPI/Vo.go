package youdaoyunAPI

type RespTran struct {
	TSpeakUrl     string        `json:"tSpeakUrl"`
	RequestId     string        `json:"requestId"`
	Query         string        `json:"query"`
	Translation   []string      `json:"translation"`
	MTerminalDict MTerminalDict `json:"mTerminalDict"`
	ErrorCode     string        `json:"errorCode"`
	Dict          Dict          `json:"dict"`
	WebDict       WebDict       `json:"webdict"`
	L             string        `json:"l"`
	IsWord        bool          `json:"isWord"`
	SpeakUrl      string        `json:"speakUrl"`
}

type MTerminalDict struct {
	URL string `json:"url"`
}

type Dict struct {
	URL string `json:"url"`
}

type WebDict struct {
	URL string `json:"url"`
}
