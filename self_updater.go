/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
	"vencordinstaller/buildinfo"
)

var IsSelfOutdated = false
var SelfUpdateCheckDoneChan = make(chan bool, 1)

// LatestInstallerRelease holds the newest release found during the update
// check; UpdateSelf downloads its assets, so the installer always updates
// from yepyepos/Installer and never from the official English repository.
var LatestInstallerRelease *GithubRelease

func init() {
	//goland:noinspection GoBoolExpressions
	if buildinfo.InstallerTag == buildinfo.VersionUnknown {
		Log.Debug("Disabling self updater as this is not a release build")
		return
	}

	go DeleteOldExecutable()

	go func() {
		Log.Debug("Checking for Installer Updates...")

		releases, err := fetchGithubReleases(InstallerReleaseUrl, InstallerReleaseUrlFallback)
		if err != nil {
			Log.Warn(L.SelfCheckFail, err)
			SelfUpdateCheckDoneChan <- false
		} else {
			LatestInstallerRelease = selectInstallerRelease(releases)
			if LatestInstallerRelease == nil {
				Log.Warn(L.SelfCheckFail, L.SelfNoValidRelease)
				SelfUpdateCheckDoneChan <- false
			} else {
				// Only newer releases trigger the update prompt. A locally
				// newer (or incomparable) tag never does, so a release that
				// is up to date can never end up in an update loop.
				IsSelfOutdated = compareVersions(LatestInstallerRelease.TagName, buildinfo.InstallerTag) > 0
				Log.Debug("Latest installer release is", LatestInstallerRelease.TagName, "local is", buildinfo.InstallerTag, "Is self outdated?", IsSelfOutdated)
				SelfUpdateCheckDoneChan <- true
			}
		}
	}()
}

// selectInstallerRelease picks the release with the highest version tag from
// the yepyepos/Installer release list (pre-releases included, drafts are
// never returned by the API). API order is not relied upon.
func selectInstallerRelease(releases []GithubRelease) *GithubRelease {
	var best *GithubRelease
	for i := range releases {
		release := &releases[i]
		if !parseVersion(release.TagName).valid {
			continue
		}
		if best == nil || compareVersions(release.TagName, best.TagName) > 0 {
			best = release
		}
	}
	return best
}

func GetInstallerDownloadLink() string {
	// Browser fallback (macOS opens this in a browser because it cannot
	// replace its own bundle). Always points at the fork's releases.
	const BaseUrl = "https://github.com/yepyepos/Installer/releases/latest/download/"
	switch runtime.GOOS {
	case "windows":
		filename := Ternary(buildinfo.UiType == buildinfo.UiTypeCli, "VencordInstallerCli.exe", "VencordInstaller.exe")
		return BaseUrl + filename
	case "darwin":
		return BaseUrl + "VencordInstaller.dmg"
	case "linux":
		return BaseUrl + "VencordInstallerCli-linux"
	default:
		return ""
	}
}

// installerAssetName returns the release asset this build updates itself
// with. Asset names intentionally match the official installer so the fork
// releases stay drop-in compatible.
func installerAssetName() string {
	switch runtime.GOOS {
	case "windows":
		return Ternary(buildinfo.UiType == buildinfo.UiTypeCli, "VencordInstallerCli.exe", "VencordInstaller.exe")
	case "linux":
		return "VencordInstallerCli-linux"
	default:
		return ""
	}
}

// installerAssetUrl resolves the download URL of this build's update asset
// from the release the update check selected (yepyepos/Installer only).
func installerAssetUrl(release *GithubRelease) (string, error) {
	name := installerAssetName()
	if name == "" {
		return "", errors.New(L.SelfNoAssetErr)
	}
	for _, ass := range release.Assets {
		if ass.Name == name {
			return ass.DownloadURL, nil
		}
	}
	return "", fmt.Errorf(L.SelfNoAssetFmt, name)
}

func CanUpdateSelf() bool {
	//goland:noinspection GoBoolExpressions
	return IsSelfOutdated && runtime.GOOS != "darwin"
}

func UpdateSelf() error {
	if !CanUpdateSelf() {
		return errors.New(L.SelfUpToDateErr)
	}

	if LatestInstallerRelease == nil {
		return errors.New(L.SelfCheckFail + " " + L.SelfNoValidRelease)
	}

	url, err := installerAssetUrl(LatestInstallerRelease)
	if err != nil {
		return err
	}

	Log.Debug("Updating self from", url)

	ownExePath, err := os.Executable()
	if err != nil {
		return err
	}

	ownExeDir := path.Dir(ownExePath)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf(L.SelfDownloadStatusFmt, res.StatusCode, res.Status)
	}

	tmp, err := os.CreateTemp(ownExeDir, "VencordInstallerUpdate")
	if err != nil {
		return fmt.Errorf(L.SelfTempFailFmt, err)
	}
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}()
	if err = tmp.Chmod(0o755); err != nil {
		return fmt.Errorf(L.SelfTempFailFmt, err)
	}

	read, err := io.Copy(tmp, res.Body)
	if err != nil {
		return err
	}

	// Completeness check (same mechanism installLatestBuilds uses): never
	// replace a working installer with a truncated download. Full checksum
	// verification is a separate future enhancement.
	if expected := res.Header.Get("Content-Length"); expected != "" && expected != strconv.FormatInt(read, 10) {
		return fmt.Errorf(L.SelfLengthMismatchFmt, expected, strconv.FormatInt(read, 10))
	}

	if err = tmp.Close(); err != nil {
		return err
	}

	if err = os.Remove(ownExePath); err != nil {
		if err = os.Rename(ownExePath, ownExePath+".old"); err != nil {
			return fmt.Errorf(L.SelfRemoveFailFmt, err)
		}
	}

	if err = os.Rename(tmp.Name(), ownExePath); err != nil {
		// Restore the previous executable so a failed replace can never
		// leave the user without a working installer.
		if _, statErr := os.Stat(ownExePath + ".old"); statErr == nil {
			if restoreErr := os.Rename(ownExePath+".old", ownExePath); restoreErr != nil {
				Log.Error("Failed to restore previous executable:", restoreErr)
			} else {
				Log.Info(L.SelfRollbackOk)
			}
		}
		return fmt.Errorf(L.SelfReplaceFailFmt, err)
	}

	return nil
}

func DeleteOldExecutable() {
	ownExePath, err := os.Executable()
	if err != nil {
		return
	}

	for attempts := 0; attempts < 10; attempts += 1 {
		err = os.Remove(ownExePath + ".old")

		if err == nil || errors.Is(err, os.ErrNotExist) {
			break
		}

		Log.Warn("Failed to remove old executable. Retrying in 1 second.", err)
		time.Sleep(1 * time.Second)
	}
}

func RelaunchSelf() error {
	attr := new(os.ProcAttr)
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	var argv []string
	if len(os.Args) > 1 {
		argv = os.Args[1:]
	} else {
		argv = []string{}
	}

	Log.Debug("Restarting self with exe", os.Args[0], "and args", argv)

	proc, err := os.StartProcess(os.Args[0], argv, attr)
	if err != nil {
		return fmt.Errorf(L.SelfStartFailFmt, err)
	}

	if err = proc.Release(); err != nil {
		return fmt.Errorf(L.SelfReleaseProcFailFmt, err)
	}

	os.Exit(0)
	return nil
}
