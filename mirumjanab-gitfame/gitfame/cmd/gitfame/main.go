//go:build !solution

package main

import (
	"os"
	
	"gitlab.com/slon/shad-go/gitfame/pkg/analyse"
	"gitlab.com/slon/shad-go/gitfame/pkg/formatter"
	"gitlab.com/slon/shad-go/gitfame/pkg/git"
	"gitlab.com/slon/shad-go/gitfame/pkg/options"
)

func main() {
	opt := options.GetOptions(os.Args[1:])
	g := git.NewGit(opt)

	data := analyse.New()
	authorSlice := data.AnalyzeGitFiles(g)
	authorSlice.Sort(opt.OrderBy())

	printer := formatter.New(opt.OutFormat(), os.Stdout)
	printer.Print(authorSlice)
}