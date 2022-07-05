package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	// "github.com/charmbracelet/bubbles/key"
	// "github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
	// "github.com/charmbracelet/lipgloss"
)

type Request struct {
	URL         string                 `json:"url"`
	Method      string                 `json:"method"`
	Body        interface{}            `json:"body"`
	ContentType string                 `json:"contentType"`
	Params      map[string]interface{} `json:"params"`
}

func (req Request) HasBody() bool {
	return req.Body != nil && req.Method != "GET" && req.Method != "DELETE"
}

type Data struct {
	Requests []Request `json:"requests"`
}

type View int

const (
	listView   View = 0
	detailView View = iota
)

type Model struct {
	Focus           View
	ListView        RequestListModel
	DetailView      RequestDetailView
	SelectedRequest Request
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case SelectMsg:
		m.DetailView.Request = msg.Request
		m.Focus = detailView
		if msg.Request.HasBody() {
			log.Println("Has Body")
			body, err := json.Marshal(msg.Request.Body)
			if err != nil {
				log.Fatalln(err)
			}

			var out bytes.Buffer
			json.Indent(&out, body, "", "\t")
			log.Println("Body: " + string(out.Bytes()))
			m.DetailView.Body.SetValue(string(out.Bytes()))
		}
		return m, nil
	}

	if m.Focus == listView {
		m.ListView, cmd = m.ListView.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.Focus == detailView {
		m.DetailView, cmd = m.DetailView.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	view := m.ListView.View()
	view += m.DetailView.View()
	return view
}

func ReadResponse(response *http.Response) string {
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

func CreateQueryURL(req *http.Request, params map[string]interface{}) string {
	query := req.URL.Query()
	for key, value := range params {
		switch value.(type) {
		case string:
			query.Add(key, value.(string))
			break
		case int:
			query.Add(key, string(rune(value.(int))))
			break
		}
	}
	return query.Encode()
}

func Send(request Request) {
	body := new(bytes.Buffer)
	if request.HasBody() {
		tmp, err := json.Marshal(request.Body)
		if err != nil {
			log.Fatalln(err)
		}
		body = bytes.NewBuffer(tmp)
	}

	req, err := http.NewRequest(request.Method, request.URL, body)
	if err != nil {
		panic(err)
	}

	req.URL.RawQuery = CreateQueryURL(req, request.Params)

	if request.ContentType != "" {
		req.Header.Set("Content-Type", request.ContentType)
	}

	log.Println("Start request: " + req.URL.String())
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	response := ReadResponse(res)
	log.Println(response)
}

func main() {
	log.Println("Postbote start")

	file, err := ioutil.ReadFile("postbote.json")
	if err != nil {
		panic(err)
	}

	data := new(Data)
	if err := json.Unmarshal(file, &data); err != nil {
		panic(err)
	}

	model := &Model{
		ListView: RequestListModel{
			Requests:         data.Requests,
			CurrentSelection: 0,
		},
		DetailView: NewDetailView(data.Requests[0]),
	}
	program := tea.NewProgram(model)

	if err := program.Start(); err != nil {
		log.Fatalln(err)
	}
}
