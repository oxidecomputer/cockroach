// Copyright 2019 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package faker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

func TestFaker(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	f := NewFaker()

	names := []string{
		`Daniel Nixon`,
		`Whitney Jimenez`,
		`Brandon Carr`,
	}
	for i, name := range names {
		assert.Equal(t, name, f.Name(rng), `testcase %d`, i)
	}

	streetAddresses := []string{
		`8339 Gary Burgs Apt. 6`,
		`67941 Lawrence Station Suite 29`,
		`29657 Ware Haven`,
	}
	for i, streetAddress := range streetAddresses {
		assert.Equal(t, streetAddress, f.StreetAddress(rng), `testcase %d`, i)
	}

	wordsets := [][]string{
		{`until`},
		{`figure`, `after`},
		{`suddenly`, `heavy`, `time`},
	}
	for i, words := range wordsets {
		assert.Equal(t, words, f.Words(rng, i+1), `testcase %d`, i)
	}

	sentences := [][]string{
		{
			`Especially claim rather town.`,
		},
		{
			`Bag other follow play agreement develop sing.`,
			`Deal great national yard various mouth.`,
		},
		{
			`During talk direction set clear direction.`,
			`Realize once thus administration.`,
			`Glass industry drop prove large age any.`,
		},
	}
	for i, sentence := range sentences {
		assert.Equal(t, sentence, f.Sentences(rng, i+1), `testcase %d`, i)
	}

	paragraphs := []string{
		`Natural purpose member institution picture address. Goal use produce drive worry process ` +
			`beautiful somebody.`,
		`Stuff home capital international. Consumer message bed story here. Contain real expert ` +
			`institution. Against ever seek become put respond. Maybe recently entire history always ` +
			`former.`,
		`Court three author ground. College walk inside coach system career newspaper. However ` +
			`health me community mission. Senior evidence form size true general compare. Teacher look ` +
			`left else.`,
	}
	for i, paragraph := range paragraphs {
		assert.Equal(t, paragraph, f.Paragraph(rng), `testcase %d`, i)
	}
}

func TestFirstToUpper(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{`foobar`, `Foobar`},
		{`fooBar`, `FooBar`},
		{`foo bar`, `Foo bar`},
		{`foo Bar`, `Foo Bar`},
		{`κόσμε`, `Κόσμε`},
	}
	for i, test := range tests {
		assert.Equal(t, test.expected, firstToUpper(test.input), `testcase %d`, i)
	}
}
