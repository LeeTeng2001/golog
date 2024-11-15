package app

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/logviewer/v2/src/pkg/parser"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
)

var (
	listStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(subtle).
			MarginRight(2)
	subtleTextStyle = lipgloss.NewStyle().
			Foreground(lightGray).
			MarginRight(1)
	boldTextStyle = lipgloss.NewStyle().
			Foreground(gray).
			Bold(true).
			MarginRight(1)
	warnLogTextStyle = lipgloss.NewStyle().
				Foreground(yellow).
				Background(darkYellow).
				Bold(true).
				MarginRight(1)
	errLogTextStyle = lipgloss.NewStyle().
			Foreground(red).
			Background(darkRed).
			Bold(true).
			MarginRight(1)
	infoLogTextStyle = lipgloss.NewStyle().
				Foreground(blue).
				Bold(true).
				MarginRight(1)
	defaultLogTextStyle = lipgloss.NewStyle().
				Foreground(darkGray).
				Bold(true).
				MarginRight(1)
	keyTextStyle = lipgloss.NewStyle().
			Foreground(lightGreen)
)

func logRender(l parser.LogLevel) string {
	switch l {
	case parser.LogInfo:
		return infoLogTextStyle.Render("│" + l.ToString() + "│")
	case parser.LogWarn:
		return warnLogTextStyle.Render("│" + l.ToString() + "│")
	case parser.LogErr:
		return errLogTextStyle.Render("│" + l.ToString() + "│")
	default:
		return defaultLogTextStyle.Render("│" + l.ToString() + "│")
	}
}

func othersRender(sortedMap []lo.Tuple2[string, any]) string {
	return strings.Join(lo.Map(sortedMap, func(t lo.Tuple2[string, any], _ int) string {
		return keyTextStyle.Render(t.A) + "=" + subtleTextStyle.Render(fmt.Sprintf("%v", t.B))
	}), " ")
}

type loglist struct {
	id     string
	height int
	width  int
	ps     parser.Parse

	curLogOffset int
	logItem      []parser.LogItem
}

func (m loglist) Init() tea.Cmd {
	return nil
}

func (m loglist) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.curLogOffset = max(m.curLogOffset-1, 0)
		case "down":
			m.curLogOffset++
		}
		// TODO, click:
		// case tea.MouseMsg:
		// 	if msg.Action != tea.MouseActionRelease || msg.Button != tea.MouseButtonLeft {
		// 		return m, nil
		// 	}
		// 	// for i, item := range m.items {
		// 	// 	// Check each item to see if it's in bounds.
		// 	// 	if zone.Get(m.id + item.name).InBounds(msg) {
		// 	// 		m.items[i].done = !m.items[i].done
		// 	// 		break
		// 	// 	}
		// 	// }
		// 	return m, nil
	}

	if items, err := m.ps.GetLogs(m.curLogOffset, m.height); err == nil {
		m.logItem = items
	}
	// writeDebugFLn("%v", m.logItem)

	return m, nil
}

func (m loglist) View() string {
	out := []string{}
	for idx, item := range m.logItem {
		content := subtleTextStyle.Render(item.TimeStamp().Format("Mon Jan 02 15:04:05.000")) +
			logRender(item.Level()) +
			subtleTextStyle.Render(item.Caller()+":") +
			boldTextStyle.Render(item.Msg()) +
			othersRender(item.SortedFields())
		out = append(out, zone.Mark(
			m.id+strconv.Itoa(idx),
			content,
		))
	}

	return listStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, out...),
	)
}
