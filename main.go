//go:generate patch vendor/github.com/clyphub/munkres/munkres.go munkres.patch
package main

import "github.com/gokapaya/cshelper/cmd"

func main() {
	cmd.Execute()
}
