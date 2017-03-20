// Copyright 2017 EcoSystem Software LLP

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
)

//packageName is the name of the package from which this is being called.
//Usage with flag: go:generate hardcodetemplates -p=admin
var packageName = flag.String("p", "", "Package name - must be set")

func main() {

	//Parse the flag for the package name
	flag.Parse()

	//Create a template holder with all .html files in the 'template' directory
	//relative to where this is run
	var templates *template.Template
	templates, _ = templates.ParseGlob("templates/*")

	//Variable for appending template contents
	allTemplates := ""

	//For each template...
	for _, t := range templates.Templates() {

		//Generate a template definition like {{ define x }} templateContents {{ end }}
		top := fmt.Sprintf(defineTop, t.Tree.Name)
		tpl := fmt.Sprint(t.Tree.Root)
		bottom := fmt.Sprint(defineBottom)

		//Append to the master template contents string
		allTemplates += (top + tpl + bottom)
	}

	//Wrap the master template string as a template literal
	allTemplates = fmt.Sprintf("`%s`", allTemplates)

	//Generate the final file
	h := fmt.Sprintf(fileContents, *packageName, allTemplates)

	//Write the file
	err := ioutil.WriteFile("templates.go", []byte(h), 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}

}

const (
	fileContents = `// Generated by jpincas/hardcodetemplates
// DONT edit

package %s

const baseTemplate = %s`

	defineTop = `{{ define "%s" }}`

	defineBottom = `{{ end }}`
)
