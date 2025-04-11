package main

import (
	"fmt"
	//"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

var (
		version = "1.0.0" 
)

func showBanner() {

	logo := `
     ██╗███████╗███╗   ███╗ ██████╗ ███╗   ██╗
     ██║██╔════╝████╗ ████║██╔═══██╗████╗  ██║
     ██║██████╗ ██╔████╔██║██║   ██║██╔██╗ ██║
██   ██║╚════██╗██║╚██╔╝██║██║   ██║██║╚██╗██║
╚█████╔╝███████║██║ ╚═╝ ██║╚██████╔╝██║ ╚████║
 ╚════╝ ╚══════╝╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝

    		JSMON/ Javascript Monitor
`
    fmt.Println(logo)
    // boldCyan := color.New(color.FgHiCyan, color.Bold)
    // banner := figure.NewFigure("JSMON", "", true)
    // boldCyan.Println(banner.String())
    // color.New(color.FgHiBlue).Println("\t\tjsmon.sh")
}

func displayVersion() {
	
	color.New(color.FgHiYellow).Printf("[INF] Current JSMON version v%s (latest)\n", version)
}

