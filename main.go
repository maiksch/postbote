package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	// "github.com/charmbracelet/bubbles/key"
	// "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

type Model struct {
	Selected int
	Requests  []Request
}

// var appStyle = lipgloss.NewStyle().Padding(1, 2)

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	view := ""

	for idx, request := range m.Requests {
		if idx == m.Selected {
			view += ">"
		} else {
			view += " "
		}
		view += fmt.Sprintf(" %s (%s)\n", request.URL, request.Method)
	}

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
		Requests: data.Requests,
		Selected: 0,
	}
	program := tea.NewProgram(model)

	if err := program.Start(); err != nil {
		log.Fatalln(err)
	}

	// for _, request := range data.Requests {
	// 	Send(request)
	// }
}
