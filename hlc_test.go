package hlc_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/alinz/hlc"
)

func random(min, max int) int {
	return rand.Intn(max-min+1) + min
}

type Message struct {
	id    int
	value int
	ts    *hlc.Timestamp
}

func TestLogic(t *testing.T) {
	clock := hlc.New()

	ts := clock.Now()

	clock.Update(ts)

	if !ts.Less(clock.Now()) {
		t.Fatal("ts is less than current clock time")
	}
}

func TestJsonHLC(t *testing.T) {
	clock := hlc.New()

	value := struct {
		Timestamp *hlc.Timestamp `json:"timestamp"`
	}{
		Timestamp: clock.Now(),
	}

	fmt.Println(value)

	b, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))

	err = json.Unmarshal(b, &value)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(value)
}
