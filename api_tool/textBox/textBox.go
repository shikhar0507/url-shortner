package textBox

import (
	"fmt"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Body struct {
	Name  string
	Id    int
	Badge string
}

func TextBox(app *tview.Application) *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetBorder(true).SetTitle("Body")
	savedText := make([]string, 1, 1)

	position := -1
	//hasTextChanged := false
	hlPos := make([]string, 1)

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		//	hasTextChanged = false
		maxLength := len(savedText)
		switch event.Key() {
		case tcell.KeyCtrlB:
			clipboard.WriteAll(textView.GetText(true))
			break
		case tcell.KeyCtrlA:
			hlPos = nil
			for i := 0; i < maxLength; i++ {
				hlPos = append(hlPos, strconv.Itoa(i))
			}
			textView.Highlight(hlPos...)
			break
		}

		switch event.Name() {

		case "Backspace2":
			//	hasTextChanged = true

			if len(textView.GetHighlights()) > 1 {
				savedText = nil
				position = 0
				hlPos = nil
			} else {
				decrementPosition(&position, 1)
				if position == 0 {
					savedText = nil
				} else {
					savedText = append(savedText[:position], savedText[position+1:]...)
				}
			}
			writeText(savedText, textView)
			setHightLight(position, textView)
			break
		case "Left":
			hlPos = nil
			decrementPosition(&position, 1)
			setHightLight(position, textView)
			break
		case "Right":
			hlPos = nil
			incrementPosition(&position, 1, maxLength)
			setHightLight(position, textView)
			break
		case "Ctrl+Left":
			decrementPosition(&position, 8)
			setHightLight(position, textView)
			break
		case "Ctrl+Right":
			incrementPosition(&position, 8, maxLength)
			setHightLight(position, textView)
			break
		case "Shift+Right":
			indexof := -1

			for i := 0; i < len(hlPos); i++ {
				if hlPos[i] == strconv.Itoa(position) {
					indexof = i
				}
			}

			if indexof != -1 {
				hlPos = append(hlPos[:indexof], hlPos[indexof+1:]...)
			} else {
				hlPos = append(hlPos, strconv.Itoa(position))
			}
			textView.Highlight(hlPos...)
			incrementPosition(&position, 1, maxLength)

			break
		case "Shift+Left":

			decrementPosition(&position, 1)
			hlPos = append(hlPos, strconv.Itoa(position))
			textView.Highlight(hlPos...)

			break
		case "Up":
			hlPos = nil
			if position == 0 {
				break
			}
			for i := position - 1; i >= 0; i-- {
				if savedText[i] == "," {
					decrementPosition(&position, position-i)
					break
				}

			}
			setHightLight(position, textView)
			break
		case "Down":
			hlPos = nil

			if position == len(savedText) {
				break
			}

			for i := position + 1; i < len(savedText); i++ {
				if savedText[i] == "," {
					//fmt.Println(i)
					incrementPosition(&position, i-position, maxLength)
					break
				}
			}
			setHightLight(position, textView)
			break
		default:
			if rune(event.Rune()) >= 32 && rune(event.Rune()) <= 126 || event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyTAB {
				join := make([]string, 1)
				if event.Name() == "Enter" {
					join[0] = fmt.Sprintf("\n")
				} else if event.Name() == "Tab" {
					join[0] = "\t"
				} else {
					join[0] = string(event.Rune())
				}

				incrementPosition(&position, 1, maxLength)
				join = append(join, savedText[position:]...)
				savedText = append(savedText[:position], join...)
				writeText(savedText, textView)
				setHightLight(position, textView)

			}

			break
		}

		return event
	})

	textView.SetChangedFunc(func() {

		app.Draw()
	})
	textView.SetBorder(true)
	return textView
}

func writeText(savedText []string, textView *tview.TextView) {
	s := ""
	for i, value := range savedText {
		s += fmt.Sprintf(`["%d"]%s[""]`, i, string(value))

	}
	textView.SetText(s)
}

func setHightLight(position int, textView *tview.TextView) {

	pos := strconv.Itoa(position)
	textView.Highlight(pos).ScrollToHighlight()

}

func incrementPosition(position *int, by, max int) {
	if *position >= max {
		*position = max
		return
	}
	*position += by
}

func decrementPosition(position *int, by int) {
	min := 0
	if *position <= min {
		*position = min
		return
	}
	*position -= by

}
