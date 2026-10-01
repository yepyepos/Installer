/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import "testing"

func TestParseVersion(t *testing.T) {
	tests := []struct {
		tag   string
		valid bool
		nums  []int
		zh    int
	}{
		{"v1.4.2", true, []int{1, 4, 2}, -1},
		{"1.4.2", true, []int{1, 4, 2}, -1},
		{"v1.4.2-zh.1", true, []int{1, 4, 2}, 1},
		{"v1.4.10-zh.2", true, []int{1, 4, 10}, 2},
		{"v2.0-zh.11", true, []int{2, 0}, 11},
		{"devbuild", false, nil, -1},
		{"Unknown", false, nil, -1},
		{"", false, nil, -1},
		{"v1.4.2-rc.1", false, nil, -1},
		{"v1.4.2-zh.x", false, nil, -1},
	}

	for _, tt := range tests {
		v := parseVersion(tt.tag)
		if v.valid != tt.valid {
			t.Errorf("parseVersion(%q).valid = %v, want %v", tt.tag, v.valid, tt.valid)
			continue
		}
		if tt.valid {
			if len(v.nums) != len(tt.nums) {
				t.Errorf("parseVersion(%q).nums = %v, want %v", tt.tag, v.nums, tt.nums)
			} else {
				for i := range tt.nums {
					if v.nums[i] != tt.nums[i] {
						t.Errorf("parseVersion(%q).nums = %v, want %v", tt.tag, v.nums, tt.nums)
					}
				}
			}
			if v.zh != tt.zh {
				t.Errorf("parseVersion(%q).zh = %d, want %d", tt.tag, v.zh, tt.zh)
			}
		}
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		// ordering required by the release strategy
		{"v1.4.2-zh.2", "v1.4.2-zh.1", 1},
		{"v1.4.2-zh.1", "v1.4.2-zh.2", -1},
		{"v1.4.2-zh.1", "v1.4.2", 1},
		{"v1.4.2", "v1.4.2-zh.1", -1},
		{"v1.4.10-zh.1", "v1.4.2-zh.1", 1},
		{"v1.4.2-zh.1", "v1.4.10-zh.1", -1},
		{"v1.5.0-zh.1", "v1.4.10-zh.1", 1},
		{"v1.4.2-zh.1", "v1.5.0-zh.1", -1},
		{"v1.4.2-zh.2", "v1.4.2-zh.2", 0},
		{"v1.4.2", "v1.4.2", 0},
		{"v1.4.2-zh.10", "v1.4.2-zh.9", 1},
		// equal cores with different segment counts
		{"v1.4", "v1.4.0", 0},
		{"v1.4.0.0.1", "v1.4", 1},
		// unparseable tags never compare as newer
		{"devbuild", "v1.4.2-zh.1", 0},
		{"v1.4.2-zh.1", "devbuild", 0},
		{"Unknown", "Unknown", 0},
		{"v1.4.2-rc.1", "v1.4.2-zh.1", 0},
	}

	for _, tt := range tests {
		if got := compareVersions(tt.a, tt.b); got != tt.want {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
