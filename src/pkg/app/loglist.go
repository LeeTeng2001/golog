package app

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lru "github.com/hashicorp/golang-lru/v2"
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
			Foreground(lightGray)
	boldTextStyle = lipgloss.NewStyle().
			Foreground(gray).
			Bold(true)
	warnLogTextStyle = lipgloss.NewStyle().
				Foreground(yellow).
				Background(darkYellow).
				Bold(true)
	errLogTextStyle = lipgloss.NewStyle().
			Foreground(red).
			Background(darkRed).
			Bold(true)
	infoLogTextStyle = lipgloss.NewStyle().
				Foreground(blue).
				Bold(true)
	defaultLogTextStyle = lipgloss.NewStyle().
				Foreground(darkGray).
				Bold(true)
	keyTextStyle = lipgloss.NewStyle().
			Foreground(lightGreen)
)

func othersRender(sortedMap []lo.Tuple2[string, any]) string {
	return strings.Join(lo.Map(sortedMap, func(t lo.Tuple2[string, any], _ int) string {
		return keyTextStyle.Render(t.A) + "=" + subtleTextStyle.Render(fmt.Sprintf("%v", t.B))
	}), " ")
}

type logCache struct {
	lastRender      string // cache
	lastOffset      int    // cache
	totalLineLength int
	renderPart      []lo.Tuple2[*lipgloss.Style, string] // seperated by single space
}

func newLogCache(item parser.LogItem) *logCache {
	// determine log style
	logStyle := &defaultLogTextStyle
	switch item.Level() {
	case parser.LogInfo:
		logStyle = &infoLogTextStyle
	case parser.LogWarn:
		logStyle = &warnLogTextStyle
	case parser.LogErr:
		logStyle = &errLogTextStyle
	}

	// compute main part
	l := &logCache{
		lastOffset: -1,
	}
	l.renderPart = append(l.renderPart,
		lo.T2(&subtleTextStyle, item.TimeStamp().Format("Mon Jan 02 15:04:05.000 ")),
		lo.T2(logStyle, "["+item.Level().ToString()+"] "),
		lo.T2(&subtleTextStyle, item.Caller()+": "),
		lo.T2(&boldTextStyle, item.Msg()+" "),
	)

	// compute other key val
	for _, t := range item.SortedFields() {
		l.renderPart = append(l.renderPart,
			lo.T2(&keyTextStyle, t.A),
			lo.T2(&subtleTextStyle, fmt.Sprintf("=%v ", t.B)),
		)
	}
	// remove last element white space
	ln := len(l.renderPart)
	l.renderPart[ln-1].B = l.renderPart[ln-1].B[:len(l.renderPart[ln-1].B)-1]

	l.totalLineLength = lo.SumBy(l.renderPart, func(item lo.Tuple2[*lipgloss.Style, string]) int {
		return len(item.B)
	})

	return l
}

func (lc *logCache) renderLine(xOffset int) (res string) {
	// fast path
	if lc.lastOffset == xOffset {
		return lc.lastRender
	}

	// save cache result
	defer func() {
		lc.lastOffset = xOffset
		lc.lastRender = res
	}()
	if xOffset >= lc.totalLineLength {
		return ""
	}

	// compute parts, TODO: optimise modify only first and last element when moving horizontally
	remainingOffset := xOffset
	for _, t := range lc.renderPart {
		part := t.B
		style := t.A
		if remainingOffset == 0 { // in bound
			res += style.Render(part)
		} else if remainingOffset >= len(part) { // not shown
			remainingOffset -= len(part)
		} else { // trim in bound
			res += style.Render(part[remainingOffset:])
			remainingOffset = 0
		}
	}
	slog.Info(res)

	return
}

type loglist struct {
	id            string
	height        int
	width         int
	curLogYOffset int
	curLogXOffset int

	renderLru *lru.Cache[int, *logCache]
	ps        parser.Parse
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
			m.curLogYOffset = max(m.curLogYOffset-1, 0)
		case "down":
			m.curLogYOffset++
		case "right":
			m.curLogXOffset += 2
		case "left":
			m.curLogXOffset = max(m.curLogXOffset-2, 0)
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

	// get current log and prerender full line content
	if items, err := m.ps.GetLogs(m.curLogYOffset, m.height); err == nil {
		for idx, item := range items {
			if !m.renderLru.Contains(m.curLogYOffset + idx) {
				m.renderLru.Add(m.curLogYOffset+idx, newLogCache(item))
			}
		}
	}

	return m, nil
}

func (m loglist) View() string {
	out := []string{}
	for offset := range m.height {
		slog.Info("xxx")
		cache, ok := m.renderLru.Get(offset + m.curLogYOffset)
		if !ok {
			break
		}
		out = append(out, zone.Mark(
			m.id+strconv.Itoa(offset),
			cache.renderLine(m.curLogXOffset),
		))
	}

	return listStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, out...),
	)
}
