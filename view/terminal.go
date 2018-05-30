package view

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	dbgPrefix             = "(xdbg)"
	defaultCodeLineLength = 10
)

// SourceFileCache caches all the source files
type SourceFileCache struct {
	Cache map[string][]string
}

var fileCache = NewSourceFileCache()

func (vfc *SourceFileCache) cacheFile(s string) error {
	fh, err := os.Open(s)

	if err != nil {
		return err
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	res := []string{}

	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	vfc.Cache[s] = res
	return nil
}

func (vfc *SourceFileCache) getLines(s string, begin int, length int) ([]string, error) {
	ce, ok := vfc.Cache[s]
	if !ok {
		if err := vfc.cacheFile(s); err != nil {
			return nil, err
		}
		ce, _ = vfc.Cache[s]
	}

	if begin >= len(ce) || begin < 0 {
		return nil,
			fmt.Errorf(
				"line number %d out of range; %s has %d lines",
				begin,
				filepath.Base(s),
				len(ce))
	}

	if (begin + length) > len(ce) {
		length = len(ce) - begin
	}

	return ce[begin:(begin + length)], nil
}

// NewSourceFileCache creates a new instance
func NewSourceFileCache() *SourceFileCache {
	return &SourceFileCache{
		Cache: map[string][]string{},
	}
}

// View facade
type View struct {
	Input *bufio.Reader
}

// NewView creates a new view instance
func NewView() *View {
	return &View{
		Input: bufio.NewReader(os.Stdin),
	}
}

// PrintInputPrefix prints the input preamble
func (v *View) PrintInputPrefix() {
	fmt.Printf("%s ", dbgPrefix)
}

// Print prints a string without any format
func (v *View) Print(s string) {
	fmt.Print(s)
}

// PrintLn prints a string without any format
func (v *View) PrintLn(s string) {
	fmt.Println(s)
}

// PrintErrorLn prints a error string
func (v *View) PrintErrorLn(s string) {
	fmt.Println(s)
}

// GetInputLine reads one line
func (v *View) GetInputLine() string {
	l, _ := v.Input.ReadString('\n')
	return l
}

// PrintSourceLn prints a string without any format
func (v *View) PrintSourceLn(s string, idx int, length int) {
	if !strings.HasPrefix(s, "file://") {
		v.PrintErrorLn("Unknown source type " + s)
		return
	}

	s = strings.TrimPrefix(s, "file://")
	lines, err := fileCache.getLines(s, idx-1, length)
	if err != nil {
		v.PrintErrorLn("Error opening file " + err.Error())
		return
	}

	padding := len(strconv.Itoa(idx + length))
	fmtLine := "%-" + strconv.Itoa(padding) + "d\t%s"

	for i := 0; i < len(lines); i++ {
		v.PrintLn(fmt.Sprintf(fmtLine, idx+i, lines[i]))
	}
}

// PrintApplicationInformation general dialog
func (v *View) PrintApplicationInformation(version string, host string, port uint16) {
	v.PrintLn(fmt.Sprintf("xdebug-cli version %s", version))
	v.PrintLn(fmt.Sprintf("listening on %s:%d for incoming connections", host, port))
	v.PrintLn("xdebug dbgp https://xdebug.org/docs-dbgp.php")
	v.PrintLn("bug reports to https://github.com/jami/xdebug-cli/issues")
	v.PrintLn("feel free to contribute!")
}
