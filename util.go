package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/jmoiron/jsonq"
	"strings"
	"unicode"
)

// -----------------------------------------------------------------------------

type Logger struct{}

var log Logger

var errorChar = color.RedString(" ✘ ")
var warnChar = color.YellowString(" ! ")
var successChar = color.GreenString(" ✓ ")

func (l Logger) Error(parts ...interface{}) {
	fmt.Println(append([]interface{}{errorChar}, parts...)...)
}

func (l Logger) Warn(parts ...interface{}) {
	fmt.Println(append([]interface{}{warnChar}, parts...)...)
}

func (l Logger) Success(parts ...interface{}) {
	fmt.Println(append([]interface{}{successChar}, parts...)...)
}

// -----------------------------------------------------------------------------

func getJQ(raw string) *jsonq.JsonQuery {
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(raw))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	return jq
}

// -----------------------------------------------------------------------------

func stripNonAlphanumeric(r rune) rune {
	if !unicode.IsLetter(r) {
		return -1
	}
	return unicode.ToLower(r)
}

func fuzzy(str1 string, str2 string) bool {
	str1 = strings.Map(stripNonAlphanumeric, str1)
	str2 = strings.Map(stripNonAlphanumeric, str2)

	return strings.Contains(str2, str1)
}
