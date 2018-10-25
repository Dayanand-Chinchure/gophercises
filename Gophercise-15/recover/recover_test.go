package recovermw

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.ibm.com/dash/dash_utils/dashtest"
)

func TestMain(m *testing.M) {
	dashtest.ControlCoverage(m)
}

func TestRenderFileWithNoQueryParams(t *testing.T) {
	handler := http.HandlerFunc(SourceCodeHandler)
	response, _ := executeRequest("Get", "/debug/", Middleware(handler))
	actualBody := response.Body.String()
	assert.Equal(t, response.Code, 400)
	assert.Contains(t, actualBody, "filepath or line no provided is wrong")
}

func TestRenderFileWithInvalidFilePath(t *testing.T) {
	handler := http.HandlerFunc(SourceCodeHandler)
	response, _ := executeRequest("Get", "/debug/?filepath=/test/test&line=20", Middleware(handler))
	actualBody := response.Body.String()
	assert.Equal(t, response.Code, 500)
	assert.Contains(t, actualBody, "no such file or directory")
}
func TestRenderFileWithInvalidNo(t *testing.T) {
	handler := http.HandlerFunc(SourceCodeHandler)
	response, _ := executeRequest("Get", "/debug/?filepath=/home/gs-1325/go/src/github.ibm.com/Dayanand-Chinchure/gophercises/recover_middlerware/main.go&line=test", Middleware(handler))
	assert.Equal(t, response.Code, 500)
}

func TestRenderFile(t *testing.T) {
	handler := http.HandlerFunc(SourceCodeHandler)
	response, _ := executeRequest("Get", "/debug/?filepath=/home/gs-1325/go/src/github.com/Dayanand-Chinchure/gophercises/recover_middleware/main.go&line=10", Middleware(handler))
	assert.Equal(t, response.Code, 200)
}

func TestMiddlewareHandler(t *testing.T) {
	handler := http.HandlerFunc(Panic)
	response, _ := executeRequest("Get", "/panic", Middleware(handler))
	actualBody := response.Body.String()
	assert.Contains(t, actualBody, "It's Panic !")

}

func executeRequest(method string, url string, handler http.Handler) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	recorder := httptest.NewRecorder()
	recorder.Result()
	handler.ServeHTTP(recorder, req)
	return recorder, err
}

func TestMakeLinks(t *testing.T) {
	input := `goroutine 10 [running]:
	runtime/debug.Stack(0xc000052ce8, 0x1, 0x1)
		/usr/local/go/src/runtime/debug/stack.go:24 +0xa7
	github.com/Dayanand-Chinchure/gophercises/recover_middleware/recover.Middleware.func1.1(0xa4e2e0, 0xc0004dee40)
		/home/gs-1325/go/src/github.com/Dayanand-Chinchure/gophercises/recover_middleware/recover/recover.go:25 +0xac`
	output := `<a href="/debug/?line=25&path=%09%2Fhome%2Fgs-1325%2Fgo%2Fsrc%2Fgithub.com%2FDayanand-Chinchure%2Fgophercises%2Frecover_middleware%2Frecover%2Frecover.go">	/home/gs-1325/go/src/github.com/Dayanand-Chinchure/gophercises/recover_middleware/recover/recover.go:25</a>`
	retVal := Links(input)
	t.Log(retVal)

	if !strings.Contains(retVal, output) {
		t.Errorf("FAIL: expected '%s' in string retured from Links. Not found", output)
	}

}
