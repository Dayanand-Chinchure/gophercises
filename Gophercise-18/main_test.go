package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Dayanand-Chinchure/gophercises/transform/primitive"
	"github.ibm.com/dash/dash_utils/dashtest"
)

func mockGenImages(rs io.ReadSeeker, ext string, opts ...Options) ([]string, error) {
	return nil, errors.New("Error Occured")
}
func mockIOCopy(w io.Writer, r io.Reader) (int64, error) {
	return -1, errors.New("Error Occured")
}

func mockPrimTmpFile(prefix, ext, dir string) (*os.File, error) {
	return nil, errors.New("Error Occured")
}
func TestGenImage(t *testing.T) {
	input, err := os.Open("./img/go.png")
	if err != nil {
		t.Errorf("TestGenImage: Failed to open input file")
	}
	outfile, err := genImage(input, "png", 10, primitive.ModeCombo)
	if err != nil {
		t.Errorf("TestGenImage: Failed to generate image")
	}
	if outfile == "" {
		t.Errorf("TestGenImage: Failed to generate image")
	}
}

func TestGenImageFail(t *testing.T) {
	var file *os.File
	outfile, err := genImage(file, "png", 10, primitive.ModeCombo)
	if err == nil || outfile != "" {
		t.Errorf("TestGenImageFail: Failed to generate image")
	}
}

func TestGenImages(t *testing.T) {
	input, err := os.Open("./img/go.png")
	if err != nil {
		t.Errorf("TestGenImages: Failed to open input file")
	}
	outs, err := genImages(input, "png", Options{shapes: 10, mode: primitive.ModeCombo})
	if err != nil || len(outs) < 1 {
		t.Errorf("TestGenImages: Failed to generate images")
	}

}

func TestGenImagesFail(t *testing.T) {
	var input *os.File
	outs, err := genImages(input, "png", Options{shapes: 10, mode: primitive.ModeCombo})
	if err == nil || len(outs) > 0 {
		t.Errorf("TestGenImagesFail: Failed to generate images")
	}

}

func fileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func TestRoot(t *testing.T) {
	go main()
	time.Sleep(10)
	resp, err := followURL("GET", "http://localhost:3000/")
	if err != nil {
		t.Errorf("TestUpload: Failed to get upload form")
	}
	if resp.StatusCode != 200 {
		t.Errorf("TestUpload: Failed to get upload form with response code = %d", resp.StatusCode)
	}
}

func TestRootFail(t *testing.T) {
	go main()
	time.Sleep(10)
	resp, err := followURL("GET", "http://localhost:3001/")
	if err == nil {
		t.Errorf("TestUpload: Failed to get upload form")
	}
	if resp != nil {
		t.Errorf("TestUpload: Failed to get upload form")
	}
}

func TestUploadBadRequest(t *testing.T) {
	go main()
	time.Sleep(10)
	resp, err := followURL("POST", "http://localhost:3000/upload")
	if err != nil {
		t.Errorf("TestUploadBadRequest: Failed to get upload form")
	}
	if resp.StatusCode != 400 {
		t.Errorf("TestUploadBadRequest: Failed to get upload form with response code = %d", resp.StatusCode)
	}
}

func TestUpload(t *testing.T) {
	go main()
	time.Sleep(10)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	b, w, err := createUploadBody()
	req, err := http.NewRequest("POST", "http://localhost:3000/upload", &b)
	if err != nil {
		t.Errorf("TestUpload: " + err.Error())
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("TestUpload: Failed to upload image")
		return
	}
	if res.StatusCode != 302 {
		t.Errorf("TestUpload: Failed to upload image with status code = %d", res.StatusCode)
		return
	}
}

func TestUploadErrorInIOCopy(t *testing.T) {
	oldIOCopy := myIOCopy
	myIOCopy = mockIOCopy
	go main()
	time.Sleep(10)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	b, w, err := createUploadBody()
	req, err := http.NewRequest("POST", "http://localhost:3000/upload", &b)
	if err != nil {
		t.Errorf("TestUploadErrorInIOCopy: " + err.Error())
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("TestUploadErrorInIOCopy: Failed to upload image")
		return
	}
	if res.StatusCode != 500 {
		t.Errorf("TestUploadErrorInIOCopy: Failed to upload image with status code = %d", res.StatusCode)
		return
	}
	myIOCopy = oldIOCopy
}

