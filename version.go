/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"strconv"
	"strings"
)

// version represents a release tag of this installer fork, e.g.
// "v1.4.2", "v1.4.2-zh.1" or "v1.4.10-zh.2".
type version struct {
	nums  []int  // numeric core segments, e.g. 1.4.2 -> [1, 4, 2]
	zh    int    // zh-CN iteration of the fork, -1 when the tag has no -zh.N suffix
	valid bool   // false when the tag could not be parsed; incomparable then
	tag   string // original tag for diagnostics
}

// parseVersion parses fork release tags. Unparseable tags (official style
// tags like "devbuild", local builds reporting "Unknown", etc.) yield
// valid=false and are never considered newer.
func parseVersion(tag string) version {
	v := version{zh: -1, tag: tag}
	t := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(tag), "v"))
	if t == "" {
		return v
	}

	if i := strings.Index(t, "-zh."); i >= 0 {
		n, err := strconv.Atoi(strings.TrimSpace(t[i+4:]))
		if err != nil || n < 0 {
			return v
		}
		v.zh = n
		t = t[:i]
	} else if strings.Contains(t, "-") {
		// Some other pre-release suffix we do not understand.
		return v
	}

	for _, part := range strings.Split(t, ".") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || n < 0 {
			return v
		}
		v.nums = append(v.nums, n)
	}

	if len(v.nums) == 0 {
		return v
	}

	v.valid = true
	return v
}

// compareVersions returns 1 if a is newer than b, -1 if a is older than b,
// and 0 if they are equal. Tags that cannot be parsed never compare greater.
func compareVersions(a, b string) int {
	va, vb := parseVersion(a), parseVersion(b)
	if !va.valid || !vb.valid {
		return 0
	}

	n := len(va.nums)
	if len(vb.nums) > n {
		n = len(vb.nums)
	}
	for i := 0; i < n; i++ {
		x, y := 0, 0
		if i < len(va.nums) {
			x = va.nums[i]
		}
		if i < len(vb.nums) {
			y = vb.nums[i]
		}
		if x != y {
			return Ternary(x > y, 1, -1)
		}
	}

	// v1.4.2 < v1.4.2-zh.1 < v1.4.2-zh.2
	if (va.zh >= 0) != (vb.zh >= 0) {
		return Ternary(va.zh >= 0, 1, -1)
	}
	if va.zh != vb.zh {
		return Ternary(va.zh > vb.zh, 1, -1)
	}

	return 0
}
