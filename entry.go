package tui

import (
	"image"
	"strings"

	wordwrap "github.com/mitchellh/go-wordwrap"
)

var _ Widget = &Entry{}

// Entry is a one-line text editor. It lets the user supply your application
// with text, e.g. to input user and password information.
type Entry struct {
	text string

	size image.Point

	focused bool

	onTextChange func(*Entry)
	onSubmit     func(*Entry)

	sizePolicyX SizePolicy
	sizePolicyY SizePolicy
}

// NewEntry returns a new Entry.
func NewEntry() *Entry {
	return &Entry{}
}

// Draw draws the entry.
func (e *Entry) Draw(p *Painter) {
	s := e.Size()

	style := "entry"
	if e.focused {
		style += ".focused"
	}

	p.WithStyledBrush(style, func(p *Painter) {
		tw := stringWidth(e.text)

		offx := tw - s.X
		if e.focused {
			offx++
		}

		text := e.text
		if tw >= s.X {
			text = text[offx:]
		}

		p.FillRect(0, 0, s.X, 1)
		p.DrawText(0, 0, text)

		if e.focused {
			p.DrawCursor(stringWidth(text), 0)
		}
	})
}

// Size returns the size of the entry.
func (e *Entry) Size() image.Point {
	return e.size
}

// MinSizeHint returns the minimum size the widget is allowed to be.
func (e *Entry) MinSizeHint() image.Point {
	return image.Point{1, 1}
}

// SizeHint returns the recommended size for the entry.
func (e *Entry) SizeHint() image.Point {
	return image.Point{10, 1}
}

// SizePolicy returns the default layout behavior.
func (e *Entry) SizePolicy() (SizePolicy, SizePolicy) {
	return e.sizePolicyX, e.sizePolicyY
}

// Resize updates the size of the entry.
func (e *Entry) Resize(size image.Point) {
	e.size = size
}

// OnEvent handles terminal events.
func (e *Entry) OnEvent(ev Event) {
	if !e.focused {
		return
	}

	if ev.Type != EventKey {
		return
	}

	if ev.Key != 0 {
		switch ev.Key {
		case KeyEnter:
			if e.onSubmit != nil {
				e.onSubmit(e)
			}
			return
		case KeySpace:
			e.text = e.text + string(' ')
			if e.onTextChange != nil {
				e.onTextChange(e)
			}
			return
		case KeyBackspace2:
			if len(e.text) > 0 {
				e.text = trimRightLen(e.text, 1)
				if e.onTextChange != nil {
					e.onTextChange(e)
				}
			}
			return
		}
	} else {
		e.text = e.text + string(ev.Ch)
		if e.onTextChange != nil {
			e.onTextChange(e)
		}
	}
}

// OnChanged sets a function to be run whenever the content of the entry has
// been changed.
func (e *Entry) OnChanged(fn func(entry *Entry)) {
	e.onTextChange = fn
}

// OnSubmit sets a function to be run whenever the user submits the entry (by
// pressing KeyEnter).
func (e *Entry) OnSubmit(fn func(entry *Entry)) {
	e.onSubmit = fn
}

// SetText sets the text content of the entry.
func (e *Entry) SetText(text string) {
	e.text = text
}

// Text returns the text content of the entry.
func (e *Entry) Text() string {
	return e.text
}

// SetSizePolicy sets the size policy for each axis.
func (e *Entry) SetSizePolicy(horizontal, vertical SizePolicy) {
	e.sizePolicyX = horizontal
	e.sizePolicyY = vertical
}

// SetFocused focuses this entry.
func (e *Entry) SetFocused(f bool) {
	e.focused = f
}

func (e *Entry) heightForWidth(w int) int {
	return len(strings.Split(wordwrap.WrapString(e.text, uint(w)), "\n"))
}
