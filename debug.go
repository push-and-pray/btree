//go:build debug
// +build debug

package main

func assert(condition func() bool, msg string) {
	if !condition() {
		panic("DEBUG ASSERTION FAILED: " + msg)
	}
}
