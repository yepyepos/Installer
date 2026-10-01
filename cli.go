//go:build cli

/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"vencordinstaller/buildinfo"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

var discords []any
var interactive = false

func isValidBranch(branch string) bool {
	switch branch {
	case "", "stable", "ptb", "canary", "auto":
		return true
	default:
		return false
	}
}

func die(msg string) {
	Log.Error(msg)
	exitFailure()
}

func main() {
	enableUTF8Console()
	InitGithubDownloader()
	discords = FindDiscords()

	// Used by log.go init func
	flag.Bool("debug", false, L.CliFlagDebug)

	var helpFlag = flag.Bool("help", false, L.CliFlagHelp)
	var versionFlag = flag.Bool("version", false, L.CliFlagVersion)
	var updateSelfFlag = flag.Bool("update-self", false, L.CliFlagUpdateSelf)
	var installFlag = flag.Bool("install", false, L.CliFlagInstall)
	var updateFlag = flag.Bool("repair", false, L.CliFlagRepair)
	var uninstallFlag = flag.Bool("uninstall", false, L.CliFlagUninstall)
	var installOpenAsarFlag = flag.Bool("install-openasar", false, L.CliFlagInstallOA)
	var uninstallOpenAsarFlag = flag.Bool("uninstall-openasar", false, L.CliFlagUninstallOA)
	var locationFlag = flag.String("location", "", L.CliFlagLocation)
	var branchFlag = flag.String("branch", "", L.CliFlagBranch)
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	if *versionFlag {
		fmt.Println("Vencord Installer Cli", buildinfo.InstallerTag, "("+buildinfo.InstallerGitHash+")")
		fmt.Println("Copyright (C) 2023 Vendicated and Vencord contributors")
		fmt.Println("License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>.")
		return
	}

	if *updateSelfFlag {
		if !<-SelfUpdateCheckDoneChan {
			die(L.CliDieUpdateCheckFail)
		}
		if err := UpdateSelf(); err != nil {
			Log.Error(L.CliSelfUpdateFail, err)
			exitFailure()
		}
		exitSuccess()
	}

	if *locationFlag != "" && *branchFlag != "" {
		die(L.CliErrFlagsExclusive)
	}

	if !isValidBranch(*branchFlag) {
		die(L.CliErrBranchInvalid)
	}

	if *installFlag || *updateFlag {
		if !<-GithubDoneChan {
			die(Ternary(*installFlag, L.CliDieFetchFailInstall, L.CliDieFetchFailRepair))
		}
	}

	install, uninstall, update, installOpenAsar, uninstallOpenAsar := *installFlag, *uninstallFlag, *updateFlag, *installOpenAsarFlag, *uninstallOpenAsarFlag
	switches := []*bool{&install, &update, &uninstall, &installOpenAsar, &uninstallOpenAsar}
	if !SliceContainsFunc(switches, func(b *bool) bool { return *b }) {
		interactive = true

		go func() {
			<-SelfUpdateCheckDoneChan
			if IsSelfOutdated {
				Log.Warn(L.CliWarnOutdated)
				Log.Warn(L.CliWarnOutdatedHint)
			}
		}()

		choices := []string{
			L.CliFlagInstall,
			L.CliFlagRepair,
			L.CliFlagUninstall,
			L.CliFlagInstallOA,
			L.CliFlagUninstallOA,
			L.CliMenuHelp,
			L.CliMenuUpdate,
			L.CliMenuQuit,
		}
		_, choice, err := (&promptui.Select{
			Label: L.CliMenuPrompt,
			Items: choices,
		}).Run()
		handlePromptError(err)

		switch choice {
		case L.CliMenuHelp:
			flag.Usage()
			return
		case L.CliMenuQuit:
			return
		case L.CliMenuUpdate:
			if err := UpdateSelf(); err != nil {
				Log.Error(L.CliSelfUpdateFail, err)
				exitFailure()
			}
			exitSuccess()
		}

		*switches[SliceIndex(choices, choice)] = true
	}

	var err error
	var errSilent error
	if install {
		errSilent = PromptDiscord("patch", *locationFlag, *branchFlag).patch()
	} else if uninstall {
		errSilent = PromptDiscord("unpatch", *locationFlag, *branchFlag).unpatch()
	} else if update {
		Log.Info(L.CliDownloading)
		err := installLatestBuilds()
		Log.Info(L.CliDone)
		if err == nil {
			errSilent = PromptDiscord("repair", *locationFlag, *branchFlag).patch()
		}
	} else if installOpenAsar {
		discord := PromptDiscord("patch", *locationFlag, *branchFlag)
		if !discord.IsOpenAsar() {
			err = discord.InstallOpenAsar()
		} else {
			die(L.CliOpenAsarAlready)
		}
	} else if uninstallOpenAsar {
		discord := PromptDiscord("patch", *locationFlag, *branchFlag)
		if discord.IsOpenAsar() {
			err = discord.UninstallOpenAsar()
		} else {
			die(L.CliOpenAsarNotInstalled)
		}
	}

	if err != nil {
		Log.Error(err)
		exitFailure()
	}
	if errSilent != nil {
		exitFailure()
	}

	exitSuccess()
}

