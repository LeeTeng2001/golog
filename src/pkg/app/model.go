// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/logviewer/v2/src/pkg/parser"
	zone "github.com/lrstanley/bubblezone"
)

// This is a modified version of this example, supporting full screen, dynamic
// resizing, and clickable models (tabs, lists, dialogs, etc).
// 	https://github.com/charmbracelet/lipgloss/blob/master/example

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
)

type model struct {
	// window property
	height int
	width  int

	// components
	tabs    tea.Model
	logList tea.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) isInitialized() bool {
	return m.height != 0 && m.width != 0
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.isInitialized() {
		if _, ok := msg.(tea.WindowSizeMsg); !ok {
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
		// toggle mouse track
		// if msg.String() == "ctrl+e" {
		// 	zone.SetEnabled(!zone.Enabled())
		// 	return m, nil
		// }
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		msg.Height -= 2
		msg.Width -= 4
		return m.propagate(msg), nil
	}

	return m.propagate(msg), nil
}

func (m *model) propagate(msg tea.Msg) tea.Model {
	// Propagate to all children.
	// m.tabs, _ = m.tabs.Update(msg)
	m.logList, _ = m.logList.Update(msg)

	// if msg, ok := msg.(tea.WindowSizeMsg); ok {
	// 	msg.Height -= m.tabs.(tabs).height + m.list1.(list).height
	// 	m.history, _ = m.history.Update(msg)
	// 	return m
	// }

	// m.history, _ = m.history.Update(msg)
	return m
}

func (m model) View() string {
	if !m.isInitialized() {
		return "initialising..."
	}

	s := lipgloss.NewStyle().MaxHeight(m.height).MaxWidth(m.width)
	return zone.Scan(s.Render(lipgloss.JoinVertical(lipgloss.Top,
		// m.tabs.View(),
		// "xxx",
		m.logList.View(),
		// lipgloss.PlaceHorizontal(
		// 	m.width, lipgloss.Center,
		// 	lipgloss.JoinHorizontal(
		// 		lipgloss.Top,
		// 		m.list1.View(), m.list2.View(), m.dialog.View(),
		// 	),
		// 	lipgloss.WithWhitespaceChars(" "),
		// ),
		// m.history.View(),
	)))
}

func Main(ps parser.Parse) error {
	// Initialize a global zone manager, so we don't have to pass around the manager
	// throughout components.
	zone.NewGlobal()

	l, err := lru.New[int, *logCache](512)
	if err != nil {
		return err
	}
	m := &model{
		// tabs: &tabs{
		// 	id:     zone.NewPrefix(), // Give each type an ID, so no zones will conflict.
		// 	height: 2,
		// 	active: "Lip Gloss",
		// 	items:  []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara"},
		// },
		logList: &loglist{
			renderLru: l,
			id:        zone.NewPrefix(),
			ps:        ps,
		},
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
