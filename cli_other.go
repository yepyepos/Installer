//go:build !windows

package main

func enableUTF8Console() {}

func IsDoubleClickRun() bool {
	return false
}
