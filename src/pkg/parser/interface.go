package parser

import (
	"time"

	"github.com/logviewer/v2/src/pkg/source"
	"github.com/samber/lo"
)

var AllParser = map[string]Parse{}

// support parsing structured log from some source
// log field with the same name MUST have the same value type
// TODO: hide field, support filtering key val
type Parse interface {
	Init(r source.Reader) error
	GetLogs(offset, maxAmt int) ([]LogItem, error)
	AvailableFields() []string
}

// should be called during init phase
func Register(id string, r Parse) {
	if _, ok := AllParser[id]; ok {
		panic("duplicate parser id: " + id)
	}
	AllParser[id] = r
}

type FieldType uint8

const (
	FieldInt FieldType = iota
	FieldFloat
	FieldStr
	FieldUnknown
)

type LogLevel uint8

const (
	LogTrace LogLevel = iota
	LogDebug
	LogInfo
	LogWarn
	LogErr
	LogUnknown
)

func (l LogLevel) ToString() string {
	switch l {
	case LogTrace:
		return "TRC"
	case LogDebug:
		return "DBG"
	case LogInfo:
		return "INF"
	case LogWarn:
		return "WRN"
	case LogErr:
		return "ERR"
	default:
		return "UNK"
	}
}

// represent single line info
type LogItem interface {
	TimeStamp() time.Time
	Level() LogLevel
	Caller() string
	Msg() string
	Line() string
	SortedFields() []lo.Tuple2[string, any]
}