func TestUploadErrorPrim(t *testing.T) {
	oldPrimTmpFile := myPrimTmpFile
	myPrimTmpFile = mockPrimTmpFile
	go main()
	time.Sleep(10)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	b, w, err := createUploadBody()
	req, err := http.NewRequest("POST", "http://localhost:3000/upload", &b)
	if err != nil {
		t.Errorf("TestUploadErrorInIOCopy: " + err.Error())
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("TestUploadErrorInIOCopy: Failed to upload image")
		return
	}
	if res.StatusCode != 500 {
		t.Errorf("TestUploadErrorInIOCopy: Failed to upload image with status code = %d", res.StatusCode)
		return
	}
	myPrimTmpFile = oldPrimTmpFile
}

func TestRender(t *testing.T) {
	go main()
	time.Sleep(50)
	res, err := followURL("GET", "http://localhost:3000/render/img/go.png")
	if err != nil || res == nil {
		t.Errorf("TestRender: " + err.Error())
		return
	}
}

func TestRenderWithMode(t *testing.T) {
	go main()
	time.Sleep(50)
	res, err := followURL("GET", "http://localhost:3000/render/img/go.png?mode=1")
	if err != nil || res == nil {
		t.Errorf("TestRender: " + err.Error())
		return
	}
}

func TestRenderWithModeAndShapes(t *testing.T) {
	go main()
	time.Sleep(10)
	res, err := followURL("GET", "http://localhost:3000/render/img/go.png?mode=1&shapes=10")
	if err != nil || res == nil {
		t.Errorf("TestRenderWithModeAndShapes: " + err.Error())
		return
	}
}

func TestRenderWithBadFilePath(t *testing.T) {
	go main()
	time.Sleep(10)
	res, _ := followURL("GET", "http://localhost:3000/render/img/got.png")
	if res.StatusCode != 500 {
		t.Errorf("TestRenderWithBadFilePath: Failed with error code =%d", res.StatusCode)
		return
	}
}

func TestRenderWithBadMode(t *testing.T) {
	go main()
	time.Sleep(10)
	res, _ := followURL("GET", "http://localhost:3000/render/img/go.png?mode=abs")
	if res.StatusCode != 400 {
		t.Errorf("TestRenderWithBadMode: Failed with error code =%d", res.StatusCode)
		return
	}
}

func TestRenderWithBadShapes(t *testing.T) {
	go main()
	time.Sleep(10)
	res, _ := followURL("GET", "http://localhost:3000/render/img/go.png?mode=1&shapes=abs")
	if res.StatusCode != 400 {
		t.Errorf("TestRenderWithBadShapes: Failed with error code =%d", res.StatusCode)
		return
	}
}

func TestRenderErrorInGenImages1(t *testing.T) {
	oldGenImages := myGenImages
	myGenImages = mockGenImages
	go main()
	time.Sleep(10)
	res, _ := followURL("GET", "http://localhost:3000/render/img/go.png")
	if res.StatusCode != 500 {
		t.Errorf("TestRenderErrorInGenImages1: Failed with error code =%d", res.StatusCode)
		return
	}
	myGenImages = oldGenImages
}

func TestRenderErrorInGenImages2(t *testing.T) {
	oldGenImages := myGenImages
	myGenImages = mockGenImages
	go main()
	time.Sleep(10)
	res, _ := followURL("GET", "http://localhost:3000/render/img/go.png?mode=1")
	if res.StatusCode != 500 {
		t.Errorf("TestRenderErrorInGenImages2: Failed with error code =%d", res.StatusCode)
		return
	}
	myGenImages = oldGenImages
}

func followURL(method, path string) (*http.Response, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	var req *http.Request
	var resp *http.Response
	var err error
	req, _ = http.NewRequest(method, path, nil)
	resp, err = client.Do(req)
	fmt.Printf("Request is %v\n", req)
	fmt.Printf("Response is %v\n", resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func createUploadBody() (bytes.Buffer, multipart.Writer, error) {
	f, _ := os.Open("./img/go.png")
	values := map[string]io.Reader{
		"file": f,
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return bytes.Buffer{}, *w, err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return bytes.Buffer{}, *w, err
		}

	}
	w.Close()
	return b, *w, nil
}

func TestMain(m *testing.M) {
	dashtest.ControlCoverage(m)
}