func exit(status int) {
	if runtime.GOOS == "windows" && IsDoubleClickRun() && interactive {
		fmt.Println(L.CliPressEnterToExit)
		var b byte
		_, _ = fmt.Scanf("%v", &b)
	}
	os.Exit(status)
}

func exitSuccess() {
	color.HiGreen(L.CliExitSuccess)
	exit(0)
}

func exitFailure() {
	color.HiRed(L.CliExitFailure)
	exit(1)
}

func handlePromptError(err error) {
	if errors.Is(err, promptui.ErrInterrupt) {
		exit(0)
	}

	Log.FatalIfErr(err)
}

func PromptDiscord(action, dir, branch string) *DiscordInstall {
	if branch == "auto" {
		for _, b := range []string{"stable", "canary", "ptb"} {
			for _, discord := range discords {
				install := discord.(*DiscordInstall)
				if install.branch == b {
					return install
				}
			}
		}
		die(L.CliNoDiscord)
	}

	if branch != "" {
		for _, discord := range discords {
			install := discord.(*DiscordInstall)
			if install.branch == branch {
				return install
			}
		}
		die(fmt.Sprintf(L.CliBranchNotFoundFmt, branch))
	}

	if dir != "" {
		if discord := ParseDiscord(dir, branch); discord != nil {
			return discord
		}

		if discord := ParseDiscordNew(dir, branch, strings.Contains(dir, "com.discordapp")); discord != nil {
			return discord
		}

		die(fmt.Sprintf(L.CliInvalidLocationFmt, dir))
	}

	items := SliceMap(discords, func(d any) string {
		install := d.(*DiscordInstall)
		//goland:noinspection GoDeprecation
		return fmt.Sprintf("%s - %s%s", strings.Title(install.branch), install.path, Ternary(install.isPatched, L.CliVencordInstalledTag, ""))
	})
	items = append(items, L.CliCustomLocationItem)

	promptLabel := ""
	switch action {
	case "patch":
		promptLabel = L.CliSelectPatch
	case "unpatch":
		promptLabel = L.CliSelectUnpatch
	default:
		promptLabel = L.CliSelectRepair
	}

	_, choice, err := (&promptui.Select{
		Label: promptLabel,
		Items: items,
	}).Run()
	handlePromptError(err)

	if choice != L.CliCustomLocationItem {
		return discords[SliceIndex(items, choice)].(*DiscordInstall)
	}

	for {
		custom, err := (&promptui.Prompt{
			Label: L.CliCustomLocationPrompt,
		}).Run()
		handlePromptError(err)

		if di := ParseDiscord(custom, ""); di != nil {
			return di
		}

		if di := ParseDiscordNew(custom, "", strings.Contains(custom, "com.discordapp")); di != nil {
			return di
		}

		Log.Error(L.CliInvalidDiscord)
	}
}

func InstallLatestBuilds() error {
	if IsDevInstall {
		return nil
	}

	return installLatestBuilds()
}

func HandleScuffedInstall() {
	fmt.Println(L.CliScuffedBody1)
	fmt.Println(L.CliScuffedBody2)
	fmt.Println(L.CliScuffedBody3)
}
