package parser

import (
	"encoding/json"
	"io"
	"sort"
	"time"

	"github.com/logviewer/v2/src/pkg/common"
	"github.com/logviewer/v2/src/pkg/common/slogx"
	"github.com/logviewer/v2/src/pkg/source"
	"github.com/mitchellh/mapstructure"
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
		item := zapLogItem{}
		mp, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: mapstructure.TextUnmarshallerHookFunc(),
			Result:     &item,
		})
		if err != nil {
			slogx.Error("create mapstructure error", err)
			return nil, err
		}

		rawJsonMap := map[string]any{}
		if err := z.decoder.Decode(&rawJsonMap); err != nil {
			if err == io.EOF {
				break
			} else { // ignore invalid line?
				slogx.Error("decode err", err)
				continue
			}
		}
		if err = mp.Decode(rawJsonMap); err != nil {
			slogx.Error("decode map error", err)
			return nil, err
		}
		z.logEntries = append(z.logEntries, &item)
	}

	if offset >= len(z.logEntries) {
		return []LogItem{}, nil
	}
	return z.logEntries[decodeFrom:min(len(z.logEntries), decodeTo)], nil
}

var _ LogItem = (*zapLogItem)(nil)

type zapLogItem struct {
	Ts              common.TimeConverter   `mapstructure:"ts"`
	Lv              string                 `mapstructure:"level"`
	Call            string                 `mapstructure:"logger"`
	Ms              string                 `mapstructure:"msg"`
	Li              string                 `mapstructure:"caller"`
	Other           map[string]interface{} `mapstructure:",remain"`
	otherSerialised []lo.Tuple2[string, any]
}

func (l *zapLogItem) TimeStamp() time.Time {
	return time.Time(l.Ts)
}
func (l *zapLogItem) Level() LogLevel {
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
func (l *zapLogItem) Caller() string { return l.Call }
func (l *zapLogItem) Msg() string    { return l.Ms }
func (l *zapLogItem) Line() string   { return l.Li }
func (l *zapLogItem) SortedFields() []lo.Tuple2[string, any] {
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
