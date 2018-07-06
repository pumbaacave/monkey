package object

import "testing"

func TestSringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is A"}
	diff2 := &String{Value: "My name is A"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have the same hash keys")
	}
	if hello1.HashKey() == diff2.HashKey() {
		t.Errorf("strings with different content have the same hash keys")
	}
}
