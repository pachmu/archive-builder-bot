package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

type ArchiveEditor interface {
	ChangeBuild(buildName string, build int) error
}

type archiveEditor struct {
	archivePath string
}

func GetNewArchiveEditor(archivePath string) ArchiveEditor {
	return &archiveEditor{
		archivePath: archivePath,
	}
}

func (a *archiveEditor) ChangeBuild(buildName string, build int) error {
	switch buildName {
	case "vpn":
		return a.changeVpnBuild(build)
	}

	return nil
}

func (a *archiveEditor) changeVpnBuild(build int) error {
	file, err := a.readMakeFile()
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`VPN_IMAGE_VERSION := .*\n`)
	file = re.ReplaceAllString(file, fmt.Sprintf("VPN_IMAGE_VERSION := %d\n", build))
	err = a.writeMakeFile(file)
	if err != nil {
		return err
	}
	return nil
}

func (a *archiveEditor) readMakeFile() (string, error) {
	b, err := ioutil.ReadFile(a.archivePath + "/makefile")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (a *archiveEditor) writeMakeFile(file string) error {
	err := ioutil.WriteFile(a.archivePath+"/makefile", []byte(file), 0644)
	if err != nil {
		return err
	}
	return nil
}
