// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hello is a trivial package for gomobile bind example.
package hello

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/mobile/asset"
	yaml "gopkg.in/yaml.v2"
)

type androidEnv struct {
	Greetstring string `yaml:displaystring`
}

func (ae *androidEnv) read() *androidEnv {
	yamlFile, err := readAsset("example.yaml")
	if err == nil {
		//		log.Printf("read android config")

		err = yaml.Unmarshal(yamlFile, ae)
		if err == nil {
			return ae
		}
	}
	return ae

}

func readAsset(name string) ([]byte, error) {
	f, errOpen := asset.Open(name)
	//var f *os.File
	//var errOpen error

	if errOpen != nil {
		return nil, errOpen
	}
	defer f.Close()
	buf, errRead := ioutil.ReadAll(f)
	if errRead != nil {
		return nil, errRead
	}
	return buf, nil
}

func Greetings(name string) string {
	var ae androidEnv
	ae.read()
	if len(ae.Greetstring) != 0 {
		return ae.Greetstring
	}
	return fmt.Sprintf("Hello, %s. Didn't find anything!", name)
}
