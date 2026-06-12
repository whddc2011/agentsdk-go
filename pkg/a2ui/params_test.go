package a2ui

import "testing"

func TestExtractMessagesFromToolParams_messagesArray(t *testing.T) {
	params := map[string]any{
		"messages": []any{
			map[string]any{"createSurface": map[string]any{"surfaceId": "main", "catalogId": "basic"}},
		},
	}
	msgs, err := ExtractMessagesFromToolParams(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 || msgs[0].Kind() != "createSurface" {
		t.Fatalf("got %d msgs kind=%q", len(msgs), msgs[0].Kind())
	}
	if msgs[1].Kind() != "updateComponents" {
		t.Fatalf("expected synthesized root update, got %q", msgs[1].Kind())
	}
}

func TestExtractMessagesFromToolParams_topLevelAction(t *testing.T) {
	params := map[string]any{
		"createSurface": map[string]any{"surfaceId": "main", "catalogId": "basic"},
	}
	msgs, err := ExtractMessagesFromToolParams(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 {
		t.Fatalf("got %d messages", len(msgs))
	}
}

func TestExtractMessagesFromToolParams_jsonl(t *testing.T) {
	params := map[string]any{
		"jsonl": "{\"createSurface\":{\"surfaceId\":\"main\"}}\n",
	}
	msgs, err := ExtractMessagesFromToolParams(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 {
		t.Fatalf("got %d messages", len(msgs))
	}
}

func TestExtractMessagesFromToolParams_messageSingular(t *testing.T) {
	params := map[string]any{
		"message": map[string]any{
			"updateComponents": map[string]any{
				"surfaceId": "main",
				"components": []any{
					map[string]any{"id": "root", "component": map[string]any{"Column": map[string]any{}}},
				},
			},
		},
	}
	msgs, err := ExtractMessagesFromToolParams(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 || msgs[0].Kind() != "createSurface" || msgs[1].Kind() != "updateComponents" {
		t.Fatalf("kind=%q len=%d", msgs[0].Kind(), len(msgs))
	}
}

func TestExtractMessagesFromToolParams_a2uiAlias(t *testing.T) {
	params := map[string]any{
		"a2ui": []any{
			map[string]any{"deleteSurface": map[string]any{"surfaceId": "main"}},
		},
	}
	msgs, err := ExtractMessagesFromToolParams(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 || msgs[0].Kind() != "deleteSurface" {
		t.Fatalf("kind=%q", msgs[0].Kind())
	}
}

func TestExtractMessagesFromToolParams_empty(t *testing.T) {
	msgs, err := ExtractMessagesFromToolParams(map[string]any{"foo": "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 0 {
		t.Fatalf("expected empty, got %d", len(msgs))
	}
}

func TestIsLikelyA2UIMessage(t *testing.T) {
	if !IsLikelyA2UIMessage([]byte(`{"createSurface":{"surfaceId":"main"}}`)) {
		t.Fatal("expected true")
	}
	if IsLikelyA2UIMessage([]byte(`{"text":"hello"}`)) {
		t.Fatal("expected false")
	}
}
