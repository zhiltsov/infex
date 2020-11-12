package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type ExtMap map[string]int64

type Pair struct {
	First  string
	Second int64
}

type PairList []Pair

func (pl PairList) Len() int {
	return len(pl)
}

func (pl PairList) Less(i, j int) bool {
	return pl[i].Second > pl[j].Second
}

func (pl PairList) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

func (pl PairList) String() (s string) {
	if len(pl) > 0 {
		for _, p := range pl {
			s += fmt.Sprintf("[%s]: %s\n", p.First, byteCount(p.Second))
		}
	}
	return
}

type Result struct {
	count   int64
	allSize int64
	extSize ExtMap
}

func (r *Result) plus(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		if ext := filepath.Ext(info.Name()); len(ext) > 0 {
			if _, ok := r.extSize[ext]; !ok {
				r.extSize[ext] = 0
			}
			r.allSize += info.Size()
			r.extSize[ext] += info.Size()
			r.count++
		}
	}

	return nil
}

func (r *Result) sort() PairList {
	pl := make(PairList, len(r.extSize))
	i := 0
	for k, v := range r.extSize {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(&pl)

	return pl
}

func (r *Result) String() (s string) {
	s = fmt.Sprintf("File count: %d\n", r.count)
	if r.count > 0 {
		line := "----------------------\n"
		s += line + fmt.Sprintf("Extension count: %d\n", len(r.extSize))
		s += fmt.Sprintf("Summary size: %s\n", byteCount(r.allSize)) + line + "\n"
	}
	return
}

func byteCount(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func nullOrExit(err error) {
	if err != nil {
		println(err)
		os.Exit(0)
	}
}

func main() {
	if len(os.Args) < 2 || len(os.Args[1]) == 0 {
		nullOrExit(errors.New("do not specify a directory path"))
	}
	rs := &Result{extSize: make(ExtMap)}
	nullOrExit(filepath.Walk(os.Args[1], rs.plus))
	print(rs.String(), rs.sort().String())
}
