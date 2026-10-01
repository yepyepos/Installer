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

func TestDistHasAllDesktopFiles(t *testing.T) {
	dir := t.TempDir()
	if distHasAllDesktopFiles(dir) {
		t.Error("empty dist must not count as complete")
	}

	names := []string{"patcher.js", "preload.js", "renderer.js", "renderer.css"}
	for _, n := range names {
		if err := os.WriteFile(path.Join(dir, n), []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if !distHasAllDesktopFiles(dir) {
		t.Error("dist with all four files must count as complete")
	}

	if err := os.Remove(path.Join(dir, "renderer.js")); err != nil {
		t.Fatal(err)
	}
	if distHasAllDesktopFiles(dir) {
		t.Error("dist missing renderer.js must not count as complete")
	}

	// extra files (package.json, .map) must not affect completeness
	if err := os.WriteFile(path.Join(dir, "renderer.js"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path.Join(dir, "package.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	if !distHasAllDesktopFiles(dir) {
		t.Error("dist with extra files must still count as complete")
	}
}
