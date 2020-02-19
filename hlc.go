package hlc

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	doubleQuote = "\""
	dash        = "-"
)

var (
	byteDash = []byte("-")
)

var (
	// ErrWrongFormat is an error which only happens during UnmarshalJSON
	// if and only if given data is not parsable to Timestamp data structure
	ErrWrongFormat = errors.New("wrong format")
)

func pt() int64 {
	return time.Now().UTC().UnixNano()
}

func max2(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func max3(a, b, c int64) int64 {
	if a < b {
		a = b
	}

	if a < c {
		a = c
	}

	return a
}

// Timestamp struct to capture counter and time in nanoseconds
type Timestamp struct {
	counter int64
	time    int64
}

func (ts *Timestamp) String() string {
	return fmt.Sprintf("%s-%s", strconv.FormatInt(ts.time, 16), strconv.FormatInt(ts.counter, 16))
}

// MarshalJSON overrides and implements how timestamp needs to be encoded in JSON
func (ts *Timestamp) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.WriteString(doubleQuote)
	buffer.WriteString(ts.String())
	buffer.WriteString(doubleQuote)

	return buffer.Bytes(), nil
}

// UnmarshalJSON overrides and imeplements how timestamp should be parsed from JSON
func (ts *Timestamp) UnmarshalJSON(data []byte) error {
	// need to remove quotes from data
	data = bytes.Trim(data, doubleQuote)

	segments := bytes.Split(data, byteDash)
	if len(segments) != 2 {
		return ErrWrongFormat
	}

	var err error
	ts.time, err = strconv.ParseInt(string(segments[0]), 16, 64)
	if err != nil {
		return err
	}

	ts.time, err = strconv.ParseInt(string(segments[1]), 16, 64)
	if err != nil {
		return err
	}

	return nil
}

// Less checks whether the given timestamp is bigger than current one
func (ts *Timestamp) Less(recv *Timestamp) bool {
	switch {
	case ts.time < recv.time:
		return true
	case ts.time == recv.time && ts.counter < recv.counter:
		return true
	default:
		return false
	}
}

// Now creates a new timestamp based on current clock
// This method should be called if sending a message or local change is occured
func (ts *Timestamp) Now() *Timestamp {
	t := ts.time
	ts.time = max2(t, pt())

	if ts.time == t {
		ts.counter++
	} else {
		ts.counter = 0
	}

	return &Timestamp{
		counter: ts.counter,
		time:    ts.time,
	}
}

// Update the current clock, this should be called once a msg is recceived from
// other nodes
func (ts *Timestamp) Update(msg *Timestamp) {
	t := ts.time
	ts.time = max3(t, msg.time, pt())
	if ts.time == t && t == msg.time {
		ts.counter = max2(ts.counter, msg.counter) + 1
	} else if ts.time == t {
		ts.counter++
	} else if ts.time == msg.time {
		ts.counter = msg.counter + 1
	} else {
		ts.counter = 0
	}
}

// New creates a brand new Clock
// This function should be called once per node
func New() *Timestamp {
	return &Timestamp{}
}
