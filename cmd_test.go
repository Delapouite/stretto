package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretto-editor/gocui"
)

func TestValidateCmd(t *testing.T) {
	g := initGui()
	defer g.Close()

	// unauthorized calls : not from the cmdline
	v, _ := g.View("main")
	assert.Panics(t, func() { validateCmd(g, v) }, "Cmdline is not the current view")
}

func TestUnknownCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vMain, _ := g.View("main")
	writeInView(v, "kl,sflk,f")
	validateCmd(g, v)
	assert.Contains(t, vMain.Footer, ErrUnknownCommand.Error(), "unknown command error expected")
}

func TestEmptyCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	validateCmd(g, v)
}

func TestOpenCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vMain, _ := g.View("main")
	filename := "Commands.md"
	writeInView(v, "o "+filename)
	validateCmd(g, v)
	f, err := os.Open(filename)
	assert.Nil(t, err, err)
	content, _ := ioutil.ReadAll(f)
	assert.Equal(t, string(content)+"\n", vMain.Buffer(), "vMain should contains the content of "+filename)

	v.EditWrite('o')
	validateCmd(g, v)
	assert.Contains(t, vMain.Footer, ErrMissingFilename.Error(), "missing argument error expected")

	writeInView(v, "o "+filename+" useless args")
	validateCmd(g, v)
	assert.Contains(t, vMain.Footer, ErrUnexpectedArgument.Error(), "unexpected argument error expected")
}

func TestCloseCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vMain, _ := g.View("main")
	writeInView(v, "c!")
	validateCmd(g, v)
	assert.Equal(t, "undefined", vMain.Title, "Title of the main view should be undefined")
	assert.Equal(t, "", vMain.Buffer(), "The buffer of the main view should be undefined")

	//TODO : add an error when there is an unexpected argument
}

func TestReplaceAllCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vMain, _ := g.View("main")
	// this test doesn't work because viewlines are empty so EditDelete doen't work
	/*	text := " foo foo foo \n foo \n \n foo"
		expected := " bar bar bar \n bar \n \n bar\n"
		writeInView(vMain, text)
		//need to fill the viewlines
		writeInView(v, "repall foo bar")
		validateCmd(g, v)

		assert.Equal(t, expected, vMain.Buffer(), "all the words equal to the pattern should be replaced")
	*/
	//without arguments
	writeInView(v, "repall")
	validateCmd(g, v)
	assert.Contains(t, vMain.Footer, ErrMissingPattern.Error(), "missing pattern or replacement for search/replace")

	//with too many arguments
	writeInView(v, "repall 1 2 3")
	validateCmd(g, v)
	assert.Contains(t, vMain.Footer, ErrUnexpectedArgument.Error(), "unexpected third argument")
}

func TestSetWrapCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vMain, _ := g.View("main")

	writeInView(v, "setwrap true")
	validateCmd(g, v)
	assert.Equal(t, true, vMain.Wrap, "wrap should be true")

	writeInView(v, "setwrap false")
	validateCmd(g, v)
	assert.Equal(t, false, vMain.Wrap, "wrap sould be false")

	writeInView(v, "setwrap ")
	validateCmd(g, v)
	assert.Contains(t, vMain.Footer, ErrWrapArgument.Error(), "missing argument error")

	writeInView(v, "setwrap useless args")
	validateCmd(g, v)
	assert.Contains(t, vMain.Footer, ErrUnexpectedArgument.Error(), "unexpected argument error")
}

func TestQuitCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	writeInView(v, "q!")
	err := validateCmd(g, v)
	assert.EqualError(t, err, gocui.ErrQuit.Error(), "Errquit should be returned from validatecmd when q! is executed")
}

func TestQuitAndSaveCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	writeInView(v, "qs")
	currentFile = "6u8Y73wHm5QWmgRPcXk96y39cL.txt"

}

func TestSaveAsCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vMain, _ := g.View("main")
	text := "This is a \n test on two lines"
	filename := "r9w92W2Cn7MTtAhuCP5si2LH356r8FrjV.txt"
	// First create the file
	writeInView(vMain, text)
	writeInView(v, "sa "+filename)
	validateCmd(g, v)
	assert.Equal(t, text, getContentFile(filename), "the save file doesn't contain the right content")

	text2 := " another text"
	text += text2
	writeInView(vMain, text2)
	// Then write in the existing file
	writeInView(v, "sa "+filename)
	validateCmd(g, v)
	assert.Equal(t, text, getContentFile(filename), "the save file doesn't contain the right content")

	os.Remove(filename)

	writeInView(v, "sa ")
	validateCmd(g, v)
	assert.Contains(t, vMain.Footer, ErrMissingFilename.Error(), "missing filename error")

	writeInView(v, "sa useless args")
	validateCmd(g, v)
	assert.Contains(t, vMain.Footer, ErrUnexpectedArgument.Error(), "missing argument error")
}

func TestSaveAndCloseCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vMain, _ := g.View("main")
	currentFile = ""
	text := "This is a \n test on two lines"
	filename := "r9w92W2Cn7MTtAhuCP5si2LH356r8FrjV.txt"
	//save with no current file
	writeInView(vMain, text)
	writeInView(v, "sc "+filename)
	validateCmd(g, v)
	assert.Equal(t, text, getContentFile(filename), "the save file doesn't contain the right content")
	assert.Equal(t, "", currentFile, "the current file name should be empty")
	assert.Equal(t, "", vMain.Buffer(), "the view should be empty")

	//save with a current file
	clearView(vMain)
	currentFile = filename
	text = "I'm trying to save \n and close an opened file"
	writeInView(vMain, text)
	writeInView(v, "sc ")
	validateCmd(g, v)
	assert.Equal(t, text, getContentFile(filename), "the save file doesn't contain the right content")
	assert.Equal(t, "", currentFile, "the current file name should be empty")
	assert.Equal(t, "", vMain.Buffer(), "the view should be empty")
	os.Remove(filename)

	//try to save without a current file name and without an argument
	currentFile = ""
	writeInView(v, "sc")
	validateCmd(g, v)
	assert.Contains(t, vMain.Footer, ErrMissingFilename.Error(), "missing filename error")
}

func TestSaveAndQuitCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vMain, _ := g.View("main")

	text := "This is a \n test on two lines"
	filename := "r9w92W2Cn7MTtAhuCP5si2LH356r8FrjV.txt"
	//save with no current file
	writeInView(vMain, text)
	writeInView(v, "sq "+filename)
	err := validateCmd(g, v)
	assert.Equal(t, text, getContentFile(filename), "the save file doesn't contain the right content")
	assert.EqualError(t, err, gocui.ErrQuit.Error(), "Errquit should be returned from validatecmd when sq is executed")
	os.Remove(filename)

	//try to save without a current file name and without an argument
	currentFile = ""
	writeInView(v, "sq")
	validateCmd(g, v)
	assert.Contains(t, vMain.Footer, ErrMissingFilename.Error(), "missing filename error ")
}

func writeInView(v *gocui.View, s string) {
	for _, c := range s {
		v.EditWrite(c)
	}
}

func getContentFile(filename string) string {
	f, _ := os.Open(filename)
	content, _ := ioutil.ReadAll(f)
	return string(content)
}
