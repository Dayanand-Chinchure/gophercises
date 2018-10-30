package primitive

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.ibm.com/dash/dash_utils/dashtest"
)

func mockTempfile(prefix, ext, dir string) (*os.File, error) {
	return nil, errors.New("Error creating temp file")
}

func mockCopy(dst io.Writer, src io.Reader) (int64, error) {
	return 0, errors.New("Error Copying content")
}

func mockPrimitive(inputFile, outputFile string, numShapes, mode Mode) (string, error) {
	return "", errors.New("Error in primitive")
}

func TestTransform(t *testing.T) {
	file, _ := os.Open("../img/volverine.jpg")
	reader, err := Transform(file, "jpg", 10, WithMode(ModeCombo))
	if err != nil || reader == nil {
		t.Error("TestTransform: " + err.Error())
		return
	}
}

func TestTransformWithErrorInTempFile(t *testing.T) {
	oldTempFile := myTempFile
	myTempFile = mockTempfile
	file, _ := os.Open("../img/volverine.jpg")
	_, err := Transform(file, "jpg", 10, WithMode(ModeCombo))
	if err == nil {
		t.Error("TestTransform:  Failed temp file code")
	}
	myTempFile = oldTempFile
}

func TestTransformWithErrorInCopy(t *testing.T) {
	oldCopy := myCopy
	myCopy = mockCopy
	file, _ := os.Open("../img/volverine.jpg")
	_, err := Transform(file, "jpg", 10, WithMode(ModeCombo))
	if err == nil {
		t.Error("TestTransform:  Failed to copy the content")
	}
	myCopy = oldCopy
}

func TestTransformWithErrorInPrimitive(t *testing.T) {
	oldPrimitive := myPrimitive
	myPrimitive = mockPrimitive
	file, _ := os.Open("../img/volverine.jpg")
	_, err := Transform(file, "jpg", 10, WithMode(ModeCombo))
	if err == nil {
		t.Error("TestTransform:  Failed to copy the content")
	}
	myPrimitive = oldPrimitive
}

func Test_primitive(t *testing.T) {
	os.MkdirAll("./temp", os.ModePerm)
	out, err := primitive("../img/volverine.jpg", "./temp/out.jpg", 10, ModeCombo)
	if err != nil || strings.TrimSpace(out) != "" {
		t.Error("Failed to trasform the image ", out)
		return
	}
	if _, err := os.Open("./temp/out.jpg"); err != nil {
		t.Error("Failed to trasform the image ")
		return
	}
	os.RemoveAll("./temp")

}

func Test_PrimitiveError(t *testing.T) {
	out, err := primitive("../img/", "./temp/out.jpg", 10, ModeCombo)
	if err == nil || strings.TrimSpace(out) == "" {
		t.Error("Failed to trasform the image ", out)
		return
	}

}

func TestTempFile(t *testing.T) {
	os.MkdirAll("./test", os.ModePerm)
	got, err := TempFile("", "txt", "./test")
	if got == nil {
		t.Errorf("TempFile() error = %v, wantErr %v", err != nil, false)
		return
	}
	defer os.RemoveAll("./test")
}

func TestTempFileWithInvalidDir(t *testing.T) {
	got, err := TempFile("", "txt", "./test")
	if got != nil {
		t.Errorf("TempFile() error = %v, wantErr %v", err != nil, true)
		return
	}
}

func TestMain(m *testing.M) {
	dashtest.ControlCoverage(m)
}
