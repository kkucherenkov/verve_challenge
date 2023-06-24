package file_processor

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"
)

type Loader interface {
	AddItem(text string) error
}

func ProcessFile(fileName string, l Loader) error {
	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println("not able to read the file", err)
		return nil
	}

	defer file.Close() //close after checking err

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println("Could not able to get the file stat")
		return err
	}

	fileSize := fileStat.Size()
	offset := fileSize - 1
	lastLineSize := 0

	for {
		b := make([]byte, 1)
		n, err := file.ReadAt(b, offset)
		if err != nil {
			fmt.Println("Error reading file ", err)
			break
		}
		char := string(b[0])
		if char == "\n" {
			break
		}
		offset--
		lastLineSize += n
	}

	lastLine := make([]byte, lastLineSize)
	_, err = file.ReadAt(lastLine, offset+1)

	if err != nil {
		fmt.Println("Not able to read last line with offset", offset, "and the lastLine size", lastLineSize)
		return err
	}

	err = process(file, l)
	if err != nil {
		return err
	}
	return nil
}

func process(f *os.File, l Loader) error {

	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 250*1024)
		return lines
	}}

	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}

	r := bufio.NewReader(f)

	var wg sync.WaitGroup

	for {
		buf := linesPool.Get().([]byte)

		n, err := r.Read(buf)
		buf = buf[:n]

		if n == 0 {
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
				break
			}
			return err
		}

		nextUntilNewline, err := r.ReadBytes('\n')

		if err != io.EOF {
			buf = append(buf, nextUntilNewline...)
		}

		wg.Add(1)
		go func() {
			processChunk(buf, &linesPool, &stringPool, l)
			wg.Done()
		}()

	}

	wg.Wait()
	return nil
}

func processChunk(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool, l Loader) {

	var wg2 sync.WaitGroup

	logs := stringPool.Get().(string)
	logs = string(chunk)

	linesPool.Put(chunk)

	linesSlice := strings.Split(logs, "\n")

	stringPool.Put(logs)

	chunkSize := 300
	n := len(linesSlice)
	noOfThread := n / chunkSize

	if n%chunkSize != 0 {
		noOfThread++
	}

	for i := 0; i < (noOfThread); i++ {

		wg2.Add(1)
		go func(s int, e int, loader Loader) {
			defer wg2.Done() //to avoid deadlocks
			for i := s; i < e; i++ {
				text := linesSlice[i]
				if len(text) == 0 {
					continue
				}
				err := loader.AddItem(text)
				if err != nil {
					fmt.Println("can't parse the line", err)
					continue
				}

			}

		}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(linesSlice)))), l)
	}

	wg2.Wait()
	linesSlice = nil
}
