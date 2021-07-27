package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"url-shortner/api_tool/textBox"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {

	app := tview.NewApplication()

	form := tview.NewForm().SetHorizontal(false)
	form.SetTitle("Request")
	form.AddInputField("URL", "http://localhost:8080/shorten", 100, nil, nil).Box.SetBorder(true)

	headerCount := 0
	form.AddButton("add-header", func() {
		label := fmt.Sprintf("%s%d", "Header", headerCount+1)
		form.AddInputField(label, "Content-Type:application/json", 100, nil, nil)
		field := form.GetFormItemByLabel(label).(*tview.InputField)
		app.SetFocus(field)
		app.SetFocus(form)
		headerCount++

	})

	form.AddButton("add-cookie", func() {
		cookieField := form.AddInputField("cookie", "", 100, nil, nil)
		app.SetFocus(cookieField)
		app.SetFocus(form)

	})
	form.AddButton("Reset", func() {

		for i := 0; i < headerCount; i++ {
			form.RemoveFormItem(form.GetFormItemIndex(fmt.Sprintf("%s%d", "Header", i+1)))

		}
		headerCount = 0
		form.GetFormItemByLabel("URL").(*tview.InputField).SetText("")

	})

	form.AddDropDown("Method", []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodOptions, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace, http.MethodPatch}, 0, nil)

	responseView := tview.NewTextView()
	responseView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyCtrlB {
			clipboard.WriteAll(responseView.GetText(true))
		}
		return event

	})

	errorView := tview.NewTextView().SetTextColor(tcell.ColorRed)
	responseView.SetBorder(true).SetTitle("Response")
	errorView.SetBorder(true).SetTitle("Error")

	txtbox := textBox.TextBox(app)
	focusList := []tview.Primitive{form, txtbox, responseView, errorView}
	focusIndex := 0
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlD:
			focusIndex++
			if focusIndex >= len(focusList) {
				focusIndex = 0
			}
			app.SetFocus(focusList[focusIndex])
			break

		}
		return event
	})

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(form, 0, 2, true).
			AddItem(txtbox, 0, 1, false), 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(responseView, 0, 3, false).
			AddItem(errorView, 0, 1, false), 0, 1, false)

	form.AddButton("Submit", func() {

		responseView.SetText("Loading...")
		resp, err := submitRequest(form, txtbox, headerCount)
		if err != nil {
			errorView.SetText(err.Error())
			return
		}
		errorView.SetText("")
		responseView.SetText(resp)

	})

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func submitRequest(form *tview.Form, txtBox *tview.TextView, headerCount int) (string, error) {

	urlString := form.GetFormItemByLabel("URL").(*tview.InputField).GetText()
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	_, method := form.GetFormItemByLabel("Method").(*tview.DropDown).GetCurrentOption()

	client := &http.Client{}

	req, err := http.NewRequest(method, parsedUrl.String(), bytes.NewBuffer([]byte(txtBox.GetText(true))))
	if err != nil {
		return "", err
	}

	for i := 0; i < headerCount; i++ {
		label := fmt.Sprintf("%s%d", "Header", i+1)
		header := strings.Split(form.GetFormItemByLabel(label).(*tview.InputField).GetText(), ":")
		req.Header.Set(strings.TrimSpace(header[0]), strings.TrimSpace(header[1]))
	}

	if form.GetFormItemByLabel("cookie") != nil {
		text := form.GetFormItemByLabel("cookie").(*tview.InputField).GetText()
		if text != "" {
			split := strings.Split(form.GetFormItemByLabel("cookie").(*tview.InputField).GetText(), ";")
			for _, value := range split {
				splitSign := strings.Split(value, "=")
				key := splitSign[0]
				v := splitSign[1]
				cookieValue := &http.Cookie{Name: key, Value: v}
				req.AddCookie(cookieValue)
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {

		return "", err
	}

	if resp.Header.Get("Content-Type") == "application/json" {
		formatted := bytes.NewBuffer([]byte(""))
		err := json.Indent(formatted, data, "", "\t")

		if err != nil {
			return fmt.Sprintf("Status Code: %d\n\n%s", resp.StatusCode, string(data)), nil
		}
		return fmt.Sprintf("Status Code: %d\n\n%s", resp.StatusCode, formatted.Bytes()), nil
	}

	return fmt.Sprintf("Status Code: %d\n\n%s", resp.StatusCode, string(data)), nil
}
