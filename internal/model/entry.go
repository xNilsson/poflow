package model

// MsgEntry represents a single translation entry in a .po file
type MsgEntry struct {
	MsgID      string   `json:"msgid"`
	MsgStr     string   `json:"msgstr"`
	Comments   []string `json:"comments,omitempty"`
	References []string `json:"references,omitempty"`
}

// IsEmpty returns true if the translation (msgstr) is empty
func (e *MsgEntry) IsEmpty() bool {
	return e.MsgStr == ""
}
