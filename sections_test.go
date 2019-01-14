package bufr

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var r, err = os.Open("./test/JUVE00 EGRR 161200")

func TestSection0(t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
	// defer r.Close()

	sect, err := NewSection0(r)
	if err != nil {
		t.Fatal(err)
	}

	expected0 := &Section0{MagicString: [4]byte{'B', 'U', 'F', 'R'}, Len: 430, Version: 3}
	if !reflect.DeepEqual(sect, expected0) {
		t.Errorf("expected %v, but found: %v", expected0, sect)
	}
	fmt.Printf("Section0: %+v\n", sect)
}

func TestSection1(t *testing.T) {
	TestSection0(t)

	sect, err := NewSection1(r)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Section1: %+v\n", sect)
	t.Fatal("OK")
}
