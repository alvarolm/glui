package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var loadShowList func(packagelist []*showItem)

func init() {
	loadPackages()
}

const (
	list_pageName    = "list"
	package_pageName = "package"
)

func main() {
	tview.Styles.ContrastBackgroundColor = 0

	app := tview.NewApplication()
	app.EnableMouse(true)
	var form *tview.Form
	var table *tview.Table
	var input *tview.InputField

	input = tview.NewInputField().
		SetLabel("Filter:").
		SetText("").
		SetFieldWidth(0).
		SetAcceptanceFunc(nil).
		SetChangedFunc(nil)

	input.SetChangedFunc(func(text string) {
		loadShowList(filterShowList(text))
	})

	form = tview.NewForm().
		AddFormItem(input)
	form.SetBackgroundColor(0)
	form.Box.SetBorderPadding(1, 0, 0, 0)
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp, tcell.KeyDown:
			app.SetFocus(table)

			switch event.Key() {
			case tcell.KeyUp:
				row, _ := table.GetSelection()
				if row > 0 {
					table.Select(row-1, 0)
				}
			case tcell.KeyDown:
				row, _ := table.GetSelection()
				rc := table.GetRowCount()
				if rc > 0 && row < rc-1 {
					table.Select(row+1, 0)
				}
			}

		case tcell.KeyEnter:
			i := (table.GetCell(table.GetSelection()).Reference).(*showItem)
			jp := loadesPackages.list[i.index]
			updatePackageView(jp)
		}

		return event
	})

	table = tview.NewTable().
		SetBorders(false)
	table.Box.SetBackgroundColor(0)
	table.Box.SetBorderPadding(0, 0, 0, 0)
	table.SetSelectable(true, true)
	table.Select(0, 0)
	loadShowList = func(packageList []*showItem) {
		table.Clear()

		for i, p := range packageList {
			table.SetCell(i, 0, tview.NewTableCell(p.impPath).SetReference(p))
		}

		if len(packageList) > 0 {
			table.Select(0, 0)
		}
	}
	loadShowList(filterShowList(""))

	table.SetSelectedFunc(func(row int, column int) {
		i := (table.GetCell(row, column).Reference).(*showItem)
		jp := loadesPackages.list[i.index]
		updatePackageView(jp)
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()

		// do nothing if manipulating results
		switch key {
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyEnter, tcell.KeyTab:
			return event
		case tcell.KeyDelete, tcell.KeyBackspace, tcell.KeyBackspace2, tcell.KeyBacktab:
			app.SetFocus(form)
			return nil
		}

		// if pressed a printable charcater then focus to filter and copy it
		if key == tcell.KeyRune {
			app.SetFocus(form)
			input.SetText(input.GetText() + string(event.Rune()))
			return nil
		}

		return event
	})

	list_flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form, 3, 0, true).
		AddItem(table, 0, 1, false)

	list_flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if form.HasFocus() {
				app.SetFocus(table)
			} else {
				app.SetFocus(form)
			}
			return nil
		}

		return event
	})
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.Stop()
		}

		return event
	})

	pages := tview.NewPages()

	package_flex := composePackageView(app, pages)

	pages.AddPage(list_pageName, list_flex, true, true)
	pages.AddPage(package_pageName, package_flex, true, false)

	if err := app.SetRoot(pages, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}
}
