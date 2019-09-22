package main

import (
  "github.com/rivo/tview"
	"github.com/gdamore/tcell"
  "net/http"
  "fmt"
  "io/ioutil"
  "time"
  "encoding/json"
  "net/url"
  "strings"
  "net/http/cookiejar"
)

func refreshProxyInfo(app *tview.Application, client *http.Client, textView *tview.TextView,
  list *tview.List) {
  // TODO: protect against consecutive calls

  // Update the UI to inform we're about to fetch data
  app.QueueUpdateDraw(func () {
    textView.SetText("Query started, please wait for the result...")
  })

  // Retrieve data and format it
  textBoxResult := "N/A"
  var proxies []Proxy

  resp, err := client.Get("http://localhost:8081/get-proxies")
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
      textBoxResult = fmt.Sprintf(
        "Got a reply from the server, but could not decode it : %s", err)
        proxies = []Proxy{} // empty the proxies in case they were filled
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

func doLogin(client *http.Client, formData url.Values,textView *tview.TextView,
  app *tview.Application) {
  // TODO: protect against consecutive calls
  newJar, _ := cookiejar.New(nil)
  client.Jar = newJar

  app.QueueUpdateDraw(func () {
    textView.SetText("Query started, please wait for the result...")
  })

  resp, err := client.PostForm("http://localhost:8081/login", formData)
  resultText := "<N/A>"
  if err != nil {
    resultText = "Could not login to the server"
  } else if resp.StatusCode != 200 {
    resultText = fmt.Sprintf("Auth failed! Got error code %d", resp.StatusCode)
  } else {
    resultText = "Login successful!"
  }
  resp.Body.Close()
  app.QueueUpdateDraw(func () {
    textView.SetText(resultText)
  })
}

func main() {
  cookieJar, _ := cookiejar.New(nil)
  transport := &http.Transport{
    IdleConnTimeout: 5*time.Second,
  }
  client := &http.Client{Transport: transport, Jar: cookieJar}

  app := tview.NewApplication()

  // Set up global layout
  pageInfo := tview.NewTextView().
                  SetDynamicColors(true).
                  SetRegions(true).
                  SetWrap(false).
                  SetText(`F1 ["1"]Login[""]  F2 ["2"]Proxies[""]  F3 ["3"]Databases[""]`)
  pageInfo.Highlight("1")

  pages := tview.NewPages()
  globalLayout := tview.NewFlex().
                  SetDirection(tview.FlexRow).
                  AddItem(pageInfo, 1, 1, false).
                  AddItem(pages, 0, 1, true)

  // Set up main page
  mainPageResultTest := tview.NewTextView()
  mainPageLoginForm := tview.NewForm().
                       AddInputField("Login", "", 20, nil, nil).
                       AddPasswordField("Password", "", 20, '*', nil).
                       AddDropDown("Server", []string{"http://localhost:8081/",}, 0, nil)
  mainPageLoginForm.AddButton("Log in", func () {
                         formData := url.Values{
                           "username": {
                             mainPageLoginForm.GetFormItemByLabel("Login").(*tview.InputField).GetText()},
                           "password": {
                             mainPageLoginForm.GetFormItemByLabel("Password").(*tview.InputField).GetText()},
                         }
                         go doLogin(client, formData, mainPageResultTest, app)
                       })
  mainPageLoginForm.SetBorder(true).SetTitle("Log in to a server").SetTitleAlign(tview.AlignLeft)
  mainPageGrid := tview.NewGrid().
                  SetBorders(true).
                  AddItem(mainPageLoginForm, 0, 0, 1, 1, 1, 1, false).
                  AddItem(mainPageResultTest, 1, 0, 1, 1, 1, 1, false)
  pages.AddPage("Login", mainPageGrid, true, true)

  // Set up proxy list page
  proxyList := tview.NewList()
  proxyInfo := tview.NewTextView().
               SetText("Press <F5> to retrieve proxy information")
  proxyListMainGrid := tview.NewFlex().
                       SetDirection(tview.FlexRow).
                       AddItem(proxyInfo, 3, 0, false).
                       AddItem(proxyList, 0, 1, true)
  pages.AddPage("Proxies", proxyListMainGrid, true, false)

  // Set up database list page
  dbList := tview.NewList()
  dbInfo := tview.NewTextView().
               SetText("Press <F5> to retrieve db information")
  dbListMainGrid := tview.NewFlex().
                       SetDirection(tview.FlexRow).
                       AddItem(dbInfo, 3, 0, false).
                       AddItem(dbList, 0, 1, true)
  pages.AddPage("Databases", dbListMainGrid, true, false)

  currentPage := 1

  // Install event handler
  app.SetInputCapture(func (event *tcell.EventKey) *tcell.EventKey {
    if event.Key() == tcell.KeyEsc {
      app.Stop()
    } else if event.Key() == tcell.KeyF1 {
      pages.SwitchToPage("Login")
      currentPage = 1
      pageInfo.Highlight("1")
      app.SetFocus(mainPageLoginForm)
    } else if event.Key() == tcell.KeyF2 {
      pages.SwitchToPage("Proxies")
      pageInfo.Highlight("2")
      currentPage = 2
      app.SetFocus(proxyList)
    } else if event.Key() == tcell.KeyF3 {
      pages.SwitchToPage("Databases")
      pageInfo.Highlight("3")
      currentPage = 3
      app.SetFocus(proxyList)
    } else if event.Key() == tcell.KeyF5 {
      if currentPage == 2 {
        go refreshProxyInfo(app, client, proxyInfo, proxyList)
      }
    }
    return event
  })

  // Run application
  if err := app.SetRoot(globalLayout, true).SetFocus(mainPageLoginForm).Run(); err != nil {
    panic(err)
  }
}
