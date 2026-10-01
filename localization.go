/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"os"
	"strings"
)

// UI holds every user visible string of the installer (GUI, CLI and shared
// error/log messages). Technical identifiers (product names, popup ids, flag
// names, URLs, file names) intentionally stay out of this struct.
type UI struct {
	// Shared GUI modals
	ModalOhNo           string
	ModalNeverSee       string
	ModalSupportHint    string
	ModalVisitSupport   string
	ModalTitlePatch     string
	ModalTitleUnpatch   string
	ModalTitleOAInstall string
	ModalTitleOAUninst  string
	ErrPermWindows      string
	ErrPermUnix         string
	ModalInstallFailed  string
	ModalInstallFailedDesc string

	// GUI main window
	FilesDirErrTitle      string
	FilesDirErrHint       string
	OfficialSourcesWarn   string
	SelectInstallPrompt   string
	NoDiscordFound        string
	SnapUnsupported       string
	AdvancedUserCheckbox  string
	VencordInstalledSfx   string
	CustomLocationRadio   string
	CustomLocationHint    string
	BtnInstall            string
	BtnReinstallRepair    string
	BtnUninstall          string
	BtnUninstallOpenAsar  string
	BtnInstallOpenAsar    string
	BtnUninstInstallOA    string
	TipInstall            string
	TipReinstall          string
	TipUninstall          string
	TipOpenAsar           string
	GithubErrorCard       string

	// GUI modals
	BtnTakeMeThere       string
	BtnOpenSettings      string
	BtnAccept            string
	BtnCancel            string
	BtnOk                string
	ModalPatchedTitle    string
	ModalPatchedDesc     string
	ModalUnpatchedTitle  string
	ModalUnpatchedDesc   string
	ModalScuffedTitle    string
	ModalScuffedDesc     string
	OpenAsarConfirmDesc  string
	ModalPermsTitle      string
	ModalPermsDesc       string
	ModalOAPatchedTitle  string
	ModalOAPatchedDesc   string
	ModalOAUnpatchedTitle string
	ModalOAUnpatchedDesc string
	ModalInvalidLocTitle string
	ModalInvalidLocDesc  string
	UpdateTitle          string
	UpdateBody           string
	BtnUpdateNow         string
	BtnLater             string
	ModalSelfUpdateFailTitle string
	ModalSelfUpdateFailDesc  string
	ModalSelfRestartFailTitle string
	ModalSelfRestartFailDesc  string

	// CLI flags & menus
	CliFlagDebug      string
	CliFlagHelp       string
	CliFlagVersion    string
	CliFlagUpdateSelf string
	CliFlagInstall    string
	CliFlagRepair     string
	CliFlagUninstall  string
	CliFlagInstallOA  string
	CliFlagUninstallOA string
	CliFlagLocation   string
	CliFlagBranch     string
	CliMenuHelp       string
	CliMenuUpdate     string
	CliMenuQuit       string
	CliMenuPrompt     string
	CliWarnOutdated   string
	CliWarnOutdatedHint string
	CliDieUpdateCheckFail string
	CliErrFlagsExclusive  string
	CliErrBranchInvalid   string
	CliDieFetchFailInstall string
	CliDieFetchFailRepair  string
	CliOpenAsarAlready    string
	CliOpenAsarNotInstalled string
	CliExitSuccess        string
	CliExitFailure        string
	CliPressEnterToExit   string
	CliNoDiscord          string
	CliBranchNotFoundFmt  string
	CliInvalidLocationFmt string
	CliVencordInstalledTag string
	CliSelectPatch        string
	CliSelectUnpatch      string
	CliSelectRepair       string
	CliCustomLocationItem string
	CliCustomLocationPrompt string
	CliInvalidDiscord     string
	CliScuffedBody1       string
	CliScuffedBody2       string
	CliScuffedBody3       string
	CliDownloading        string
	CliDone               string

	// Self updater
	SelfCheckFail         string
	SelfNoValidRelease    string
	SelfUpToDateErr       string
	SelfNoAssetErr        string
	SelfNoAssetFmt        string
	SelfDownloadStatusFmt string
	SelfTempFailFmt       string
	SelfRemoveFailFmt     string
	SelfReplaceFailFmt    string
	SelfStartFailFmt      string
	SelfReleaseProcFailFmt string
	SelfLengthMismatchFmt  string
	SelfRollbackOk         string
	CliSelfUpdateFail      string

	// Shared log & error messages
	LogPatchingFmt     string
	LogAlreadyPatchedFmt string
	LogPatched         string
	LogUnpatchingFmt   string
	LogUnpatched       string
	LogUndoPatchFailed string
	LogUndoUnpatchFailed string
	LogUndoPatchBricked  string
	LogUndoUnpatchBricked string
	LogUndidAllChanges string
	LogCreateFail      string
	LogDownloadFail    string
	LogDownloadToFail  string
	LogFetchDataFailed string
	LogContractViolation string
	LogTmpBackupFail     string
	LogCreateRequestFail string
	LogSendFail          string
	LogFetchFallbackFmt  string
	LogNonOKStatus       string
	LogDecodeFail        string
	ErrUnpatchPatchedFmt string
	ErrBusyDiscord     string
	ErrFlatpakFmt      string
	ErrNoAsarFmt       string
	ErrOpenAsarFetchFmt string
	ErrNoBackup        string
	ErrNoDesktopAssets string
	ErrMissingFiles    string
}

