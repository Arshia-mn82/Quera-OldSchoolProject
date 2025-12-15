package protocol_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"

	"OldSchool/internal/transport/protocol"
)

func TestWriteResponse_WritesNDJSON(t *testing.T) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	resp := protocol.Response{
		Status:  true,
		Message: "ok",
		Data: map[string]any{
			"id":   1,
			"name": "S1",
		},
	}

	if err := protocol.WriteResponse(w, resp); err != nil {
		t.Fatalf("WriteResponse error: %v", err)
	}

	out := buf.String()
	if len(out) == 0 || out[len(out)-1] != '\n' {
		t.Fatalf("expected newline-delimited JSON, got: %q", out)
	}

	// remove trailing '\n' and validate JSON
	line := out[:len(out)-1]

	var decoded protocol.Response
	if err := json.Unmarshal([]byte(line), &decoded); err != nil {
		t.Fatalf("response is not valid JSON: %v, raw=%q", err, line)
	}

	if decoded.Status != true || decoded.Message != "ok" {
		t.Fatalf("unexpected decoded response: %+v", decoded)
	}
}

func TestReadRequest_ReadsNDJSON(t *testing.T) {
	in := `{"method":"/school/create","data":{"name":"MIT"}}` + "\n"
	r := bufio.NewReader(bytes.NewBufferString(in))

	req, err := protocol.ReadRequest(r)
	if err != nil {
		t.Fatalf("ReadRequest error: %v", err)
	}

	if req.Method != "/school/create" {
		t.Fatalf("expected method /school/create, got %q", req.Method)
	}

	// Data is RawMessage; verify it’s valid JSON
	var payload map[string]any
	if err := json.Unmarshal(req.Data, &payload); err != nil {
		t.Fatalf("data invalid JSON: %v", err)
	}
	if payload["name"] != "MIT" {
		t.Fatalf("expected name MIT, got %v", payload["name"])
	}
}
