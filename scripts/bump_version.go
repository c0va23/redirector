//usr/bin/env go run "$0" "$@"; exit $?
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	scriptPath := os.Args[0]
	if len(os.Args) < 3 {
		panic(fmt.Sprintf("Run as:\n%s VERSION <level>", scriptPath))
	}
	versionFilePath := os.Args[1]
	level := os.Args[2]
	data, err := ioutil.ReadFile(versionFilePath)
	if nil != err {
		panic(err)
	}

	version := strings.TrimSpace(string(data))
	fmt.Printf("Old version: %s\n", version)

	versionParts := strings.Split(version, ".")
	if 3 != len(versionParts) {
		panic(fmt.Sprintf(`Invalid version: "%s"`, version))
	}

	major, majorErr := strconv.ParseUint(versionParts[0], 10, 32)
	if nil != majorErr {
		panic(majorErr)
	}
	minor, minorErr := strconv.ParseUint(versionParts[1], 10, 32)
	if nil != minorErr {
		panic(minorErr)
	}
	patch, patchErr := strconv.ParseUint(versionParts[2], 10, 32)
	if nil != patchErr {
		panic(patchErr)
	}

	switch level {
	case "major":
		{
			major++
			minor = 0
			patch = 0
			break
		}
	case "minor":
		{
			minor++
			patch = 0
			break
		}
	case "patch":
		{
			patch++
			break
		}
	default:
		panic(fmt.Sprintf("Unknown level: %s", level))
	}

	newVersion := fmt.Sprintf("%d.%d.%d\n", major, minor, patch)
	fmt.Printf("New version: %s", newVersion)

	if err := ioutil.WriteFile(versionFilePath, ([]byte)(newVersion), os.ModePerm); nil != err {
		panic(err)
	}
}
