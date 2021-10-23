// BSD 3-Clause License
//
// Copyright (c) 2019, Michael McLoughlin
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package ext

import (
	"regexp"
	"strings"
	"unicode"
)

// none is the zero rune.
var none rune

// Fixed data structures required for formula processing.
var (
	// Symbol replacer.
	replacer *strings.Replacer

	// Regular expressions for super/subscripts.
	supregexp *regexp.Regexp
	subregexp *regexp.Regexp

	// Rune replacement maps.
	super = map[rune]rune{}
	sub   = map[rune]rune{}
)

// casedChar represents a character with its superscript and subscript variants.
type casedChar struct {
	Char  rune
	Super rune
	Sub   rune
}

// chars is the table of super/subscriptable characters.
var chars = []casedChar{
	{'0', '\u2070', '\u2080'},
	{'1', '\u00B9', '\u2081'},
	{'2', '\u00B2', '\u2082'},
	{'3', '\u00B3', '\u2083'},
	{'4', '\u2074', '\u2084'},
	{'5', '\u2075', '\u2085'},
	{'6', '\u2076', '\u2086'},
	{'7', '\u2077', '\u2087'},
	{'8', '\u2078', '\u2088'},
	{'9', '\u2079', '\u2089'},
	{'a', '\u1d43', '\u2090'},
	{'b', '\u1d47', none},
	{'c', '\u1d9c', none},
	{'d', '\u1d48', none},
	{'e', '\u1d49', '\u2091'},
	{'f', '\u1da0', none},
	{'g', '\u1d4d', none},
	{'h', '\u02b0', '\u2095'},
	{'i', '\u2071', '\u1d62'},
	{'j', '\u02b2', '\u2c7c'},
	{'k', '\u1d4f', '\u2096'},
	{'l', '\u02e1', '\u2097'},
	{'m', '\u1d50', '\u2098'},
	{'n', '\u207f', '\u2099'},
	{'o', '\u1d52', '\u2092'},
	{'p', '\u1d56', '\u209a'},
	{'q', none, none},
	{'r', '\u02b3', '\u1d63'},
	{'s', '\u02e2', '\u209b'},
	{'t', '\u1d57', '\u209c'},
	{'u', '\u1d58', '\u1d64'},
	{'v', '\u1d5b', '\u1d65'},
	{'w', '\u02b7', none},
	{'x', '\u02e3', '\u2093'},
	{'y', '\u02b8', none},
	{'z', none, none},
	{'A', '\u1d2c', none},
	{'B', '\u1d2e', none},
	{'C', none, none},
	{'D', '\u1d30', none},
	{'E', '\u1d31', none},
	{'F', none, none},
	{'G', '\u1d33', none},
	{'H', '\u1d34', none},
	{'I', '\u1d35', none},
	{'J', '\u1d36', none},
	{'K', '\u1d37', none},
	{'L', '\u1d38', none},
	{'M', '\u1d39', none},
	{'N', '\u1d3a', none},
	{'O', '\u1d3c', none},
	{'P', '\u1d3e', none},
	{'Q', none, none},
	{'R', '\u1d3f', none},
	{'S', none, none},
	{'T', '\u1d40', none},
	{'U', '\u1d41', none},
	{'V', '\u2c7d', none},
	{'W', '\u1d42', none},
	{'X', none, none},
	{'Y', none, none},
	{'Z', none, none},
	{'+', '\u207A', '\u208A'},
	{'-', '\u207B', '\u208B'},
	{'=', '\u207C', '\u208C'},
	{'(', '\u207D', '\u208D'},
	{')', '\u207E', '\u208E'},
}

// charclass builds a regular expression character class from a list of runes.
func charclass(runes []rune) string {
	return strings.ReplaceAll("["+string(runes)+"]", "-", `\-`)
}

// subsupreplacer builds a replacement function that applies the repl rune map
// to a matched super/subscript.
func subsupreplacer(repl map[rune]rune) func(string) string {
	return func(s string) string {
		var runes []rune
		for i, r := range s {
			if i == 0 || unicode.IsSpace(r) {
				runes = append(runes, r)
			} else if repl[r] != none {
				runes = append(runes, repl[r])
			}
		}
		return string(runes)
	}
}

func init() {
	// Build symbol replacer.
	var oldnew []string
	for symbol, r := range symbols {
		oldnew = append(oldnew, symbol, string([]rune{r}))
	}
	replacer = strings.NewReplacer(oldnew...)

	// Build super/subscript character classes and replacement maps.
	var superclass, subclass []rune
	for _, char := range chars {
		if char.Super != none {
			superclass = append(superclass, char.Char)
			super[char.Char] = char.Super
		}
		if char.Sub != none {
			subclass = append(subclass, char.Char)
			sub[char.Char] = char.Sub
		}
	}

	// Build regular expressions.
	supregexp = regexp.MustCompile(`(\b[A-Za-z0-9]|[)\pL\pS^A-Za-z0-9])\^(\d+|\{` + charclass(superclass) + `+\}|` + charclass(superclass) + `\s)`)
	subregexp = regexp.MustCompile(`(\b[A-Za-z]|\pS|\p{Greek})_(\d+\b|\{` + charclass(subclass) + `+\})`)
}

// formula processes a formula in s, writing the result to w.
func formula(s string) string {
	// Replace symbols.
	s = replacer.Replace(s)

	// Replace superscripts.
	s = supregexp.ReplaceAllStringFunc(s, subsupreplacer(super))

	// Replace subscripts.
	s = subregexp.ReplaceAllStringFunc(s, subsupreplacer(sub))

	return s
}
