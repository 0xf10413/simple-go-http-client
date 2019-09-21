package main

import (
  "log"
  "time"
  "fmt"
  "io"
  "net/http"
  "encoding/json"
  "strings"
  "io/ioutil"
  ui "github.com/gizak/termui/v3"
  "github.com/gizak/termui/v3/widgets"
)

func main() {
  if err := ui.Init(); err != nil {
  log.Fatalf("failed to initialize termui: %v", err)
}
  defer ui.Close()

  // Set up main page
  mainPageHello := widgets.NewParagraph()
  mainPageHello.Title = "Short Demo"
  mainPageHello.Text = "Welcome to this short demo!"
  mainPageGrid := ui.NewGrid()
  mainPageGrid.Set(
    ui.NewRow(1.0,
      ui.NewCol(1.0, mainPageHello),
    ),
  )
  termWidth, termHeight := ui.TerminalDimensions()
  mainPageGrid.SetRect(0, 3, termWidth, termHeight)


  // Set up proxy list page
  proxyList := widgets.NewList()
  proxyList.Title = "Proxies"
  proxyList.WrapText = true
  proxyListMainGrid := ui.NewGrid()
  proxyListMainGrid.Set(
    ui.NewRow(1.0,
     ui.NewCol(1.0, proxyList),
   ),
  )
  proxyListMainGrid.SetRect(0, 3, termWidth, termHeight)

  // Set up tab pane
  tabpane := widgets.NewTabPane("Home page (1)", "Demo (2)")
  tabpane.SetRect(0, 0 ,termWidth, 3)
  tabpane.Border = true

  renderTab := func() {
    switch tabpane.ActiveTabIndex {
    case 0:
      ui.Render(mainPageGrid)
    case 1:
      ui.Render(proxyListMainGrid)
    }
  }

  uiEvents := ui.PollEvents()
  ticker := time.NewTicker(time.Second).C
  ui.Render(mainPageGrid, tabpane)
  renderTab()
  for {
    select {
    case e := <-uiEvents:
      switch e.ID {
      case "q":
        return
      case "1":
        tabpane.ActiveTabIndex = 0
        ui.Clear()
        ui.Render(tabpane)
        renderTab()
      case "2":
        tabpane.ActiveTabIndex = 1
        ui.Clear()
        ui.Render(tabpane)
        renderTab()
      case "r":
        proxyList.Rows = []string{
          "Please wait while we retrieve the data…",
        }
        ui.Clear()
        ui.Render(tabpane)
        renderTab()
        resp, err := http.Get("http://localhost:8081/get-proxies")
        if err != nil {
          proxyList.Rows = []string{
            fmt.Sprintf("Error during data retrieval : %s", err.Error()),
          }
        } else {
          defer resp.Body.Close()
          body, _ := ioutil.ReadAll(resp.Body)
          bodyStr := string(body)
          decoder := json.NewDecoder(strings.NewReader(bodyStr))
          var proxyArray []Proxy
          if err := decoder.Decode(&proxyArray); err == io.EOF {
            break;
          }  else if err != nil {
            log.Fatal(err)
          }
          var proxies []string
          for _, proxy := range proxyArray {
            proxies = append(proxies, fmt.Sprintf("%s: %s-%s",
              proxy.Name, proxy.IP, proxy.Port))
          }
          proxyList.Rows = proxies
        }
        ui.Clear()
        ui.Render(tabpane)
        renderTab()
     //case "<Resize>":
     //  payload := e.Payload.(ui.Resize)
     //  width, height := payload.Width, payload.Height
     //  mainPageGrid.SetRect(0, 0, width, height)
     //  ui.Render(mainPageGrid)
      }
    case <- ticker:
      //ui.Render(mainPageGrid)
    }
  }
}
