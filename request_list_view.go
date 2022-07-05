package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RequestListModel struct {
	Requests         []Request
	CurrentSelection int
}

func (m RequestListModel) Init() tea.Cmd {
	return nil
}

func (m RequestListModel) View() string {
	style := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
	content := ""

	for idx, request := range m.Requests {
		row := ""
		if idx == m.CurrentSelection {
			row += "> "
		} else {
			row += "  "
		}
		newline := ""
		if idx < len(m.Requests)-1 {
			newline = "\n"
		}
		row += fmt.Sprintf("%s (%s)%s", request.URL, request.Method, newline)
		content += row
	}

	return style.Render(content) + "\n"
}

type SelectMsg struct {
	Request Request
}

func selectItem(request Request) tea.Cmd {
	return func() tea.Msg {
		return SelectMsg{
			Request: request,
		}
	}
}

func (m RequestListModel) Update(msg tea.Msg) (RequestListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", tea.KeyDown.String():
			if m.CurrentSelection < len(m.Requests)-1 {
				m.CurrentSelection++
			}
			return m, nil
		case "k", tea.KeyUp.String():
			if m.CurrentSelection > 0 {
				m.CurrentSelection--
			}
			return m, nil
		case tea.KeyEnter.String():
			return m, selectItem(m.Requests[m.CurrentSelection])
		}
	}

	return m, nil
}
