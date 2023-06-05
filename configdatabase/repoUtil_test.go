package configdatabase

import (
	"testing"
)

func TestGetKeyIndexInfo(t *testing.T) {

	actual := getKeyIndexInfo(0, "success/fail")
	expected := "success"
	if actual != expected {
		t.Errorf("Test failed. Expected: %s, but got: %s", expected, actual)
	}
}

func TestSortLabels(t *testing.T) {

	actual := sortLabels("b;a")
	expected := "a;b"

	if actual != expected {
		t.Errorf("Test failed. Expected: %s, but got: %s", expected, actual)
	}
}

func TestDbKeyGen(t *testing.T) {

	actual := dbKeyGen("info", "param1", "param2")
	expected := "info/param1/param2"

	if actual != expected {
		t.Errorf("Test failed. Expected: %s, but got: %s", expected, actual)
	}
}

func TestGenerateLabelString(t *testing.T) {
	m := make(map[string]*string)
	value := "v"
	m["a"] = &value
	m["b"] = &value

	actual := generateLabelString(m)
	expected := "a:v;b:v"

	if actual != expected {
		t.Errorf("Test failed. Expected: %s, but got: %s", expected, actual)
	}
}

type idTest struct {
	id, expect string
}

var idTests = []idTest{
	{"1234", "----"},
	{"123", "------------------------------------"},
	{"1235", "------------------------------------"},
}

func TestCreateId(t *testing.T) {
	for i := range idTests {
		ex := CreateId(idTests[i].id)
		if len(ex) != len(idTests[i].expect) {
			t.Errorf("Test failed. Expected: %d, but got: %d", len(idTests[i].expect), len(ex))
		}
	}
}
