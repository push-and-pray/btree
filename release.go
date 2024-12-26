//go:build !debug
// +build !debug

package main

func assert(_ func() bool, _ string) {
	// No-op in release mode
}
