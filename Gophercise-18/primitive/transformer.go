package primitive

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

//Mode is mode of image transformation
type Mode int

const (
	//ModeCombo is prmitive mode
	ModeCombo Mode = iota
	//ModeTriangle is prmitive mode
	ModeTriangle
	//ModeRect is prmitive mode
	ModeRect
	//ModeEllipse is prmitive mode
	ModeEllipse
	//ModeCircle is prmitive mode
	ModeCircle
	//ModeRotatedRect is prmitive mode
	ModeRotatedRect
	//ModeBeziers is prmitive mode
	ModeBeziers
	//ModeRotatedEllipse is prmitive mode
	ModeRotatedEllipse
	//ModePolygon is prmitive mode
	ModePolygon
)

var (
	myPrimitive = primitive
	myTempFile  = TempFile
	myCopy      = io.Copy
)

//WithMode returs the function, the underlying function return the diffrent mode options
func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

//Transform is mechanism wich accept image and transform it using defined mode and shapes
func Transform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	var args []string
	for _, opt := range opts {
		args = append(args, opt()...)
	}
	//Create the file to store input io reader
	in, err := myTempFile("in_", ext, "")
	if err == nil {
		//Create the file to store transformed image
		defer os.Remove(in.Name())
		//Creating temp file to store out image
		out, err := myTempFile("out_", ext, "")
		defer os.Remove(out.Name())
		if err == nil {
			//Read image into file
			_, err = myCopy(in, image)
			if err != nil {
				return nil, errors.New("primitive: Error while reading image into file")
			}

			//Transform image using primitive api
			stdOutput, err := myPrimitive(in.Name(), out.Name(), Mode(numShapes), ModeCombo)
			if err != nil {
				return nil, errors.New("primitive: Error in transform image")
			}
			fmt.Println(stdOutput)
			b := bytes.NewBuffer(nil)
			//Coping the content of image output buffer
			_, err = myCopy(b, out)
			if err == nil {
				return b, nil
			}
		}
	}
	return nil, errors.New("primitive: Error in creating input temp file")

}

//primitive preovides the functionality to ececute all commands in prmitive package
func primitive(inputFile, outputFile string, numShapes, mode Mode) (string, error) {
	argsStr := fmt.Sprintf("-i %s -o %s -n %d -m %d", inputFile, outputFile, numShapes, mode)
	cmd := exec.Command("primitive", strings.Fields(argsStr)...)
	b, err := cmd.CombinedOutput()
	return string(b), err
}

//TempFile creates the file on disk and return its pointer
func TempFile(prefix, ext, dir string) (*os.File, error) {
	in, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return nil, errors.New("primitive: Error while creating temp file")
	}
	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}
