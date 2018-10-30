package recovermw

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

//Middleware function
func Middleware(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, Links(string(stack)))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		app.ServeHTTP(w, r)
	}
}

//SourceCodeHandler display the source code on browser with specified line number
func SourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("filepath")
	lineno := r.URL.Query().Get("line")
	line, err := strconv.Atoi(lineno)
	var fileRead []byte
	if len(filePath) <= 0 && len(lineno) <= 0 && err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "filepath or line no provided is wrong")
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileRead, err = ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lexer := lexers.Get("go")
	iterator, err := lexer.Tokenise(nil, string(fileRead))
	style := styles.Get("dracula")
	formatter := html.New(html.TabWidth(2), html.WithLineNumbers(), html.HighlightLines([][2]int{{line, line}}))

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<style>pre { font-size: 1.2em; }</style>")
	formatter.Format(w, style, iterator)
}

//Links create link function
func Links(stack string) string {
	lines := strings.Split(stack, "\n")
	for li, line := range lines {
		if len(line) == 0 || line[0] != '\t' {
			continue
		}
		file := ""
		for i, ch := range line {
			if ch == ':' {
				file = line[1:i]
				break
			}
		}
		var lineStr strings.Builder
		for i := len(file) + 2; i < len(line); i++ {
			if line[i] < '0' || line[i] > '9' {
				break
			}
			lineStr.WriteByte(line[i])
		}
		v := url.Values{}
		v.Set("path", file)
		v.Set("line", lineStr.String())
		lines[li] = "\t<a href=\"/debug/?" + v.Encode() + "\">" + file + ":" + lineStr.String() + "</a>" + line[len(file)+2+len(lineStr.String()):]
	}
	return strings.Join(lines, "\n")
}

//Panic return the panic
func Panic(w http.ResponseWriter, r *http.Request) {
	panic("It's Panic !")
}
