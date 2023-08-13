package main

import (
	"github.com/pterm/pterm"
)

type RPC struct {
	Method  string
	Payload string
}

func replRun() RPC {
	primary := pterm.NewStyle(pterm.FgLightCyan, pterm.Bold)
	secondary := pterm.NewStyle(pterm.FgLightGreen, pterm.Italic)

	method, _ := pterm.DefaultInteractiveTextInput.
		WithMultiLine(false).
		WithTextStyle(primary).
		Show("method  ")
	payload, _ := pterm.DefaultInteractiveTextInput.
		WithMultiLine(false).
		WithTextStyle(secondary).
		Show("payload ")

	return RPC{
		Method:  method,
		Payload: payload,
	}
}
