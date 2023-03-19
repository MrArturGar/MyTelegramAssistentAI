package tests

import (
	"MyTelegramAssistentAI/src/layers/messageLayer"
	"strings"
	"testing"
)

func SpecCommandReplacer(t *testing.T) {
	requests := []string{
		"Hello world",
		"Hello world\n\n/dall-e picture\n",
		"Hello /dall-e picture\n world",
		"/dall-e picture!",
		"/dall-e picture\nHello world!",
		"Hello /dall-e picture. world",
		"Hello world /dall-e picture.",
	}

	for i, s := range requests {
		s1, s2 := messageLayer.SpecCommandReplacer(&s)
		if strings.Contains(*s1, "picture") || strings.Contains(*s2, "Hello") || strings.Contains(*s1, "Hello") {
			t.Errorf("Fail data, %s", i)
		}
	}

}
