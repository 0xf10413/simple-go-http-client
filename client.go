package main

import (
  "github.com/rivo/tview"
	"github.com/gdamore/tcell"
  "net/http"
  "fmt"
  "io/ioutil"
  "time"
  "encoding/json"
  "strings"
)

func refreshProxyInfo(app *tview.Application, textView *tview.TextView,
  list *tview.List) {
  // TODO: protect against consecutive calls

  // Update the UI to inform we're about to fetch data
  app.QueueUpdateDraw(func () {
    textView.SetText("Query started, please wait for the result...")
  })

  // Retrieve data and format it
  textBoxResult := "N/A"
  var proxies []Proxy

  resp, err := http.Get("http://localhost:8081/get-proxies")
  if err != nil {
    textBoxResult = fmt.Sprintf("Could not retrieve data: %s", err)
    goto finish
  } else {
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      textBoxResult = fmt.Sprintf("Could not read data from the server: %s", err)
      goto finish
    }

    // Compute new server list
    bodyStr := string(body)
    decoder := json.NewDecoder(strings.NewReader(bodyStr))
    if err := decoder.Decode(&proxies); err != nil{
      textBoxResult = "Got a reply from the server, but could not decode it"
    } else {
      textBoxResult = "Data successfully retrieved"
    }
  }

finish:
  // Prefix the result by current time for tracability
  textBoxResult = fmt.Sprintf("%s: %s", time.Now().String(), textBoxResult)

  // Update text box and list of servers (if applicable)
  app.QueueUpdateDraw(func () {
    if len(proxies) > 0 {
      list.Clear()
    }
    for index, proxy := range proxies {
      list.AddItem(proxy.Name, fmt.Sprintf("%s:%s", proxy.IP, proxy.Port),
      rune(index+int('a')), nil)
    }
    textView.SetText(textBoxResult)
  })
}

func main() {
  app := tview.NewApplication()

  // Set up global layoyt
  pageInfo := tview.NewTextView().
                  SetDynamicColors(true).
                  SetRegions(true).
                  SetWrap(false).
                  SetText(`F1 ["1"]HomePage[""]  F2 ["2"]Proxies[""]`)
  pageInfo.Highlight("1")

  pages := tview.NewPages()
  globalLayout := tview.NewFlex().
                  SetDirection(tview.FlexRow).
                  AddItem(pageInfo, 1, 1, false).
                  AddItem(pages, 0, 1, true)

  // Set up main page
  mainPageHello := tview.NewBox().
                  SetBorder(true).
                  SetTitle("Short demo")
  mainPageGrid := tview.NewGrid().
                  SetBorders(true).
                  AddItem(mainPageHello, 0, 0, 1, 1, 1, 1, false)
  pages.AddPage("HomePage", mainPageGrid, true, true)

  // Set up proxy list page
  proxyList := tview.NewList()
  proxyInfo := tview.NewTextView().
               SetText("Press <F5> to retrieve proxy information")
  proxyListMainGrid := tview.NewGrid().
                       SetRows(3, 0, 3).
                       SetColumns(30, 0, 30).
                       SetBorders(true).
                       AddItem(proxyInfo, 0, 0, 1, 3, 0, 0, true).
                       AddItem(proxyList, 1, 0, 1, 3, 0, 0, true)
  pages.AddPage("Proxies", proxyListMainGrid, true, false)

  // Install event handler
  app.SetInputCapture(func (event *tcell.EventKey) *tcell.EventKey {
    if event.Key() == tcell.KeyEsc {
      app.Stop()
    } else if event.Key() == tcell.KeyF1 {
      pages.SwitchToPage("HomePage")
      pageInfo.Highlight("1")
    } else if event.Key() == tcell.KeyF2 {
      pages.SwitchToPage("Proxies")
      pageInfo.Highlight("2")
      app.SetFocus(proxyList)
    } else if event.Key() == tcell.KeyF5 {
      go refreshProxyInfo(app, proxyInfo, proxyList)
    }
    return event
  })

  // Run application
  if err := app.SetRoot(globalLayout, true).Run(); err != nil {
    panic(err)
  }
}
