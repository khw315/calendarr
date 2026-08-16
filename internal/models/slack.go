package models

type SlackText struct {
	Type string `json:"type"` // "mrkdwn" or "plain_text"
	Text string `json:"text"`
}

type SlackField struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SlackBlock struct {
	Type   string       `json:"type"`
	Text   *SlackText   `json:"text,omitempty"`
	Fields []SlackField `json:"fields,omitempty"`
}

type SlackAttachment struct {
	Color  string       `json:"color,omitempty"`
	Blocks []SlackBlock `json:"blocks,omitempty"`
}

type SlackPayload struct {
	Text        string            `json:"text,omitempty"`
	Blocks      []SlackBlock      `json:"blocks,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}
