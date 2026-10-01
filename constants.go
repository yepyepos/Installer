/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"image/color"
	"vencordinstaller/buildinfo"
)

// The zh-CN fork installs Vencord exclusively from yepyepos/Vencord.
// We use the releases LIST endpoint (not /releases/latest) because the newest
// published release must be selected by actually shipping the required desktop
// dist assets, and pre-releases have to be selectable too.
// The fallback intentionally points at the same endpoint: if GitHub API access
// fails, the installer must fail loudly instead of silently falling back to the
// official English Vencord (Vendicated/Vencord / vencord.dev).
const ReleaseUrl = "https://api.github.com/repos/yepyepos/Vencord/releases?per_page=20"
const ReleaseUrlFallback = "https://api.github.com/repos/yepyepos/Vencord/releases?per_page=20"
// Self-update source: the zh-CN fork updates itself exclusively from
// yepyepos/Installer. Like the Vencord source above we use the releases LIST
// endpoint so pre-releases are selectable and the newest release is picked by
// actual version comparison. The fallback points at the same endpoint on
// purpose: if GitHub API access fails the installer must fail loudly instead
// of silently updating to the official English installer.
const InstallerReleaseUrl = "https://api.github.com/repos/yepyepos/Installer/releases?per_page=20"
const InstallerReleaseUrlFallback = "https://api.github.com/repos/yepyepos/Installer/releases?per_page=20"

var UserAgent = "VencordInstaller/" + buildinfo.InstallerGitHash + " (https://github.com/yepyepos/Installer)"

var (
	DiscordGreen        = color.RGBA{0, 133, 69, 0xff}
	DiscordGreenHovered = color.RGBA{0, 108, 55, 0xff}
	DiscordRed          = color.RGBA{210, 45, 57, 0xff}
	DiscordRedHovered   = color.RGBA{169, 35, 46, 0xff}
	DiscordBlue         = color.RGBA{88, 101, 242, 0xff}
	DiscordBlueHovered  = color.RGBA{68, 82, 187, 0xff}
	DiscordYellow       = color.RGBA{0xfe, 0xe7, 0x5c, 0xff}
)

var LinuxDiscordNames = []string{
	"Discord",
	"DiscordPTB",
	"DiscordCanary",
	"DiscordDevelopment",
	"discord",
	"discordptb",
	"discordcanary",
	"discorddevelopment",
	"discord-ptb",
	"discord-canary",
	"discord-development",
	// Flatpak
	"com.discordapp.Discord",
	"com.discordapp.DiscordPTB",
	"com.discordapp.DiscordCanary",
	"com.discordapp.DiscordDevelopment",
}
