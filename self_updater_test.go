/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"strings"
	"testing"
)

// TestRejectOfficialRepository asserts every installer-release related URL in
// the codebase points at the fork, so a self update can never reach the
// official English installer repository.
func TestRejectOfficialRepository(t *testing.T) {
	urls := []string{
		InstallerReleaseUrl,
		InstallerReleaseUrlFallback,
		// Vencord desktop dist source (I3) must stay on the fork too
		ReleaseUrl,
		ReleaseUrlFallback,
	}

	for _, u := range urls {
		if strings.Contains(u, "api.github.com/repos/Vencord/") ||
			strings.Contains(u, "repos/Vendicated/") ||
			strings.Contains(u, "vencord.dev") {
			t.Errorf("URL %q still references the official repository", u)
		}
		if !strings.Contains(u, "yepyepos") {
			t.Errorf("URL %q does not point at the fork", u)
		}
	}

	// Vencord dist source and installer source must be separate repos.
	if !strings.Contains(ReleaseUrl, "yepyepos/Vencord") {
		t.Errorf("Vencord dist source changed unexpectedly: %q", ReleaseUrl)
	}
	if !strings.Contains(InstallerReleaseUrl, "yepyepos/Installer") {
		t.Errorf("installer update source changed unexpectedly: %q", InstallerReleaseUrl)
	}
}

func TestSelectInstallerRelease(t *testing.T) {
	releases := []GithubRelease{
		{TagName: "v1.4.2-zh.1", Name: "old", PublishedAt: "2026-10-01T00:00:00Z"},
		{TagName: "v1.4.2-zh.2", Name: "new", PublishedAt: "2026-10-02T00:00:00Z"},
		{TagName: "devbuild", Name: "not-a-version", PublishedAt: "2026-10-03T00:00:00Z"},
		{TagName: "v1.4.10-zh.1", Name: "newest", PublishedAt: "2026-10-02T12:00:00Z"},
	}

	got := selectInstallerRelease(releases)
	if got == nil {
		t.Fatal("selectInstallerRelease returned nil")
	}
	if got.TagName != "v1.4.10-zh.1" {
		t.Errorf("selected %q, want v1.4.10-zh.1 (version order, not api order)", got.TagName)
	}

	if got := selectInstallerRelease([]GithubRelease{{TagName: "devbuild"}, {TagName: "Unknown"}}); got != nil {
		t.Errorf("expected nil for releases without valid tags, got %q", got.TagName)
	}
	if got := selectInstallerRelease(nil); got != nil {
		t.Error("expected nil for empty list")
	}
}

func TestInstallerAssetUrl(t *testing.T) {
	// asset URL resolution is covered indirectly via installerAssetName on
	// the current platform; assert the miss path produces a clear error.
	release := &GithubRelease{TagName: "v1.4.2-zh.2"}
	if name := installerAssetName(); name != "" {
		_, err := installerAssetUrl(release)
		if err == nil {
			t.Errorf("expected error for missing asset %q", name)
		} else if !strings.Contains(err.Error(), name) {
			t.Errorf("error %q should mention the missing asset name", err.Error())
		}
	}

	release.Assets = []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	}{
		{Name: "VencordInstaller.exe", DownloadURL: "https://github.com/yepyepos/Installer/releases/download/v1/VencordInstaller.exe"},
		{Name: "VencordInstallerCli.exe", DownloadURL: "https://github.com/yepyepos/Installer/releases/download/v1/VencordInstallerCli.exe"},
	}
	if name := installerAssetName(); name != "" {
		url, err := installerAssetUrl(release)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(url, "yepyepos/Installer") {
			t.Errorf("asset url %q not from fork", url)
		}
	}
}
