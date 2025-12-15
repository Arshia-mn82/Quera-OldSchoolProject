package protocol

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrEmptyLine     = errors.New("empty message")
	ErrMessageTooBig = errors.New("message too big")
)

const MaxLineBytes = 1 << 20

type Request struct {
	Method string          `json:"method,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}
type Response struct {
	Status  bool        `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func readLineLimited(r *bufio.Reader, max int) ([]byte, error) {
	var buf []byte

	for {
		part, isPrefix, err := r.ReadLine()
		if err != nil {
			return nil, err
		}

		buf = append(buf, part...)
		if len(buf) > max {
			return nil, ErrMessageTooBig
		}

		if !isPrefix {
			break
		}
	}

	if n := len(buf); n > 0 && buf[n-1] == '\r' {
		buf = buf[:n-1]
	}

	return buf, nil
}

func ReadRequest(r *bufio.Reader) (*Request, error) {
	line, err := readLineLimited(r, MaxLineBytes)
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, ErrEmptyLine
	}

	var req Request
	if err := json.Unmarshal(line, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func WriteResponse(w *bufio.Writer, resp Response) error {
	b, err := json.Marshal(resp)
	if err != nil {
		return nil
	}
	if _, err := w.Write(b); err != nil {
		return err
	}

	if err := w.WriteByte('\n'); err != nil {
		return err
	}
	return w.Flush()
}

func IsEOF(err error) bool {
	return errors.Is(err, io.EOF)
}
