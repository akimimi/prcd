package main

import (
	"testing"
	"time"
)

func resetDedupCache() {
	dedupCacheMu.Lock()
	defer dedupCacheMu.Unlock()
	dedupCache = make(map[string]time.Time)
}

func TestIsDuplicateMessage_Disabled(t *testing.T) {
	resetDedupCache()
	settings.dedupWindowSeconds = 0
	payload := []byte(`{"hook_name":"merge_request_hooks","action":"merge"}`)
	if isDuplicateMessage(payload) {
		t.Fatalf("first call should not be duplicate")
	}
	if isDuplicateMessage(payload) {
		t.Fatalf("dedup disabled (window=0) should never report duplicate")
	}
}

func TestIsDuplicateMessage_DetectsDuplicateInWindow(t *testing.T) {
	resetDedupCache()
	settings.dedupWindowSeconds = 10
	payload := []byte(`{"hook_name":"merge_request_hooks","action":"merge","id":1}`)

	if isDuplicateMessage(payload) {
		t.Fatalf("first call should not be duplicate")
	}
	if !isDuplicateMessage(payload) {
		t.Fatalf("identical payload within window should be duplicate")
	}
}

func TestIsDuplicateMessage_DifferentPayloadsNotDuplicate(t *testing.T) {
	resetDedupCache()
	settings.dedupWindowSeconds = 10
	a := []byte(`{"id":1}`)
	b := []byte(`{"id":2}`)

	if isDuplicateMessage(a) {
		t.Fatalf("first call for a should not be duplicate")
	}
	if isDuplicateMessage(b) {
		t.Fatalf("different payload b should not be duplicate")
	}
}

func TestIsDuplicateMessage_ExpiresAfterWindow(t *testing.T) {
	resetDedupCache()
	// 1 秒窗口便于测试
	settings.dedupWindowSeconds = 1
	payload := []byte(`{"id":"expire-test"}`)

	if isDuplicateMessage(payload) {
		t.Fatalf("first call should not be duplicate")
	}
	// 等待超过窗口
	time.Sleep(1100 * time.Millisecond)
	if isDuplicateMessage(payload) {
		t.Fatalf("payload after window expiry should not be duplicate")
	}
}
