package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Response struct {
	Status  bool            `json:"status,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type Request struct {
	Method string      `json:"method,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

// Used to read IDs from create endpoints where router returns models (School/Person/Class)
// and protocol marshals them to JSON.
type hasID struct {
	ID uint `json:"ID"`
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "server address host:port")
	pause := flag.Duration("pause", 0, "pause between requests (e.g. 200ms)")
	flag.Parse()

	fmt.Println("== OldSchool Scenario Client ==")
	fmt.Println("Connecting to:", *addr)

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Println("Dial error:", err)
		os.Exit(1)
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	// ---- helpers ----
	send := func(method string, payload any) (Response, error) {
		if *pause > 0 {
			time.Sleep(*pause)
		}

		req := Request{Method: method, Data: payload}
		b, err := json.Marshal(req)
		if err != nil {
			return Response{}, err
		}

		// newline-delimited JSON
		if _, err := w.Write(b); err != nil {
			return Response{}, err
		}
		if err := w.WriteByte('\n'); err != nil {
			return Response{}, err
		}
		if err := w.Flush(); err != nil {
			return Response{}, err
		}

		line, err := r.ReadBytes('\n')
		if err != nil {
			return Response{}, err
		}
		line = bytes.TrimSpace(line)

		var resp Response
		if err := json.Unmarshal(line, &resp); err != nil {
			return Response{}, fmt.Errorf("invalid response JSON: %w | raw=%s", err, string(line))
		}

		// Pretty print I/O
		fmt.Println("\n--- REQUEST ---")
		fmt.Println(prettyJSON(b))
		fmt.Println("--- RESPONSE ---")
		fmt.Println(prettyJSON(line))

		return resp, nil
	}

	mustOK := func(step string, resp Response) {
		if !resp.Status {
			fmt.Printf("❌ %s failed: %s\n", step, resp.Message)
			os.Exit(1)
		}
		fmt.Printf("✅ %s OK\n", step)
	}

	mustFail := func(step string, resp Response, msgContains string) {
		if resp.Status {
			fmt.Printf("❌ %s expected failure but got success\n", step)
			os.Exit(1)
		}
		if msgContains != "" && !strings.Contains(strings.ToLower(resp.Message), strings.ToLower(msgContains)) {
			fmt.Printf("❌ %s expected message containing %q but got %q\n", step, msgContains, resp.Message)
			os.Exit(1)
		}
		fmt.Printf("✅ %s failed as expected (%s)\n", step, resp.Message)
	}

	extractID := func(resp Response) uint {
		var x hasID
		// For create endpoints, resp.Data is JSON of your model struct, which includes "ID"
		if err := json.Unmarshal(resp.Data, &x); err != nil {
			fmt.Println("Cannot parse ID from Data:", err, "raw:", string(resp.Data))
			os.Exit(1)
		}
		if x.ID == 0 {
			fmt.Println("Parsed ID=0, raw:", string(resp.Data))
			os.Exit(1)
		}
		return x.ID
	}

	// ---- scenario state ----
	var (
		s1ID, s2ID             uint
		t1ID, t2ID             uint
		student1ID, student2ID uint
		c1ID, c2ID, c3ID        uint
	)

	// =========================
	// SCENARIO STARTS HERE
	// =========================

	// 1) Create schools
	resp, err := send("/school/create", map[string]any{"name": "S1"})
	mustNoErr("Create school S1", err)
	mustOK("Create school S1", resp)
	s1ID = extractID(resp)

	resp, err = send("/school/create", map[string]any{"name": "S2"})
	mustNoErr("Create school S2", err)
	mustOK("Create school S2", resp)
	s2ID = extractID(resp)

	// 1b) Duplicate school (should fail)
	resp, err = send("/school/create", map[string]any{"name": "S1"})
	mustNoErr("Duplicate school S1", err)
	mustFail("Duplicate school S1", resp, "exists")

	// 2) Create people
	resp, err = send("/person/create", map[string]any{"name": "T1", "role": "teacher"})
	mustNoErr("Create teacher T1", err)
	mustOK("Create teacher T1", resp)
	t1ID = extractID(resp)

	resp, err = send("/person/create", map[string]any{"name": "T2", "role": "teacher"})
	mustNoErr("Create teacher T2", err)
	mustOK("Create teacher T2", resp)
	t2ID = extractID(resp)

	resp, err = send("/person/create", map[string]any{"name": "Stu1", "role": "student"})
	mustNoErr("Create student Stu1", err)
	mustOK("Create student Stu1", resp)
	student1ID = extractID(resp)

	resp, err = send("/person/create", map[string]any{"name": "Stu2", "role": "student"})
	mustNoErr("Create student Stu2", err)
	mustOK("Create student Stu2", resp)
	student2ID = extractID(resp)

	// 2b) Invalid role (should fail)
	resp, err = send("/person/create", map[string]any{"name": "BadRole", "role": "admin"})
	mustNoErr("Create invalid role", err)
	mustFail("Create invalid role", resp, "invalid")

	// 3) Create classes
	resp, err = send("/class/create", map[string]any{
		"name":       "C1",
		"school_id":  s1ID,
		"teacher_id": t1ID,
	})
	mustNoErr("Create class C1", err)
	mustOK("Create class C1", resp)
	c1ID = extractID(resp)

	resp, err = send("/class/create", map[string]any{
		"name":       "C2",
		"school_id":  s1ID,
		"teacher_id": t1ID,
	})
	mustNoErr("Create class C2", err)
	mustOK("Create class C2", resp)
	c2ID = extractID(resp)

	resp, err = send("/class/create", map[string]any{
		"name":       "C3",
		"school_id":  s2ID,
		"teacher_id": t2ID,
	})
	mustNoErr("Create class C3", err)
	mustOK("Create class C3", resp)
	c3ID = extractID(resp)

	// 3b) Create class with student as teacher (should fail)
	resp, err = send("/class/create", map[string]any{
		"name":       "BadClass",
		"school_id":  s1ID,
		"teacher_id": student1ID,
	})
	mustNoErr("Create class with student teacher", err)
	mustFail("Create class with student teacher", resp, "role")

	// 4) Enroll student1 into C1 (S1) — OK
	resp, err = send("/class/add/student", map[string]any{"student_id": student1ID, "class_id": c1ID})
	mustNoErr("Enroll Stu1 -> C1", err)
	mustOK("Enroll Stu1 -> C1", resp)

	// 4b) Duplicate enrollment — should fail
	resp, err = send("/class/add/student", map[string]any{"student_id": student1ID, "class_id": c1ID})
	mustNoErr("Duplicate enroll Stu1 -> C1", err)
	mustFail("Duplicate enroll Stu1 -> C1", resp, "duplicate")

	// 4c) Student1 enroll in another class in SAME school (C2 in S1) — OK
	resp, err = send("/class/add/student", map[string]any{"student_id": student1ID, "class_id": c2ID})
	mustNoErr("Enroll Stu1 -> C2", err)
	mustOK("Enroll Stu1 -> C2", resp)

	// 4d) Student1 enroll in a class in DIFFERENT school (C3 in S2) — should fail
	resp, err = send("/class/add/student", map[string]any{"student_id": student1ID, "class_id": c3ID})
	mustNoErr("Enroll Stu1 -> C3 (different school)", err)
	mustFail("Enroll Stu1 -> C3 (different school)", resp, "different school")

	// 4e) Student2 enroll in S2 (C3) — OK (first enrollment sets their school)
	resp, err = send("/class/add/student", map[string]any{"student_id": student2ID, "class_id": c3ID})
	mustNoErr("Enroll Stu2 -> C3", err)
	mustOK("Enroll Stu2 -> C3", resp)

	// 5) WhoAmI for teacher T1 (should list classes they teach: C1, C2)
	resp, err = send("/who/am/i", map[string]any{"id": t1ID})
	mustNoErr("WhoAmI teacher T1", err)
	mustOK("WhoAmI teacher T1", resp)

	// 5b) WhoAmI for student1 (should list class IDs enrolled: C1, C2)
	resp, err = send("/who/am/i", map[string]any{"id": student1ID})
	mustNoErr("WhoAmI student Stu1", err)
	mustOK("WhoAmI student Stu1", resp)

	// 5c) WhoAmI with unknown id (should fail)
	resp, err = send("/who/am/i", map[string]any{"id": uint(999999)})
	mustNoErr("WhoAmI unknown", err)
	mustFail("WhoAmI unknown", resp, "not found")

	// 6) Unknown method (should fail)
	resp, err = send("/unknown/method", map[string]any{"x": 1})
	mustNoErr("Unknown method", err)
	mustFail("Unknown method", resp, "unknown")

	// =========================
	// SCENARIO END
	// =========================

	fmt.Println("\n=========================")
	fmt.Println("✅ Scenario completed OK.")
	fmt.Println("=========================")
	fmt.Printf("Created IDs:\n")
	fmt.Printf("  Schools: S1=%d S2=%d\n", s1ID, s2ID)
	fmt.Printf("  Teachers: T1=%d T2=%d\n", t1ID, t2ID)
	fmt.Printf("  Students: Stu1=%d Stu2=%d\n", student1ID, student2ID)
	fmt.Printf("  Classes: C1=%d C2=%d C3=%d\n", c1ID, c2ID, c3ID)
	fmt.Println("\nNow open your server DB file in Beekeeper Studio and check tables:")
	fmt.Println("  schools, people, classes, enrollments")
}

func mustNoErr(step string, err error) {
	if err != nil {
		fmt.Printf("❌ %s error: %v\n", step, err)
		os.Exit(1)
	}
}

func prettyJSON(b []byte) string {
	var out bytes.Buffer
	if err := json.Indent(&out, b, "", "  "); err != nil {
		return string(b)
	}
	return out.String()
}
