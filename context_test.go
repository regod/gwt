package gwt

import (
	"net/http"
	"testing"
)

func TestContextQuery(t *testing.T) {
	app := New()
	req_path := "/ab/c/d?id=123&tag=xyz"
	r, _ := http.NewRequest("GET", req_path, nil)
	ctx := NewContext(app, nil, r)

	test_table := map[string]string{
		"id":  "123",
		"tag": "xyz",
	}
	for key, expect := range test_table {
		actual := ctx.Query().Get(key)
		if actual != expect {
			t.Errorf("actual: %s; expect: %s", actual, expect)
		}
	}
}

func TestContextStore(t *testing.T) {
	app := New()
	r, _ := http.NewRequest("GET", "/", nil)
	ctx := NewContext(app, nil, r)
	test_table := map[string]interface{}{
		"num":  123,
		"str":  "abc",
		"bool": true,
	}
	for key, val := range test_table {
		ctx.SetStore(key, val)
	}
	for key, expect := range test_table {
		actual := ctx.GetStore(key)
		if actual != expect {
			t.Errorf("actual: %s; expect: %s", actual, expect)
		}
	}

}
