package cui

// MsgBar is a struct to represent a single row horizontal
// bar on the screen.
type MsgBar struct {
	row       int
	title     string
	message   string
	showTitle bool
}

// SetTitle sets the title for the MsgBar in question
func (m *MsgBar) SetTitle(s string) {
	m.title = s
}

// SetMessage sets the message for the MsgBar in question
func (m *MsgBar) SetMessage(s string) {
	m.message = s
}

// ClearAll clears all text from the message bar (title and message)
func (m MsgBar) ClearAll() {
	MoveCursorTo(m.row, 1)
	Clear("line")
}

// ClearMessage clears all message text while leaving the title in place
func (m *MsgBar) ClearMessage() {
	MoveCursorTo(m.row, len(m.title)+1)
	Clear("right")
}
