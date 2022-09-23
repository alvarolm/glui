package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var updatePackageView func(jp *jsonPackage)

func composePackageView(app *tview.Application, pages *tview.Pages) tview.Primitive {
	textView := tview.NewTextView().
		SetDynamicColors(true)
	textView.SetBackgroundColor(0)
	textView.SetBorder(true)

	table := tview.NewTable().
		SetBorders(false)
	table.Box.SetBackgroundColor(0)
	table.Box.SetBorderPadding(0, 0, 0, 0)
	table.SetSelectable(true, true)

	table.SetSelectedFunc(func(row int, column int) {
		c := table.GetCell(row, column)
		jp := (table.GetCell(row, column).Reference).(*jsonPackage)

		path := filepath.Join(jp.Dir, c.Text)

		editorCmd := os.Getenv("EDITOR")

		exitCode := 0
		var err error

		if len(editorCmd) == 0 {
			err = errors.New("missing 'EDITOR' environment variable")
			exitCode = 1
			goto exit
		}

		if editorCmd, err = exec.LookPath(editorCmd); err != nil {
			exitCode = 1
			goto exit
		}

		err = exec.Command(editorCmd, path).Run()

	exit:

		if err != nil {
			app.Stop()
			fmt.Println(err.Error())
			exitCode = 1
		}
		os.Exit(exitCode)
	})

	updatePackageView = func(jp *jsonPackage) {
		textView.SetText(jp.headerText())
		textView.SetTitle(" " + jp.ImportPath + " ")

		table.Clear()
		for i, goFilePath := range jp.GoFiles {
			table.SetCell(i, 0, tview.NewTableCell(goFilePath).SetReference(jp))
		}

		pages.SwitchToPage(package_pageName)
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textView, 0, 2, false).
		AddItem(table, 0, 3, true)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBackspace, tcell.KeyBackspace2, tcell.KeyBacktab:
			pages.SwitchToPage(list_pageName)
		case tcell.KeyTab:
			if textView.HasFocus() {
				app.SetFocus(table)
			} else {
				app.SetFocus(textView)
			}
			return nil
		}

		return event
	})

	return flex
}
