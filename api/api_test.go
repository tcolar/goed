package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/ui"
)

const test_port = 4241

func init() {
	core.Testing = true
	core.InitHome()
	core.Ed = ui.NewMockEditor()
	core.Ed.Start("../test_data/file1.txt")
}

func TestApiGets(t *testing.T) {
	api := Api{}
	go api.Start(test_port)

	body, err := get("/foobar")
	assert.NotNil(t, err)

	body, err = get("/api_version")
	assert.Nil(t, err)
	assert.Equal(t, core.ApiVersion, body, "api_version")

	body, err = get("/v1/version")
	assert.Nil(t, err)
	assert.Equal(t, core.Version, body, "version")

	body, err = get("/v1/cur_view")
	assert.Nil(t, err)
	assert.Equal(t, body, "1")

	body, err = get("/v1/view/1/title")
	assert.Nil(t, err)
	assert.Equal(t, body, "file1.txt")

	body, err = get("/v1/view/1/workdir")
	assert.Nil(t, err)
	d, _ := filepath.Abs("../test_data")
	assert.Equal(t, body, d)

	body, err = get("/v1/view/1/src_loc")
	p, _ := filepath.Abs("../test_data/file1.txt")
	assert.Nil(t, err)
	assert.Equal(t, body, p)

	body, err = get("/v1/view/1/dirty")
	assert.Nil(t, err)
	assert.Equal(t, body, "0")

	body, err = get("/v1/view/1/selections")
	assert.Nil(t, err)
	assert.Equal(t, body, "")

	s := core.Ed.CurView().Selections()
	sel := core.NewSelection(1, 1, 2, 10)
	sel2 := core.NewSelection(3, 3, 5, 5)
	*s = append(*s, *sel, *sel2)
	body, err = get("/v1/view/1/selections")
	assert.Nil(t, err)
	assert.Equal(t, body, "1 1 2 10\n3 3 5 5\n")

	body, err = get("/v1/view/1/line_count")
	assert.Nil(t, err)
	assert.Equal(t, body, "12")
}

func get(url string) (string, error) {
	response, err := http.Get(fmt.Sprintf("http://localhost:%d%s", test_port, url))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return string(body), fmt.Errorf("Error %d", response.StatusCode)
	}
	return string(body), nil
}
