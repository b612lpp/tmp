package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var item []string
	err := dirTree2(path, printFiles, &item, out)

	return err

}

func dirTree2(path string, printFiles bool, item *[]string, out io.Writer) error {
	var last bool

	f, err := os.Open(path)
	ss, _ := sortS(f, printFiles)

	for i := range ss {
		if i == len(ss)-1 {
			last = true
		}
		//fmt.Print(path, item)
		printItem(ss[i], formTab(path, last, item), out)
		dirTree2(path+"/"+ss[i].Name(), printFiles, item, out)

	}

	return err
}

func formTab(p string, last bool, item *[]string) string {
	var (
		pref string
		buf  []string
	)

	buf = append(buf, *item...)

	strings.ReplaceAll(p, "./", "")
	parts := strings.SplitAfter(p, "/")

	*item = []string{}
	if last == true {
		pref = "└───"
	} else {
		pref = "├───"
	}

	switch diff := len(parts) - len(buf); {
	case diff > 0:
		for i := range buf {
			if buf[i] == "├───" || buf[i] == "│	" {
				*item = append(*item, "│	")
			}
			if buf[i] == "└───" || buf[i] == "	" {
				*item = append(*item, "	")
			}
		}
		*item = append(*item, pref)
	case diff < 0:
		for i := 0; i < len(parts)-1; i++ {
			if buf[i] == "├───" || buf[i] == "│	" {
				*item = append(*item, "│	")
			}
			if buf[i] == "└───" || buf[i] == "	" {
				*item = append(*item, "	")
			}
		}
		*item = append(*item, pref)
	case diff == 0:
		for i := 0; i < len(parts)-1; i++ {
			if buf[i] == "├───" || buf[i] == "│	" {
				*item = append(*item, "│	")
			}
			if buf[i] == "└───" || buf[i] == "	" {
				*item = append(*item, "	")
			}
		}
		*item = append(*item, pref)
	}

	return strings.Join(*item, "")
}

func printItem(f os.DirEntry, h string, out io.Writer) {
	var value string
	if f.IsDir() {
		//hj := strings.Join(h, "")
		fmt.Fprintf(out, "%v%s\n", h, f.Name())
	} else {
		i, _ := f.Info()

		if i.Size() == 0 {
			value = "empty"
		} else {
			value = strconv.FormatInt(i.Size(), 10) + "b"
		}
		//hj := strings.Join(h, "")
		fmt.Fprintf(out, "%v%s (%v)\n", h, f.Name(), value)
	}
}

func sortS(f *os.File, fv bool) ([]os.DirEntry, string) {
	var lastF string
	var swfFile, swfFolder []os.DirEntry
	ff, _ := f.ReadDir(-1)
	for i := range ff {
		if ff[i].IsDir() {
			swfFolder = append(swfFolder, ff[i])
			sort.Slice(swfFolder, func(i, j int) bool {
				return swfFolder[i].Name() < swfFolder[j].Name()
			})
		} else {
			swfFile = append(swfFile, ff[i])
			sort.Slice(swfFile, func(i, j int) bool {
				return swfFile[i].Name() < swfFile[j].Name()
			})
		}

	}
	if !fv {
		if len(swfFolder) > 0 {
			lastF = swfFolder[len(swfFolder)-1].Name()
			return swfFolder, lastF
		}
	} else {
		if len(swfFolder) > 0 {
			lastF = swfFolder[len(swfFolder)-1].Name()

		}
		swfFolder = append(swfFolder, swfFile...)

	}
	return swfFolder, lastF
}
