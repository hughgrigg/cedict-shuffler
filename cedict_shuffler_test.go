package main

import (
	"testing"
)

func TestMakeCedictEntry(t *testing.T) {
	entry := makeCedictEntry("這類 这类 [zhe4 lei4] /this kind (of)/")
	if entry.traditional != "這類" {
		t.Error("Expected 這類 but got ", entry.traditional)
	}
	if entry.simplified != "这类" {
		t.Error("Expected 这类 but got ", entry.simplified)
	}
	if entry.pinyin != "zhe4 lei4" {
		t.Error("Expected 這類 but got ", entry.pinyin)
	}
	if entry.definition != "this kind (of)" {
		t.Error("Expected Asia-Pacific but got ", entry.definition)
	}
}

func TestVsToUmlaut(t *testing.T) {
	replaced := vsToUmlaut("lv: LV:")
	if replaced != "lu\u0308 LU\u0308" {
		t.Error("Expected lu\u0308 LU\u0308 but got ", replaced)
	}
}

func TestToneMark(t *testing.T) {
	numberedToMarked := map[string]string{
		"a1":    "ā",
		"qiao1": "qiāo",
		"ba2":   "bá",
		"tuo2":  "tuó",
		"wa3":   "wǎ",
		"guo3":  "guǒ",
		"ku4":   "kù",
		"cuo4":  "cuò",
		"le5":   "le",
		"kang5": "kang",
	}
	mark := toneMarker()
	for numbered, marked := range numberedToMarked {
		tone, letters := toneAndLetters(numbered)
		testMarked := mark(tone, letters)
		if testMarked != marked {
			t.Error("Expected "+marked+" but got ", testMarked)
		}
	}
}

func TestToneColour(t *testing.T) {
	colourer := toneColourer()
	test := colourer(3, "test")
	if test != "\x1b[32mtest\x1b[0m" {
		t.Error("Expected \x1b[32mtest\x1b[0m but got ", test)
	}
}

func TestToneAndLetters(t *testing.T) {
	tone, letters := toneAndLetters("gong1")
	if tone != 1 {
		t.Error("Expected 1 but got ", tone)
	}
	if letters != "gong" {
		t.Error("Expected gong but got ", letters)
	}
}
