package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func build(app *tview.Application) *tview.Application {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("todo app\n<C-q> -> exit app")

	listField := tview.NewList().ShowSecondaryText(false)
	listField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		eventRune := event.Rune()
		if eventRune == rune('j') {
			listField.SetCurrentItem(min(listField.GetItemCount(), listField.GetCurrentItem()+1))
		} else if eventRune == rune('k') {
			listField.SetCurrentItem(max(0, listField.GetCurrentItem()-1))
		}
		return event
	})

	inputField := tview.NewInputField()

	inputField.SetDoneFunc(func(key tcell.Key) {
		text := inputField.GetText()
		if key == tcell.KeyEnter {
			if text != "" {
				listField.AddItem(text, "", rune('□'), nil)
				inputField.SetText("")
			}
		}
	})

	grid := tview.NewGrid().
		SetRows(3, 1, -1, 1).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(inputField, 1, 0, 1, 3, 0, 0, true).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(listField, 2, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Footer"), 3, 0, 1, 3, 0, 0, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		if key == tcell.KeyF1 {
			if inputField.HasFocus() {
				app.SetFocus(listField)
			} else {
				app.SetFocus(inputField)
			}
		} else if key == tcell.KeyF2 && listField.GetItemCount() > 0 && listField.HasFocus() {
			listField.RemoveItem(listField.GetCurrentItem())
		} else if key == tcell.KeyCtrlQ {
			app.Stop()
		} else {
		}
		return event
	})

	// SetTitleColor(tcell.ColorWhite)
	app.SetRoot(grid, true)
	// app.SetFocus(help)
	return app
}

func main() {
	app := tview.NewApplication()
	app = build(app)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

