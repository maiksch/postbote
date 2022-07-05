package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RequestDetailView struct {
	Request Request
	Body    textarea.Model
}

func NewDetailView(request Request) RequestDetailView {
	textarea := textarea.New()
	log.Println("YEET")
	if request.HasBody() {
		log.Println("Has Body")
		body, err := json.Marshal(request.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var out bytes.Buffer
		json.Indent(&out, body, "", "\t")
		log.Println("Body: " + string(out.Bytes()))
		textarea.SetValue(string(out.Bytes())) 
	}

	return RequestDetailView{
		Request: request,
		Body:    textarea,
	}
}

func (m RequestDetailView) Init() tea.Cmd {
	return nil
}

func (m RequestDetailView) Update(msg tea.Msg) (RequestDetailView, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String():
			if !m.Body.Focused() {
				cmds = append(cmds, m.Body.Focus())
			}
		case tea.KeyEsc.String():
			if m.Body.Focused() {
				m.Body.Blur()
			}
		}
	}

	m.Body, cmd = m.Body.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m RequestDetailView) View() string {
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	view := fmt.Sprintf("URL: %s\n", m.Request.URL)
	view += fmt.Sprintf("Method: %s\n", m.Request.Method)

	view += "Params:\n"
	for key, value := range m.Request.Params {
		switch value := value.(type) {
		case string:
			view += fmt.Sprintf("\t%s = %s\n", key, value)
			break
		case int:
			view += fmt.Sprintf("\t%s = %d\n", key, value)
			break
		}
	}

	view += fmt.Sprintf("Focus: %s", m.Body.Focused())

	view += fmt.Sprintf("\n\n%s", m.Body.View())

	return style.Render(view)
}
