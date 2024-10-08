package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func build(app *tview.Application) *tview.Application {
	fotter := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	fotter.SetText("hoge")

	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("todo app\n<C-q> -> exit app")

	listField := tview.NewList().ShowSecondaryText(false)
	doneListField := tview.NewList().ShowSecondaryText(false)
	inputFieldExplanation := "F1 -> move TODO area"
	listFieldExplanation := "F1 -> move input area\td -> move item to DONE\tx -> remove item\tAlt + → -> move DONE area"
	doneListFiekdExplanation := "F1 -> move input area\to -> move item to TODO\tx -> remove item\tAlt + ← -> move TODO area"
	fotter.SetText(inputFieldExplanation)

	listField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		eventRune := event.Rune()
		eventKey := event.Key()
		if eventRune == rune('j') {
			listField.SetCurrentItem(min(listField.GetItemCount(), listField.GetCurrentItem()+1))
		} else if eventRune == rune('k') {
			listField.SetCurrentItem(max(0, listField.GetCurrentItem()-1))
		} else if eventRune == rune('d') && listField.GetItemCount() > 0 {
			mainText, _ := listField.GetItemText(listField.GetCurrentItem())
			doneListField.InsertItem(0, mainText, "", 0, nil)
			listField.RemoveItem(listField.GetCurrentItem())
			doneListField.SetCurrentItem(0)
		} else if eventRune == rune('x') && listField.GetItemCount() > 0 {
			listField.RemoveItem(listField.GetCurrentItem())
		}
		if (event.Modifiers()&tcell.ModAlt == 0) && eventKey == tcell.KeyRight {
			app.SetFocus(doneListField)
			fotter.SetText(doneListFiekdExplanation)
		}
		return event
	})

	doneListField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		eventRune := event.Rune()
		eventKey := event.Key()
		if eventRune == rune('j') {
			doneListField.SetCurrentItem(min(doneListField.GetItemCount(), doneListField.GetCurrentItem()+1))
		} else if eventRune == rune('k') {
			doneListField.SetCurrentItem(max(0, doneListField.GetCurrentItem()-1))
		} else if eventRune == rune('x') && doneListField.GetItemCount() > 0 {
			doneListField.RemoveItem(doneListField.GetCurrentItem())
		} else if eventRune == rune('o') && doneListField.GetItemCount() > 0 {
			mainText, _ := doneListField.GetItemText(doneListField.GetCurrentItem())
			listField.InsertItem(0, mainText, "", 0, nil)
			doneListField.RemoveItem(doneListField.GetCurrentItem())
			listField.SetCurrentItem(0)
		}
		if (event.Modifiers()&tcell.ModAlt == 0) && eventKey == tcell.KeyLeft {
			app.SetFocus(listField)
			fotter.SetText(listFieldExplanation)
		}
		return event
	})

	inputField := tview.NewInputField()

	inputField.SetDoneFunc(func(key tcell.Key) {
		text := inputField.GetText()
		if key == tcell.KeyEnter {
			if text != "" {
				listField.AddItem(text, "", 0, nil)
				inputField.SetText("")
			}
		}
	})

	innerGrid := tview.NewGrid().
		SetRows(-1).
		SetColumns(30, -1, -1).
		AddItem(listField, 0, 0, 1, 3, 0, 0, false).
		AddItem(doneListField, 0, 3, 1, 3, 0, 0, false)

	grid := tview.NewGrid().
		SetRows(3, 1, -1, 1).
		SetColumns(30, -1, -1).
		SetBorders(true).
		AddItem(inputField, 1, 0, 1, 3, 0, 0, true).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(innerGrid, 2, 0, 1, 3, 0, 0, false).
		AddItem(fotter, 3, 0, 1, 3, 0, 0, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		if key == tcell.KeyF1 {
			if inputField.HasFocus() {
				app.SetFocus(listField)
				fotter.SetText(listFieldExplanation)
			} else {
				app.SetFocus(inputField)
				fotter.SetText(inputFieldExplanation)
			}
		} else if key == tcell.KeyCtrlQ {
			app.Stop()
		} else {
		}
		return event
	})

	listField.SetBorder(true).SetTitle("TODO")
	doneListField.SetBorder(true).SetTitle("DONE")
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

