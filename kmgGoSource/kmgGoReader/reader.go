package kmgGoReader

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

type Reader struct {
	buf     []byte //需要读入的数据
	pos     int    //当前位置
	filePos *FilePos
}

func NewReader(buf []byte, filePos *FilePos) *Reader {
	return &Reader{
		buf:     buf,
		filePos: filePos,
	}
}

func NewReaderWithPosFile(filename string, content []byte) *Reader {
	pos := NewPosFile(filename, content)
	return NewReader(content, pos)
}

func (r *Reader) Pos() int {
	return r.pos
}
func (r *Reader) BufToCurrent(start int) []byte {
	return r.buf[start:r.pos]
}

func (r *Reader) IsEof() bool {
	return r.pos >= len(r.buf)
}
func (r *Reader) ReadByte() byte {
	//if r.IsEof() {
	//	panic(r.GetFileLineInfo() + " unexcept EOF")
	//}
	out := r.buf[r.pos]
	r.pos++
	return out
}

func (r *Reader) NextByte() byte {
	return r.buf[r.pos]
}

func (r *Reader) IsMatchAfter(s []byte) bool {
	return len(r.buf)-r.pos >= len(s) && r.buf[r.pos] == s[0] && bytes.Equal(r.buf[r.pos:r.pos+len(s)], s)
}

// 读取到某个字符,或者读取到结束(该字符会已经被读过)
func (r *Reader) ReadUntilByte(b byte) []byte {
	startPos := r.pos
	for {
		if r.IsEof() {
			return r.buf[startPos:]
		}
		if r.ReadByte() == b {
			return r.buf[startPos:r.pos]
		}
	}
}

// 回调返回真的时候,停止读取,(这个回调提到的字符串也包含在内)
func (r *Reader) ReadUntilRuneCb(cb func(run rune) bool) []byte {
	startPos := r.pos
	for {
		if r.IsEof() {
			return r.buf[startPos:]
		}
		run, size := utf8.DecodeRune(r.buf[r.pos:])
		r.pos += size
		if cb(run) {
			return r.buf[startPos:r.pos]
		}
	}
}

// 读取到某个字符串,或者读取到结束(该字符串会已经被读过)
func (r *Reader) ReadUntilString(s []byte) []byte {
	startPos := r.pos
	for {
		if r.IsEof() {
			return r.buf[startPos:]
		}
		if r.IsMatchAfter(s) {
			r.pos += len(s)
			return r.buf[startPos:r.pos]
		}
		r.pos++
	}
}

func (r *Reader) ReadAllSpace() {
	for {
		if r.IsEof() {
			return
		}
		run, size := utf8.DecodeRune(r.buf[r.pos:])
		if !unicode.IsSpace(run) {
			return
		}
		r.pos += size
	}
}

func (r *Reader) ReadAllSpaceWithoutLineBreak() {
	for {
		if r.IsEof() {
			return
		}
		run, size := utf8.DecodeRune(r.buf[r.pos:])
		if unicode.IsSpace(run) && run != '\n' {
			r.pos += size
		} else {
			return
		}
	}
}

func (r *Reader) ReadRune() rune {
	run, size := utf8.DecodeRune(r.buf[r.pos:])
	r.pos += size
	return run
}

func (r *Reader) UnreadRune() rune {
	run, size := utf8.DecodeLastRune(r.buf[:r.pos])
	if size == 0 {
		panic(r.GetFileLineInfo() + " [UnreadRune] last is not valid utf8 code.")
	}
	r.pos -= size
	return run
}

func (r *Reader) UnreadByte() {
	r.pos -= 1
}

func (r *Reader) MustReadMatch(s []byte) {
	if !r.IsMatchAfter(s) {
		panic(r.GetFileLineInfo() + " [MustReadMatch] not match " + string(s))
	}
	r.pos += len(s)
}

func (r *Reader) MustReadWithSize(size int) []byte {
	if r.IsEof() {
		panic(r.GetFileLineInfo() + " unexpect EOF")
	}
	output := r.buf[r.pos : r.pos+size]
	r.pos += size
	return output
}

func (r *Reader) GetFileLineInfo() string {
	return r.filePos.PosString(r.pos)
}