// English is the fallback locale: every field must be non empty so an
// unavailable locale can never yield blank UI text.
var English = UI{
	ModalOhNo:                "Oh No :(",
	ModalNeverSee:            "You should never see this",
	ModalSupportHint:         "If this issue persists, visit https://vencord.dev/support for help.",
	ModalVisitSupport:        "If this issue persists, visit: https://vencord.dev/support",
	ModalTitlePatch:          "Failed to patch this Install.",
	ModalTitleUnpatch:        "Failed to unpatch this Install.",
	ModalTitleOAInstall:      "Failed to install OpenAsar on this Install.",
	ModalTitleOAUninst:       "Failed to uninstall OpenAsar from this Install.",
	ErrPermWindows:           "Permission denied. Make sure your Discord is fully closed (from the tray)!",
	ErrPermUnix:              "Permission denied. Maybe try running me as Administrator/Root?",
	ModalInstallFailed:       "Failed to install the latest Vencord builds from GitHub",
	ModalInstallFailedDesc:   "If this issue persists, visit https://vencord.dev/support for help.",

	FilesDirErrTitle:         "Error: Failed to create: ",
	FilesDirErrHint:          "Resolve this error, then restart me!",
	OfficialSourcesWarn:      "**Github** and **vencord.dev** are the only official places to get Vencord. Any other site claiming to be us is malicious.\n" +
		"If you downloaded from any other source, you should delete / uninstall everything immediately, run a malware scan and change your Discord password.",
	SelectInstallPrompt:      "Please select an install to patch",
	NoDiscordFound:           "No Discord installs found. You first need to install Discord.",
	SnapUnsupported:          " snap is not supported.",
	AdvancedUserCheckbox:     "I am an advanced user and have Discord installed at a different location",
	VencordInstalledSfx:      " (Vencord Installed)",
	CustomLocationRadio:      "Custom Install Location",
	CustomLocationHint:       "The custom location",
	BtnInstall:               "Install",
	BtnReinstallRepair:       "Reinstall / Repair",
	BtnUninstall:             "Uninstall",
	BtnUninstallOpenAsar:     "Uninstall OpenAsar",
	BtnInstallOpenAsar:       "Install OpenAsar",
	BtnUninstInstallOA:       "(Un-)Install OpenAsar",
	TipInstall:               "Patch the selected Discord Install",
	TipReinstall:             "Reinstall & Update Vencord",
	TipUninstall:             "Unpatch the selected Discord Install",
	TipOpenAsar:              "Manage OpenAsar",
	GithubErrorCard:          "Failed to fetch Info from GitHub. If this issue persists, visit https://vencord.dev/support for help.",

	BtnTakeMeThere:           "Take me there!",
	BtnOpenSettings:          "Open Settings",
	BtnAccept:                "Accept",
	BtnCancel:                "Cancel",
	BtnOk:                    "Ok",
	ModalPatchedTitle:        "Installed!",
	ModalPatchedDesc:         "Vencord was successfully installed!",
	ModalUnpatchedTitle:      "Uninstalled",
	ModalUnpatchedDesc:       "Vencord has been uninstalled!",
	ModalScuffedTitle:        "Hold On!",
	ModalScuffedDesc:         "You have a broken Discord Install.\n" +
		"Sometimes Discord decides to install to the wrong location for some reason!\n" +
		"You need to fix this before patching, otherwise Vencord will likely not work.\n\n" +
		"Use the below button to jump there and delete any folder called Discord or Squirrel.\n" +
		"If the folder is now empty, feel free to go back a step and delete that folder too.\n" +
		"Then see if Discord still starts. If not, reinstall it",
	OpenAsarConfirmDesc:      "OpenAsar is an open-source alternative of Discord desktop's app.asar.\n" +
		"Vencord is in no way affiliated with OpenAsar.\n" +
		"You're installing OpenAsar at your own risk. If you run into issues with OpenAsar,\n" +
		"no support will be provided, join the OpenAsar Server instead!\n\n" +
		"To install OpenAsar, press Accept and click 'Install OpenAsar' again.",
	ModalPermsTitle:          "Insufficient Permissions",
	ModalPermsDesc:           "Permission denied. Please grant the installer permissions in the settings.",
	ModalOAPatchedTitle:      "Successfully Installed OpenAsar",
	ModalOAPatchedDesc:       "If Discord is still open, fully close it first. Then start it again and verify OpenAsar installed successfully!",
	ModalOAUnpatchedTitle:    "Successfully Uninstalled OpenAsar",
	ModalOAUnpatchedDesc:     "If Discord is still open, fully close it first. Then start it again and it should be back to stock!",
	ModalInvalidLocTitle:     "Invalid Location",
	ModalInvalidLocDesc:      "The specified location is not a valid Discord install.\nMake sure you select the base folder.\n\nHint: Discord snap is not supported. use flatpak or .deb",
	UpdateTitle:              "Your Installer is outdated!",
	UpdateBody:               "Would you like to update now?\n\n" +
		"Once you press Update Now, the new installer will automatically be downloaded.\n" +
		"The installer will temporarily seem unresponsive. Just wait!\n" +
		"Once the update is done, the Installer will automatically reopen.\n\n" +
		"On MacOs, Auto updates are not supported, so it will instead open in browser.",
	BtnUpdateNow:             "Update Now",
	BtnLater:                 "Later",
	ModalSelfUpdateFailTitle: "Failed to update self!",
	ModalSelfUpdateFailDesc:  "Please manually download the latest Installer.",
	ModalSelfRestartFailTitle: "Failed to restart self!",
	ModalSelfRestartFailDesc:  "Please manually restart the Installer.",

	CliFlagDebug:       "Enable debug info",
	CliFlagHelp:        "View usage instructions",
	CliFlagVersion:     "View the program version",
	CliFlagUpdateSelf:  "Update me to the latest version",
	CliFlagInstall:     "Install Vencord",
	CliFlagRepair:      "Repair Vencord",
	CliFlagUninstall:   "Uninstall Vencord",
	CliFlagInstallOA:   "Install OpenAsar",
	CliFlagUninstallOA: "Uninstall OpenAsar",
	CliFlagLocation:    "The location of the Discord install to modify",
	CliFlagBranch:      "The branch of Discord to modify [auto|stable|ptb|canary]",
	CliMenuHelp:        "View Help Menu",
	CliMenuUpdate:      "Update Vencord Installer",
	CliMenuQuit:        "Quit",
	CliMenuPrompt:      "What would you like to do? (Press Enter to confirm)",
	CliWarnOutdated:    "Your installer is outdated.",
	CliWarnOutdatedHint: "To update, select the 'Update Vencord Installer' option to update, or run with --update-self",
	CliDieUpdateCheckFail: "Can't update self because checking for updates failed",
	CliErrFlagsExclusive:  "The 'location' and 'branch' flags are mutually exclusive.",
	CliErrBranchInvalid:   "The 'branch' flag must be one of the following: [auto|stable|ptb|canary]",
	CliDieFetchFailInstall: "Not installing as fetching release data failed. If this issue persists, see https://vencord.dev/support",
	CliDieFetchFailRepair:  "Not updating as fetching release data failed. If this issue persists, see https://vencord.dev/support",
	CliOpenAsarAlready:     "OpenAsar already installed",
	CliOpenAsarNotInstalled: "OpenAsar not installed",
	CliExitSuccess:         "✔ Success!",
	CliExitFailure:         "❌ Failed! If this issue persists, see https://vencord.dev/support",
	CliPressEnterToExit:    "Press Enter to exit",
	CliNoDiscord:           "No Discord install found. Before proceeding, make sure Discord is installed. snap is not supported!",
	CliBranchNotFoundFmt:   "Discord %s not found",
	CliInvalidLocationFmt:  "%s is not a valid Discord install. Hint: snap is not supported",
	CliVencordInstalledTag: " [Vencord Installed]",
	CliSelectPatch:         "Select Discord install to patch (Press Enter to confirm)",
	CliSelectUnpatch:       "Select Discord install to unpatch (Press Enter to confirm)",
	CliSelectRepair:        "Select Discord install to repair (Press Enter to confirm)",
	CliCustomLocationItem:  "Custom Location",
	CliCustomLocationPrompt: "Custom Discord Location",
	CliInvalidDiscord:      "Invalid Discord install!",
	CliScuffedBody1:        "Hold On!",
	CliScuffedBody2:        "You have a broken Discord Install.",
	CliScuffedBody3:        "Please reinstall Discord before proceeding! Otherwise, Vencord will likely not work.",
	CliDownloading:         "Downloading latest Vencord files...",
	CliDone:                "Done!",

	SelfCheckFail:          "Failed to check for installer updates:",
	SelfNoValidRelease:     "no yepyepos/Installer release with a valid version tag found",
	SelfUpToDateErr:        "Cannot update self. Either no update available or macos",
	SelfNoAssetErr:         "update asset not found",
	SelfNoAssetFmt:         "update asset '%s' not found in the latest release",
	SelfDownloadStatusFmt:  "Failed to download update (status code %d): %s",
	SelfTempFailFmt:        "Failed to create tempfile: %w",
	SelfRemoveFailFmt:      "Failed to remove/rename own executable: %w",
	SelfReplaceFailFmt:     "Failed to replace self with updated executable. Please manually redownload the installer: %w",
	SelfStartFailFmt:       "Failed to start new process: %w",
	SelfReleaseProcFailFmt: "Failed to release new process: %w",
	SelfLengthMismatchFmt:  "Unexpected end of input. Content-Length was %s, but I only read %s",
	SelfRollbackOk:         "Update failed, previous executable restored",
	CliSelfUpdateFail:      "Failed to update self:",

	LogPatchingFmt:        "Patching %s...",
	LogAlreadyPatchedFmt:  "%s is already patched. Unpatching first...",
	LogPatched:            "Successfully patched",
	LogUnpatchingFmt:      "Unpatching %s...",
	LogUnpatched:          "Successfully unpatched",
	LogUndoPatchFailed:    "Failed to patch. Undoing partial patch",
	LogUndoUnpatchFailed:  "Failed to unpatch. Undoing partial unpatch",
	LogUndoPatchBricked:   "Failed to undo partial patch. This install is probably bricked.",
	LogUndoUnpatchBricked: "Failed to undo partial unpatch. This install is probably bricked.",
	LogUndidAllChanges:    "Successfully undid all changes",
	LogCreateFail:         "Failed to create",
	LogDownloadFail:       "Failed to download",
	LogDownloadToFail:     "Failed to download to",
	LogFetchDataFailed:    "Failed to fetch Vencord release data:",
	LogContractViolation:  "Release contract violation: downloaded patcher.js hash does not match the hash advertised by the release!",
	LogTmpBackupFail:      "Failed to delete temporary app.asar (patch folder) backup. This is whatever but you might want to delete it manually.",
	LogCreateRequestFail:  "Failed to create Request",
	LogSendFail:           "Failed to send Request",
	LogFetchFallbackFmt:   "Failed to fetch %s (status code %d). Trying fallback url %s",
	LogNonOKStatus:        "returned Non-OK status",
	LogDecodeFail:         "Failed to decode GitHub JSON Response",
	ErrUnpatchPatchedFmt:  "patch: Failed to unpatch already patched install '%s':\n%v",
	ErrBusyDiscord:        "Cannot patch because Discord's files are used by a different process." +
		"\nMake sure you close Discord before trying to patch!",
	ErrFlatpakFmt:        "Failed to grant Discord Flatpak access to %s: %v",
	ErrNoAsarFmt:         "Install at %s has no asar file",
	ErrOpenAsarFetchFmt:  "Failed to fetch OpenAsar - %s: %s",
	ErrNoBackup:          "No app.asar.backup. Reinstall Discord",
	ErrNoDesktopAssets:   "no yepyepos/Vencord release ships the required Vencord desktop files (patcher.js, preload.js, renderer.js, renderer.css)",
	ErrMissingFiles:      "Couldn't find all required files",
}

