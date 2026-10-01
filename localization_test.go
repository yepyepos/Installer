/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// nonEmptyFields returns the names of all string fields that are blank.
func nonEmptyFields(v *UI) []string {
	var blank []string
	rv := reflect.ValueOf(*v)
	for i := 0; i < rv.NumField(); i++ {
		if rv.Field(i).Kind() != reflect.String {
			continue
		}
		if strings.TrimSpace(rv.Field(i).String()) == "" {
			blank = append(blank, rv.Type().Field(i).Name)
		}
	}
	return blank
}

func TestEnglishFallbackComplete(t *testing.T) {
	if blank := nonEmptyFields(&English); len(blank) > 0 {
		t.Errorf("English fallback has blank strings: %v", blank)
	}
}

func TestChineseComplete(t *testing.T) {
	if blank := nonEmptyFields(&Chinese); len(blank) > 0 {
		t.Errorf("Chinese locale has blank strings: %v", blank)
	}
}

func TestZhAndEnLocalesDiffer(t *testing.T) {
	// At least the vast majority of strings must actually be translated,
	// not identical copies of the English text.
	rvEn, rvZh := reflect.ValueOf(English), reflect.ValueOf(Chinese)
	same, total := 0, 0
	for i := 0; i < rvEn.NumField(); i++ {
		total++
		if rvEn.Field(i).String() == rvZh.Field(i).String() {
			same++
		}
	}
	if same == total {
		t.Errorf("Chinese locale is identical to English (%d/%d fields)", same, total)
	}
	t.Logf("%d/%d fields translated", total-same, total)
}

func TestZhLocaleContainsCJK(t *testing.T) {
	samples := []string{L.BtnInstall, L.BtnUninstall, L.SelectInstallPrompt, L.CliExitSuccess}
	for _, s := range samples {
		if !strings.ContainsFunc(s, func(r rune) bool { return r >= 0x4E00 && r <= 0x9FFF }) {
			t.Errorf("expected CJK characters in %q", s)
		}
	}
}

func TestPickLanguage(t *testing.T) {
	t.Setenv("VENCORD_INSTALLER_LANG", "en")
	if pickLanguage() != &English {
		t.Error("VENCORD_INSTALLER_LANG=en must select English")
	}

	t.Setenv("VENCORD_INSTALLER_LANG", "zh-CN")
	if pickLanguage() != &Chinese {
		t.Error("zh-CN (and anything unrecognised) must select Chinese")
	}

	t.Setenv("VENCORD_INSTALLER_LANG", "")
	if pickLanguage() != &Chinese {
		t.Error("default language must be Chinese")
	}
}

func TestFormatPlaceholdersValid(t *testing.T) {
	// Strings with % verbs must be usable with fmt.Sprintf: verify with
	// representative args that no "%!x(MISSING)" artifacts appear.
	cases := []string{
		fmt.Sprintf(L.CliBranchNotFoundFmt, "ptb"),
		fmt.Sprintf(L.LogPatchingFmt, "C:\\Discord"),
		fmt.Sprintf(L.ErrUnpatchPatchedFmt, "C:\\Discord", "boom"),
		fmt.Sprintf(L.ErrFlatpakFmt, "/x", "boom"),
		fmt.Sprintf(L.ErrNoAsarFmt, "/x"),
		fmt.Sprintf(L.ErrOpenAsarFetchFmt, "403", "403 Forbidden"),
		fmt.Sprintf(L.LogFetchFallbackFmt, "u", 429, "f"),
	}
	for _, out := range cases {
		if strings.Contains(out, "%!") {
			t.Errorf("format produced error artifact: %q", out)
		}
	}
}
