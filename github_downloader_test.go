/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import "testing"

func releaseWithAssets(names []string) GithubRelease {
	r := GithubRelease{}
	for _, n := range names {
		r.Assets = append(r.Assets, struct {
			Name        string `json:"name"`
			DownloadURL string `json:"browser_download_url"`
		}{Name: n})
	}
	return r
}

func TestExtractDesktopHash(t *testing.T) {
	tests := []struct {
		name    string
		release GithubRelease
		want    string
	}{
		{
			name: "contract marker in body",
			release: GithubRelease{
				Name: "Vencord 1.15.7 — 简体中文本地化版",
				Body: "some notes\nVencord-Desktop-Hash: 4c73c063\nmore notes",
			},
			want: "4c73c063",
		},
		{
			name: "contract marker is case insensitive",
			release: GithubRelease{
				Name: "x",
				Body: "VENCORD-DESKTOP-HASH: ABCD123",
			},
			want: "abcd123",
		},
		{
			name: "no marker falls back to upstream name format",
			release: GithubRelease{
				Name: "DevBuild 7f0c10c",
			},
			want: "7f0c10c",
		},
		{
			name: "no marker and no space in name returns whole name",
			release: GithubRelease{
				Name: "NoSpaceName",
			},
			want: "NoSpaceName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractDesktopHash(&tt.release); got != tt.want {
				t.Errorf("ExtractDesktopHash() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHasAllDesktopAssets(t *testing.T) {
	all := []string{"patcher.js", "preload.js", "renderer.js", "renderer.css"}
	if HasAllDesktopAssets(&GithubRelease{}) {
		// sanity: GithubRelease zero value must not qualify
		t.Error("empty release unexpectedly qualifies")
	}
	if r := releaseWithAssets(all); !HasAllDesktopAssets(&r) {
		t.Error("release with all 4 assets should qualify")
	}
	if r := releaseWithAssets([]string{"patcher.js", "patcher.js.map", "preload.js", "renderer.js", "renderer.css", "renderer.js.map"}); !HasAllDesktopAssets(&r) {
		t.Error("release with .map siblings should qualify")
	}
	missing := []string{"patcher.js", "preload.js", "renderer.js"}
	if r := releaseWithAssets(missing); HasAllDesktopAssets(&r) {
		t.Error("release missing renderer.css should NOT qualify")
	}
	userjs := []string{"Vencord.user.js", "extension-chrome.zip", "extension-firefox.zip"}
	if r := releaseWithAssets(userjs); HasAllDesktopAssets(&r) {
		t.Error("browser-extension-only release should NOT qualify")
	}
}

func TestSelectVencordRelease(t *testing.T) {
	stableNoAssets := GithubRelease{
		Name:        "Vencord 1.15.7 — 简体中文本地化版",
		TagName:     "v1.15.7-zh.4",
		PublishedAt: "2026-09-30T09:00:32Z",
		Assets: []struct {
			Name        string `json:"name"`
			DownloadURL string `json:"browser_download_url"`
		}{
			{Name: "extension-chrome.zip"}, {Name: "extension-firefox.zip"}, {Name: "Vencord.user.js"},
		},
	}
	poc := GithubRelease{
		Name:        "Vencord zh-CN Desktop Installer PoC",
		TagName:     "desktop-poc-v1",
		PublishedAt: "2026-09-30T15:37:17Z",
		Assets: []struct {
			Name        string `json:"name"`
			DownloadURL string `json:"browser_download_url"`
		}{
			{Name: "patcher.js"}, {Name: "preload.js"}, {Name: "renderer.js"}, {Name: "renderer.css"},
		},
	}
	officialStyle := GithubRelease{
		Name:        "DevBuild deadbeef",
		TagName:     "devbuild",
		PublishedAt: "2026-01-01T00:00:00Z",
		Assets: []struct {
			Name        string `json:"name"`
			DownloadURL string `json:"browser_download_url"`
		}{
			{Name: "patcher.js"}, {Name: "preload.js"}, {Name: "renderer.js"}, {Name: "renderer.css"},
		},
	}

	t.Run("picks newest qualifying release regardless of api order", func(t *testing.T) {
		got, err := selectVencordRelease([]GithubRelease{stableNoAssets, poc})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.TagName != "desktop-poc-v1" {
			t.Errorf("selected %q, want desktop-poc-v1", got.TagName)
		}
	})

	t.Run("skips releases without desktop assets", func(t *testing.T) {
		if _, err := selectVencordRelease([]GithubRelease{stableNoAssets}); err == nil {
			t.Error("expected error when no release ships desktop assets")
		}
	})

	t.Run("empty list errors", func(t *testing.T) {
		if _, err := selectVencordRelease(nil); err == nil {
			t.Error("expected error for empty release list")
		}
	})

	t.Run("selection is content based, not name based", func(t *testing.T) {
		// an official-style "DevBuild <hash>" release must not be chosen over a
		// newer qualifying zh release merely by name; selection follows published_at
		got, err := selectVencordRelease([]GithubRelease{officialStyle, poc})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.TagName != "desktop-poc-v1" {
			t.Errorf("selected %q, want desktop-poc-v1", got.TagName)
		}
	})
}
