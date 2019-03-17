package cui

import (

)


type MsgBar struct {
	row						int
	title					string
	message				string
	showTitle			bool
}

func (m *MsgBar) SetTitle(s string) { 
	m.title = s 
}

func (m *MsgBar) SetMessage(s string) { 
	m.message = s 
}

func (m MsgBar) ClearAll() { 
	MoveCursorTo(m.row, 1)
	Clear("line")
}

func (m *MsgBar) ClearMessage() {
	MoveCursorTo(m.row, len(m.title) + 1)
	Clear("right")
}
