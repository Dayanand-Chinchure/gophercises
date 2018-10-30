package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Dayanand-Chinchure/gophercises/transform/primitive"
)

var (
	myGenImages   = genImages
	myIOCopy      = io.Copy
	myPrimTmpFile = primitive.TempFile
)

func main() {
	mux := http.NewServeMux()
	//Root of the server which present the upload image form to user
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<html><body>
		<form action="/upload"
			enctype="multipart/form-data" method="post">
		<p>
			Please specify a file:<br>
			<input type="file" name="file" size="40">
		</p>	
		<div>
		<input type="submit" value="Send">
		</div>
		</form>
		</body></html>
		`
		fmt.Fprint(w, html)
	})
	//This path provides the functionity to upload the image on server
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		ext := filepath.Ext(fileHeader.Filename)[1:]
		onDiskFile, err := myPrimTmpFile("in_", ext, "./img/")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer onDiskFile.Close()
		_, err = myIOCopy(onDiskFile, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/render/"+filepath.Base(onDiskFile.Name()), http.StatusFound)
	})
	//Renders the diffrent versions of transformed images on users browser
	mux.HandleFunc("/render/", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("./img/" + filepath.Base(r.URL.Path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		ext := filepath.Ext(file.Name())[1:]
		modeStr := r.FormValue("mode")
		//Checks if mode is provided, if not render the images with diffrent mode combinations
		if modeStr == "" {
			renderChoicesWithModes(w, r, file, ext)
			return
		}
		mode, err := strconv.Atoi(modeStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		nStr := r.FormValue("shapes")
		//Checks if shapes is provided, if not render the images with diffrent shapes value combinations
		if nStr == "" {
			renderWithShapes(w, r, file, ext, primitive.Mode(mode))
			return
		}
		numShapes, err := strconv.Atoi(nStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = numShapes
		http.Redirect(w, r, "/img/"+filepath.Base(file.Name()), http.StatusFound)
	})
	//Creates file server
	fs := http.FileServer(http.Dir("./img/"))
	mux.Handle("/img/", http.StripPrefix("/img", fs))
	http.ListenAndServe(":3000", mux)

}

//Options data struct for mode and mmo of shapes
type Options struct {
	mode   primitive.Mode
	shapes int
}

//renderChoicesWithModes creates multiple versions of image by applying combination of various number
//of modes with same shape and render them on user browser
func renderChoicesWithModes(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, ext string) {
	opts := []Options{
		{mode: primitive.ModeCombo, shapes: 10},
		{mode: primitive.ModeEllipse, shapes: 10},
		{mode: primitive.ModePolygon, shapes: 10},
		{mode: primitive.ModeBeziers, shapes: 10},
	}
	images, err := myGenImages(rs, ext, opts...)
	if err != nil {
		//panic(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
			{{range .}}
				<a href="/render/{{.Name}}?mode={{.Mode}}">
					<img style="width: 20%;" src="/img/{{.Name}}">
				</a>
			{{end}}
			</body></html>`
	mytemplate := template.Must(template.New("").Parse(html))

	type placeholder struct {
		Mode primitive.Mode
		Name string
	}
	var placeholders []placeholder
	for i, img := range images {
		placeholders = append(placeholders, placeholder{Name: filepath.Base(img), Mode: opts[i].mode})
	}
	mytemplate.Execute(w, placeholders)

}

//renderWithShapes creates multiple versions of image by applying combination of various number
//of shapes with same mode and render them on user browser
func renderWithShapes(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, ext string, mode primitive.Mode) {
	opts := []Options{
		{mode: mode, shapes: 11},
		{mode: mode, shapes: 22},
		{mode: mode, shapes: 33},
		{mode: mode, shapes: 44},
	}
	images, err := myGenImages(rs, ext, opts...)
	if err != nil {
		//panic(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
			{{range .}}
				<a href="/render/{{.Name}}?mode={{.Mode}}&shapes={{.Shapes}}">
					<img style="width: 20%;" src="/img/{{.Name}}">
				</a>
			{{end}}
			</body></html>`
	mytemplate := template.Must(template.New("").Parse(html))

	type placeholder struct {
		Mode   primitive.Mode
		Name   string
		Shapes int
	}
	var placeholders []placeholder
	for i, img := range images {
		placeholders = append(placeholders, placeholder{Name: filepath.Base(img), Mode: opts[i].mode, Shapes: opts[i].shapes})
	}
	mytemplate.Execute(w, placeholders)

}

//genImsges provide the fuctionality to transform multiple images
func genImages(rs io.ReadSeeker, ext string, opts ...Options) ([]string, error) {
	var images []string
	for _, opt := range opts {
		rs.Seek(0, 0)
		img, err := genImage(rs, ext, opt.shapes, opt.mode)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

//genImage transform incoming image into some out file and return its path
func genImage(r io.Reader, ext string, numShapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.Transform(r, ext, numShapes, primitive.WithMode(mode))
	if err == nil {
		outFile, err := myPrimTmpFile("out_", ext, "./img/")
		if err == nil {
			defer outFile.Close()
			myIOCopy(outFile, out)
			return outFile.Name(), nil
		}
	}
	return "", err
}
