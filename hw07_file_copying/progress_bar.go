package main

import (
	"fmt"
	"io"
	"time"
)

type ProgressBar struct {
	total      int64
	position   int64
	percentage int64
	processed  int64
	start      time.Time
	placed     bool
}

func NewProgressBar(t int64) *ProgressBar {
	pb := new(ProgressBar)
	pb.SetTotal(t)
	return pb
}

func (pb *ProgressBar) Place() {
	if pb.placed {
		return
	}
	pb.placed = true
	pb.start = time.Now()
	fmt.Printf("Start: 0%%")
	pb.render()
}

func (pb *ProgressBar) SetTotal(t int64) {
	pb.total = t
}

func (pb *ProgressBar) Add(i int) {
	pb.Add64(int64(i))
}

func (pb *ProgressBar) Add64(i int64) {
	pb.position += i
	pb.percentage = pb.position * 100 / pb.total

	if i != 0 {
		pb.render()
	}
}

func (pb *ProgressBar) Finish() {
	pb.position = pb.total
	pb.percentage = 100
	pb.render()
}

func (pb *ProgressBar) render() {
	if !pb.placed {
		return
	}

	if pb.processed < pb.percentage {
		for pb.processed < pb.percentage {
			pb.processed++
			pb.print(pb.processed)
		}
	}
}

func (pb *ProgressBar) print(v int64) {
	switch {
	case v == 0:
		fmt.Printf("0%%")
	case v == 100:

		fmt.Printf("100%%. ")
		elapsedTime := time.Since(pb.start)
		if elapsedTime%time.Second != 0 {
			elapsedTime = elapsedTime - (elapsedTime % time.Second) + time.Second
		}

		fmt.Printf("Done %v.\n", elapsedTime.String())
	case v%10 == 0:
		fmt.Printf("%v%%", v)
	case v%3 == 0:
		fmt.Printf(".")
	}
}

func (pb *ProgressBar) Reader(r io.Reader) *Reader {
	return &Reader{r, pb}
}