// Chinese is the default locale of the zh-CN installer.
var Chinese = UI{
	ModalOhNo:                "出错了 :(",
	ModalNeverSee:            "你不应该看到这条消息",
	ModalSupportHint:         "如果此问题持续出现，请访问 https://vencord.dev/support 获取帮助。",
	ModalVisitSupport:        "如果此问题持续出现，请访问：https://vencord.dev/support",
	ModalTitlePatch:          "修补此 Discord 安装失败。",
	ModalTitleUnpatch:        "卸载此 Discord 安装失败。",
	ModalTitleOAInstall:      "在此安装上安装 OpenAsar 失败。",
	ModalTitleOAUninst:       "从此安装卸载 OpenAsar 失败。",
	ErrPermWindows:           "权限被拒绝。请确保 Discord 已完全关闭（包括托盘图标）！",
	ErrPermUnix:              "权限被拒绝。请尝试以管理员/Root 身份运行？",
	ModalInstallFailed:       "从 GitHub 安装最新 Vencord 构建失败",
	ModalInstallFailedDesc:   "如果此问题持续出现，请访问 https://vencord.dev/support 获取帮助。",

	FilesDirErrTitle:         "错误：创建目录失败：",
	FilesDirErrHint:          "解决此问题后，请重新运行安装器！",
	OfficialSourcesWarn:      "**GitHub** 和 **vencord.dev** 是获取 Vencord 的唯一官方渠道，其他任何声称是我们的网站都是恶意的。\n" +
		"如果你从其他来源下载了 Vencord，请立即删除/卸载所有相关内容，进行恶意软件扫描，并修改你的 Discord 密码。",
	SelectInstallPrompt:      "请选择要修补的 Discord 安装",
	NoDiscordFound:           "未找到 Discord 安装，请先安装 Discord。",
	SnapUnsupported:          "（不支持 snap。）",
	AdvancedUserCheckbox:     "我是高级用户，Discord 安装在其他位置",
	VencordInstalledSfx:      "（已安装 Vencord）",
	CustomLocationRadio:      "自定义安装位置",
	CustomLocationHint:       "输入自定义位置的路径",
	BtnInstall:               "安装",
	BtnReinstallRepair:       "重新安装 / 修复",
	BtnUninstall:             "卸载",
	BtnUninstallOpenAsar:     "卸载 OpenAsar",
	BtnInstallOpenAsar:       "安装 OpenAsar",
	BtnUninstInstallOA:       "安装/卸载 OpenAsar",
	TipInstall:               "修补选中的 Discord 安装",
	TipReinstall:             "重新安装并更新 Vencord",
	TipUninstall:             "卸载选中的 Discord 安装",
	TipOpenAsar:              "管理 OpenAsar",
	GithubErrorCard:          "从 GitHub 获取信息失败。如果此问题持续出现，请访问 https://vencord.dev/support 获取帮助。",

	BtnTakeMeThere:           "打开该目录",
	BtnOpenSettings:          "打开设置",
	BtnAccept:                "接受",
	BtnCancel:                "取消",
	BtnOk:                    "确定",
	ModalPatchedTitle:        "安装成功！",
	ModalPatchedDesc:         "Vencord 已成功安装！",
	ModalUnpatchedTitle:      "已卸载",
	ModalUnpatchedDesc:       "Vencord 已成功卸载！",
	ModalScuffedTitle:        "请稍等！",
	ModalScuffedDesc:         "检测到损坏的 Discord 安装。\n" +
		"有时 Discord 会因某些原因安装到错误的位置！\n" +
		"修补前需要先解决此问题，否则 Vencord 很可能无法正常工作。\n\n" +
		"点击下方按钮跳转到该目录，删除所有名为 Discord 或 Squirrel 的文件夹。\n" +
		"如果删除后文件夹已为空，可以返回上一级目录并删除该空文件夹。\n" +
		"然后确认 Discord 是否还能启动；若无法启动，请重新安装 Discord。",
	OpenAsarConfirmDesc:      "OpenAsar 是 Discord 桌面版 app.asar 的开源替代品。\n" +
		"Vencord 与 OpenAsar 没有任何关联。\n" +
		"安装 OpenAsar 的风险由你自行承担；如果遇到 OpenAsar 相关问题，\n" +
		"我们将不提供支持，请前往 OpenAsar 社区寻求帮助！\n\n" +
		"如需安装 OpenAsar，请先点击“接受”，然后再次点击“安装 OpenAsar”。",
	ModalPermsTitle:          "权限不足",
	ModalPermsDesc:           "权限被拒绝。请在系统设置中为安装器授予相应权限。",
	ModalOAPatchedTitle:      "OpenAsar 安装成功",
	ModalOAPatchedDesc:       "如果 Discord 仍在运行，请先将其完全关闭，然后重新启动，并确认 OpenAsar 已成功安装！",
	ModalOAUnpatchedTitle:    "OpenAsar 卸载成功",
	ModalOAUnpatchedDesc:     "如果 Discord 仍在运行，请先将其完全关闭，然后重新启动，Discord 应恢复原状！",
	ModalInvalidLocTitle:     "无效的位置",
	ModalInvalidLocDesc:      "指定的位置不是有效的 Discord 安装。\n请确保选择的是其根目录。\n\n提示：不支持 snap 安装的 Discord，请使用 flatpak 或 .deb 版本。",
	UpdateTitle:              "安装器已有新版本！",
	UpdateBody:               "要立即更新吗？\n\n" +
		"点击“立即更新”后，新版安装器将自动下载。\n" +
		"期间安装器可能暂时无响应，请耐心等待！\n" +
		"更新完成后，安装器会自动重新打开。\n\n" +
		"macOS 不支持自动更新，将改为在浏览器中打开。",
	BtnUpdateNow:             "立即更新",
	BtnLater:                 "稍后",
	ModalSelfUpdateFailTitle: "自动更新失败！",
	ModalSelfUpdateFailDesc:  "请手动下载最新版安装器。",
	ModalSelfRestartFailTitle: "重启安装器失败！",
	ModalSelfRestartFailDesc:  "请手动重启安装器。",

	CliFlagDebug:       "启用调试信息",
	CliFlagHelp:        "查看使用说明",
	CliFlagVersion:     "查看程序版本",
	CliFlagUpdateSelf:  "将安装器自身更新到最新版本",
	CliFlagInstall:     "安装 Vencord",
	CliFlagRepair:      "修复 Vencord",
	CliFlagUninstall:   "卸载 Vencord",
	CliFlagInstallOA:   "安装 OpenAsar",
	CliFlagUninstallOA: "卸载 OpenAsar",
	CliFlagLocation:    "要修改的 Discord 安装位置",
	CliFlagBranch:      "要修改的 Discord 分支 [auto|stable|ptb|canary]",
	CliMenuHelp:        "查看帮助",
	CliMenuUpdate:      "更新 Vencord 安装器",
	CliMenuQuit:        "退出",
	CliMenuPrompt:      "请选择要执行的操作（按回车确认）",
	CliWarnOutdated:    "安装器已有新版本。",
	CliWarnOutdatedHint: "如需更新，请选择“更新 Vencord 安装器”，或使用 --update-self 参数重新运行",
	CliDieUpdateCheckFail: "无法自更新：检查更新失败",
	CliErrFlagsExclusive:  "--location 与 --branch 参数互斥，不能同时使用。",
	CliErrBranchInvalid:   "--branch 参数必须为以下之一：[auto|stable|ptb|canary]",
	CliDieFetchFailInstall: "获取 Release 数据失败，已取消安装。如果此问题持续出现，请访问 https://vencord.dev/support",
	CliDieFetchFailRepair:  "获取 Release 数据失败，已取消修复。如果此问题持续出现，请访问 https://vencord.dev/support",
	CliOpenAsarAlready:     "OpenAsar 已安装",
	CliOpenAsarNotInstalled: "OpenAsar 未安装",
	CliExitSuccess:         "✔ 操作成功！",
	CliExitFailure:         "❌ 操作失败！如果此问题持续出现，请访问 https://vencord.dev/support",
	CliPressEnterToExit:    "按回车键退出",
	CliNoDiscord:           "未找到 Discord 安装。请先安装 Discord（不支持 snap）！",
	CliBranchNotFoundFmt:   "未找到 Discord %s 安装",
	CliInvalidLocationFmt:  "%s 不是有效的 Discord 安装（提示：不支持 snap）",
	CliVencordInstalledTag: " [已安装 Vencord]",
	CliSelectPatch:         "选择要修补的 Discord 安装（按回车确认）",
	CliSelectUnpatch:       "选择要卸载的 Discord 安装（按回车确认）",
	CliSelectRepair:        "选择要修复的 Discord 安装（按回车确认）",
	CliCustomLocationItem:  "自定义位置",
	CliCustomLocationPrompt: "自定义 Discord 位置",
	CliInvalidDiscord:      "无效的 Discord 安装！",
	CliScuffedBody1:        "请稍等！",
	CliScuffedBody2:        "检测到损坏的 Discord 安装。",
	CliScuffedBody3:        "请先重新安装 Discord 再继续！否则 Vencord 很可能无法正常工作。",
	CliDownloading:         "正在下载最新 Vencord 文件…",
	CliDone:                "完成！",

	SelfCheckFail:          "检查安装器更新失败：",
	SelfNoValidRelease:     "没有找到包含有效版本号的 yepyepos/Installer Release",
	SelfUpToDateErr:        "无法自更新：已是最新版本，或当前平台不支持自动更新",
	SelfNoAssetErr:         "未找到更新文件",
	SelfNoAssetFmt:         "最新 Release 中未找到更新文件 %s",
	SelfDownloadStatusFmt:  "下载更新失败（状态码 %d）：%s",
	SelfTempFailFmt:        "创建临时文件失败：%w",
	SelfRemoveFailFmt:      "移除/重命名自身可执行文件失败：%w",
	SelfReplaceFailFmt:     "替换自身失败，请手动重新下载安装器：%w",
	SelfStartFailFmt:       "启动新进程失败：%w",
	SelfReleaseProcFailFmt: "释放新进程句柄失败：%w",
	SelfLengthMismatchFmt:  "下载不完整：Content-Length 为 %s，实际读取 %s",
	SelfRollbackOk:         "更新失败，已恢复旧版安装器",
	CliSelfUpdateFail:      "安装器自更新失败：",

	LogPatchingFmt:        "正在修补 %s…",
	LogAlreadyPatchedFmt:  "%s 已安装过 Vencord，正在先卸载旧补丁…",
	LogPatched:            "修补成功",
	LogUnpatchingFmt:      "正在卸载 %s…",
	LogUnpatched:          "卸载成功",
	LogUndoPatchFailed:    "修补失败，正在回滚已执行的更改",
	LogUndoUnpatchFailed:  "卸载失败，正在回滚已执行的更改",
	LogUndoPatchBricked:   "回滚失败，此 Discord 安装可能已损坏。",
	LogUndoUnpatchBricked: "回滚失败，此 Discord 安装可能已损坏。",
	LogUndidAllChanges:    "已成功回滚所有更改",
	LogCreateFail:         "创建失败",
	LogDownloadFail:       "下载失败",
	LogDownloadToFail:     "下载写入失败",
	LogFetchDataFailed:    "获取 Vencord Release 数据失败：",
	LogContractViolation:  "Release 契约校验失败：下载的 patcher.js 构建哈希与 Release 声明不一致！",
	LogTmpBackupFail:      "删除临时 app.asar 备份失败。影响不大，但你可以稍后手动删除它。",
	LogCreateRequestFail:  "创建请求失败",
	LogSendFail:           "发送请求失败",
	LogFetchFallbackFmt:   "获取 %s 失败（状态码 %d），正在尝试备用地址 %s",
	LogNonOKStatus:        "返回了非正常状态码",
	LogDecodeFail:         "解析 GitHub JSON 响应失败",
	ErrUnpatchPatchedFmt:  "修补失败：卸载 %s 上已有补丁时出错：\n%v",
	ErrBusyDiscord:        "无法修补：Discord 的文件正被其他进程占用。\n请先完全关闭 Discord，再尝试修补！",
	ErrFlatpakFmt:         "授予 Discord Flatpak 访问 %s 的权限失败：%v",
	ErrNoAsarFmt:          "在 %s 中未找到 asar 文件",
	ErrOpenAsarFetchFmt:   "下载 OpenAsar 失败 - %s: %s",
	ErrNoBackup:           "未找到 app.asar.backup，请重新安装 Discord",
	ErrNoDesktopAssets:    "没有找到包含全部所需 Vencord 桌面文件（patcher.js、preload.js、renderer.js、renderer.css）的 yepyepos/Vencord Release",
	ErrMissingFiles:       "未能获取全部必需的 Vencord 文件",
}

// pickLanguage selects the UI locale. The zh-CN installer defaults to
// Simplified Chinese; setting VENCORD_INSTALLER_LANG=en switches back to
// English. Anything unrecognised falls back to Chinese with English fields
// never being blank (see English above).
func pickLanguage() *UI {
	lang := strings.TrimSpace(os.Getenv("VENCORD_INSTALLER_LANG"))
	if strings.EqualFold(lang, "en") || strings.EqualFold(lang, "english") {
		return &English
	}
	return &Chinese
}

// L is the locale used by every user visible string in the app.
var L = pickLanguage()
