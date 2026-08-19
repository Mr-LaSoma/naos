package tokens

import "fmt"

type Position struct {
	FileName string
	Offset   int
	Line     int
	Column   int
}

func NewPosition(filename string) Position {
	return Position{
		FileName: filename,
		Offset:   0,
		Line:     1,
		Column:   1,
	}
}

func PositionGoForward(pos *Position, shouldWrap bool) {
	pos.Offset++
	pos.Column++
	if shouldWrap {
		pos.Line++
		pos.Column = 1
	}
}

func (pos Position) IsValid() bool { return pos.Line > 0 && pos.Column > 0 }

func (pos Position) String() string {
	return fmt.Sprintf(
		"%s:%d:%d",
		pos.FileName, pos.Line, pos.Column,
	)
}
