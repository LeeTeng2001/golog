package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/logviewer/v2/src/pkg/common"
	"github.com/logviewer/v2/src/pkg/source"
	"github.com/samber/lo"
)

var _ Parse = (*Zap)(nil)

func init() {
	Register("zap", &Zap{})
}

type Zap struct {
	logEntries []LogItem
	seenFields []string
	decoder    *json.Decoder
}

func (z *Zap) Init(r source.Reader) error {
	z.decoder = json.NewDecoder(r)
	z.seenFields = make([]string, 0)
	z.logEntries = make([]LogItem, 0)
	return nil
}

func (z *Zap) AvailableFields() []string {
	return z.seenFields
}

func (z *Zap) GetLogs(offset int, maxAmt int) ([]LogItem, error) {
	decodeFrom := min(len(z.logEntries), offset)
	decodeTo := offset + maxAmt
	linesToDecode := decodeTo - decodeFrom

	// might hit eof
	for range linesToDecode {
		item := logItem{}
		if err := z.decoder.Decode(&item); err != nil {
			fmt.Println("decode err", err)
			if err == io.EOF {
				fmt.Println("EOF")
				break
			} else { // ignore invalid line?
				continue
			}
		}
		fmt.Printf("%v\n", item)
		z.logEntries = append(z.logEntries, &item)
	}

	if offset >= len(z.logEntries) {
		return []LogItem{}, nil
	}
	return z.logEntries[decodeFrom:min(len(z.logEntries), decodeTo)], nil
}

var _ LogItem = (*logItem)(nil)

type logItem struct {
	Ts              common.TimeConverter   `json:"ts"`
	Lv              string                 `json:"level"`
	Call            string                 `json:"logger"`
	Ms              string                 `json:"msg"`
	Li              string                 `json:"caller"`
	Other           map[string]interface{} `json:"-"`
	otherSerialised []lo.Tuple2[string, any]
}

func (l *logItem) TimeStamp() time.Time {
	return time.Time(l.Ts)
}
func (l *logItem) Level() LogLevel {
	switch l.Lv {
	case "trace":
		return LogTrace
	case "debug":
		return LogDebug
	case "info":
		return LogInfo
	case "warn":
		return LogWarn
	case "error":
		return LogErr
	default:
		return LogUnknown
	}
}
func (l *logItem) Caller() string { return l.Call }
func (l *logItem) Msg() string    { return l.Ms }
func (l *logItem) Line() string   { return l.Li }
func (l *logItem) SortedFields() []lo.Tuple2[string, any] {
	// lazy serialisation
	if l.otherSerialised == nil {
		keys := lo.Keys(l.Other)
		sort.Strings(keys)
		for _, k := range keys {
			l.otherSerialised = append(l.otherSerialised, lo.T2(k, l.Other[k]))
		}
	}
	return l.otherSerialised
}
