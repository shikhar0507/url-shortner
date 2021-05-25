package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gdamore/tcell/v2"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"
)

func main() {

	//	box := tview.NewBox().SetTitle("starting").SetBorder(true).SetTitleColor(tcell.Color126)
	app := tview.NewApplication()

	urlField := tview.NewInputField().
		SetLabel("URL: ").
		SetFieldWidth(50).
		SetDoneFunc(func(key tcell.Key) {

		})
	urlField.SetBorder(true)
	urlField.SetRect(0, 1, 50, 1)
	urlField2 := tview.NewInputField().
		SetLabel("URL 2: ").
		SetFieldWidth(50).
		SetDoneFunc(func(key tcell.Key) {

		})

	urlField2.Box.SetRect(100, 50, 50, 1)
	urlField2.SetBorder(true).SetBorderColor(tcell.Color191)
	urlField2.Box.SetBorderPadding(0, 0, 0, 0)

	form := tview.NewForm()
	form.AddInputField("url", "", 50, nil, nil)
	form.AddInputField("url 2", "", 50, nil, nil)

	buffer := femto.NewBufferFromString(string("Hello this is shikhar"), "/home/xanadu/")

	root := femto.NewView(buffer)
	root.SetRuntimeFiles(runtime.Files)
	root.SetRect(0, 0, 100, 40)

	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			app.Stop()
			return nil
		}
		return event
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		return event
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 1, false).
		AddItem(urlField, 0, 1, true)
	flex.AddItem(urlField2, 0, 1, false)
	flex.Box.SetBorderPadding(1, 1, 1, 1)
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
	return

	fmt.Println("Enter url name")
	var url, contentType, method string
	var body string = ""
	fmt.Scanf("%s", &url)

	fmt.Println("Method :")
	fmt.Scanf("%s", &method)

	fmt.Println("Content-Type :")
	fmt.Scanf("%s", &contentType)

	if method == http.MethodPatch || method == http.MethodPut || method == http.MethodPost {

		fmt.Println("Body :")
		fmt.Scanf("%s", &body)
	}

	buff := bytes.NewBuffer([]byte(body))
	fmt.Println("body", body)
	req, err := http.NewRequest(method, url, buff)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("error in client", err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("error parsing response", err)
	}
	fmt.Println(resp.StatusCode, string(respBody))
}
