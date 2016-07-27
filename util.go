package main

import (
	"fmt"
	"github.com/fatih/color"
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
