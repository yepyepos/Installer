/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"os"
	path "path/filepath"
	"testing"
)

func touch(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(path.Dir(p), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("stub"), 0644); err != nil {
		t.Fatal(err)
	}
}

// An interrupted Discord update leaves an app-* folder with an empty
// resources directory behind. ParseDiscord must skip it and pick the newest
// folder that actually contains Discord asar files.
func TestParseDiscordSkipsEmptyUpdateStub(t *testing.T) {
	base := t.TempDir()

	touch(t, path.Join(base, "app-1.0.9259", "resources", "app.asar"))
	if err := os.MkdirAll(path.Join(base, "app-1.0.9260", "resources"), 0755); err != nil {
		t.Fatal(err)
	}

	di := ParseDiscord(base, "")
	if di == nil {
		t.Fatal("ParseDiscord returned nil, want the valid install")
	}
	if want := path.Join(base, "app-1.0.9259", "resources", "app"); di.appPath != want {
		t.Errorf("appPath = %q, want %q", di.appPath, want)
	}
	if di.isPatched {
		t.Error("install without _app.asar must not be reported as patched")
	}
}

// An asar-less stub must not win over a patched older install either.
func TestParseDiscordPrefersPatchedOverStub(t *testing.T) {
	base := t.TempDir()

	res := path.Join(base, "app-1.0.1", "resources")
	touch(t, path.Join(res, "app.asar"))
	touch(t, path.Join(res, "_app.asar"))
	if err := os.MkdirAll(path.Join(base, "app-1.0.2", "resources"), 0755); err != nil {
		t.Fatal(err)
	}

	di := ParseDiscord(base, "")
	if di == nil {
		t.Fatal("ParseDiscord returned nil")
	}
	if di.appPath != path.Join(res, "app") {
		t.Errorf("appPath = %q, want %q", di.appPath, path.Join(res, "app"))
	}
	if !di.isPatched {
		t.Error("install with _app.asar must be reported as patched")
	}
}

// A folder without any app-* entries is not an install.
func TestParseDiscordRejectsDirWithoutApps(t *testing.T) {
	if di := ParseDiscord(t.TempDir(), ""); di != nil {
		t.Errorf("expected nil for dir without app-* folders, got %+v", di)
	}
}
