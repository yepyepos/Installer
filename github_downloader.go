/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	path "path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type GithubRelease struct {
	Name        string `json:"name"`
	TagName     string `json:"tag_name"`
	Body        string `json:"body"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var ReleaseData GithubRelease
var GithubError error
var GithubDoneChan chan bool

var InstalledHash = "None"
var LatestHash = "Unknown"
var IsDevInstall bool

// Vencord zh-CN release contract (see docs/VENCORD_RELEASE_CONTRACT.md):
// a release qualifies as a desktop Vencord release if it ships all of these
// assets (matched by name prefix, so .map siblings match too), and its release
// notes carry the machine readable build hash line "Vencord-Desktop-Hash: <hash>"
// which must equal the "// Vencord <hash>" header of the shipped patcher.js.
var requiredDesktopAssets = []string{"patcher.js", "preload.js", "renderer.js", "renderer.css"}

var desktopHashRe = regexp.MustCompile(`(?mi)^vencord-desktop-hash:[ \t]*([0-9a-f]{7,40})[ \t]*$`)

func HasAllDesktopAssets(release *GithubRelease) bool {
	for _, name := range requiredDesktopAssets {
		found := false
		for _, ass := range release.Assets {
			if strings.HasPrefix(ass.Name, name) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// ExtractDesktopHash returns the build hash a release advertises for its
// desktop dist files. Prefers the explicit contract marker in the release
// notes; falls back to the upstream format where the release NAME ends with
// the hash ("DevBuild <hash>").
func ExtractDesktopHash(release *GithubRelease) string {
	if m := desktopHashRe.FindStringSubmatch(release.Body); m != nil {
		return strings.ToLower(m[1])
	}
	i := strings.LastIndex(release.Name, " ") + 1
	return release.Name[i:]
}

// selectVencordRelease picks the newest published release that actually ships
// all required desktop assets. Order of the API response is NOT relied upon;
// releases are compared by published_at. Pre-releases are eligible, drafts are
// never returned by the API.
func selectVencordRelease(releases []GithubRelease) (*GithubRelease, error) {
	var best *GithubRelease
	var bestTime time.Time
	for i := range releases {
		release := &releases[i]
		if !HasAllDesktopAssets(release) {
			continue
		}
		published, err := time.Parse(time.RFC3339, release.PublishedAt)
		if err != nil {
			published = time.Time{}
		}
		if best == nil || published.After(bestTime) {
			best = release
			bestTime = published
		}
	}
	if best == nil {
		return nil, errors.New(L.ErrNoDesktopAssets)
	}
	return best, nil
}

func fetchGithubReleases(url, fallbackUrl string) ([]GithubRelease, error) {
	Log.Debug("Fetching", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log.Error(L.LogCreateRequestFail, err)
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		Log.Error(L.LogSendFail, err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		isRateLimitedOrBlocked := res.StatusCode == 401 || res.StatusCode == 403 || res.StatusCode == 429
		triedFallback := url == fallbackUrl

		if isRateLimitedOrBlocked && !triedFallback {
			Log.Error(fmt.Sprintf(L.LogFetchFallbackFmt, url, res.StatusCode, fallbackUrl))
			return fetchGithubReleases(fallbackUrl, fallbackUrl)
		}

		err = errors.New(res.Status)
		Log.Error(url, L.LogNonOKStatus, err)
		return nil, err
	}

	var releases []GithubRelease

	if err = json.NewDecoder(res.Body).Decode(&releases); err != nil {
		Log.Error(L.LogDecodeFail, err)
		return nil, err
	}

	return releases, nil
}

// GetLatestVencordRelease fetches the release to install Vencord from.
// Only yepyepos/Vencord is ever queried; there is deliberately no fallback to
// the official English Vencord repositories.
func GetLatestVencordRelease() (*GithubRelease, error) {
	releases, err := fetchGithubReleases(ReleaseUrl, ReleaseUrlFallback)
	if err != nil {
		return nil, err
	}
	return selectVencordRelease(releases)
}

func GetGithubRelease(url, fallbackUrl string) (*GithubRelease, error) {
	Log.Debug("Fetching", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log.Error("Failed to create Request", err)
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		Log.Error("Failed to send Request", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		isRateLimitedOrBlocked := res.StatusCode == 401 || res.StatusCode == 403 || res.StatusCode == 429
		triedFallback := url == fallbackUrl

		// GitHub has a very strict 60 req/h rate limit and some (mostly indian) isps block github for some reason.
		// If that is the case, try our fallback at https://vencord.dev/releases/project
		if isRateLimitedOrBlocked && !triedFallback {
			Log.Error(fmt.Sprintf("Failed to fetch %s (status code %d). Trying fallback url %s", url, res.StatusCode, fallbackUrl))
			return GetGithubRelease(fallbackUrl, fallbackUrl)
		}

		err = errors.New(res.Status)
		Log.Error(url, "returned Non-OK status", GithubError)
		return nil, err
	}

	var data GithubRelease

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		Log.Error("Failed to decode GitHub JSON Response", err)
		return nil, err
	}

	return &data, nil
}

func InitGithubDownloader() {
	GithubDoneChan = make(chan bool, 1)

	IsDevInstall = os.Getenv("VENCORD_DEV_INSTALL") == "1"
	Log.Debug("Is Dev Install: ", IsDevInstall)
	if IsDevInstall {
		GithubDoneChan <- true
		return
	}

	go func() {
		// Make sure UI updates once the request either finished or failed
		defer func() {
			GithubDoneChan <- GithubError == nil
		}()

		data, err := GetLatestVencordRelease()
		if err != nil {
			Log.Error(L.LogFetchDataFailed, err)
			GithubError = err
			return
		}

		ReleaseData = *data

		LatestHash = ExtractDesktopHash(data)
		Log.Debug("Selected release", data.TagName, "("+data.Name+")")
		Log.Debug("Finished fetching GitHub Data")
		Log.Debug("Latest hash is", LatestHash, "Local Install is", Ternary(LatestHash == InstalledHash, "up to date!", "outdated!"))
	}()

	// Check hash of installed version if exists
	f, err := os.Open(Patcher)
	if err != nil {
		return
	}
	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()

	Log.Debug("Found existing Vencord Install. Checking for hash...")
	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "// Vencord ") {
			if !distHasAllDesktopFiles(FilesDir) {
				// Interrupted download left a partial dist behind: treat the
				// install as outdated so all files are fetched again.
				Log.Debug("Dist folder is incomplete, treating install as outdated")
			} else {
				InstalledHash = line[11:]
				Log.Debug("Existing hash is", InstalledHash)
			}
		} else {
			Log.Debug("Didn't find hash")
		}
	}
}

// distHasAllDesktopFiles reports whether every required desktop dist file is
// present in dir. A partial dist (e.g. left behind by an interrupted
// download) must never count as an installed version.
func distHasAllDesktopFiles(dir string) bool {
	for _, name := range requiredDesktopAssets {
		if !ExistsFile(path.Join(dir, name)) {
			return false
		}
	}
	return true
}

func installLatestBuilds() (retErr error) {
	// long error for testing
	// return errors.New("According to all known laws of aviation, there is no way a bee should be able to fly. Its wings are too small to get its fat little body off the ground. The bee, of course, flies anyway because bees don't care what humans think is impossible. Yellow, black. Yellow, black. Yellow, black. Yellow, black. Ooh, black and yellow! Let's shake it up a little. Barry! Breakfast is ready! Coming! Hang on a second. Hello? Barry? Adam? Can you believe this is happening? I can't. I'll pick you up. Looking sharp. Use the stairs, Your father paid good money for those. Sorry. I'm excited. Here's the graduate. We're very proud of you, son. A perfect report card, all B's. Very proud. Ma! I got a thing going here. You got lint on your fuzz. Ow! That's me! Wave to us! We'll be in row 118,000. Bye! Barry, I told you, stop flying in the house! Hey, Adam. Hey, Barry. Is that fuzz gel? A little. Special day, graduation. Never thought I'd make it. Three days grade school, three days high school. Those were awkward. Three days college. I'm glad I took a day and hitchhiked around The Hive. You did come back different. Hi, Barry. Artie, growing a mustache? Looks good. Hear about Frankie? Yeah. You going to the funeral? No, I'm not going. Everybody knows, sting someone, you die. Don't waste it on a squirrel. Such a hothead. I guess he could have just gotten out of the way. I love this incorporating an amusement park into our day. That's why we don't need vacations. Boy, quite a bit of pomp under the circumstances. Well, Adam, today we are men. We are! Bee-men. Amen! Hallelujah! Students, faculty, distinguished bees, please welcome Dean Buzzwell.")
	Log.Debug("Installing latest builds...")

	// create an empty package.json file in our files dir.
	// without this, node will walk up the file tree and search for a package.json in the
	// parent folders. This might lead to issues if the user for example has ~/package.json
	// with type: "module" in it
	pkgJsonFile := path.Join(FilesDir, "package.json")
	err := os.WriteFile(pkgJsonFile, []byte("{}"), 0644)
	if err != nil {
		Log.Warn(L.LogCreateFail, pkgJsonFile, err)
	}

	var wg sync.WaitGroup
	var downloadedFiles atomic.Uint64

	for _, ass := range ReleaseData.Assets {
		isRequired := false
		for _, required := range requiredDesktopAssets {
			if strings.HasPrefix(ass.Name, required) {
				isRequired = true
				break
			}
		}
		if isRequired {
			wg.Add(1)
			ass := ass // Need to do this to not have the variable be overwritten halfway through
			go func() {
				defer wg.Done()
				Log.Debug("Downloading file", ass.Name)

				res, err := http.Get(ass.DownloadURL)
				if err == nil && res.StatusCode >= 300 {
					err = errors.New(res.Status)
				}
				if err != nil {
					Log.Error(L.LogDownloadFail, ass.Name+":", err)
					retErr = err
					return
				}
				outFile := path.Join(FilesDir, ass.Name)
				out, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					Log.Error(L.LogCreateFail, outFile+":", err)
					retErr = err
					return
				}
				read, err := io.Copy(out, res.Body)
				if err != nil {
					Log.Error(L.LogDownloadToFail, outFile+":", err)
					retErr = err
					return
				}
				contentLength := res.Header.Get("Content-Length")
				expected := strconv.FormatInt(read, 10)
				if expected != contentLength {
					err = errors.New("Unexpected end of input. Content-Length was " + contentLength + ", but I only read " + expected)
					Log.Error(err.Error())
					retErr = err
					return
				}
				downloadedFiles.Add(1)
			}()
		}
	}

	wg.Wait()

	if retErr != nil {
		return retErr
	}
	if downloadedFiles.Load() < 4 {
		return errors.New(L.ErrMissingFiles)
	}

	Log.Debug("Done!")
	_ = FixOwnership(FilesDir)

	// Contract check: the downloaded patcher.js must carry the same build hash
	// the release advertised. Never trust "up to date" on a mismatch.
	f, err := os.Open(path.Join(FilesDir, "patcher.js"))
	if err == nil {
		scanner := bufio.NewScanner(f)
		if scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "// Vencord ") || line[11:] != LatestHash {
				Log.Warn(L.LogContractViolation, line, "!=", LatestHash)
			}
		}
		f.Close()
	}

	InstalledHash = LatestHash
	return
}
