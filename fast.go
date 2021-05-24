package main

import (
	"bufio"
	"fmt"
	"github.com/CyganFx/Profiling-pprof/domain"
	"io"
	"os"
	"strings"
	"sync"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	seenBrowsers := make(map[string]bool, 1000)
	foundUsers := ""

	user := &domain.User{}

	lineIdx := -1
	for {
		lineIdx++
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		err = user.UnmarshalJSON(line)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false
		var wg sync.WaitGroup
		var mutex sync.Mutex

		wg.Add(2)
		{
			go search("Android", user.Browsers, seenBrowsers, &isAndroid, &wg, &mutex)
			go search("MSIE", user.Browsers, seenBrowsers, &isMSIE, &wg, &mutex)
		}
		wg.Wait()

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.Replace(user.Email, "@", " [at] ", -1)
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", lineIdx, user.Name, email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

func search(word string, source []string, seenBrowsers map[string]bool, flag *bool, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()
	for _, browser := range source {
		if strings.Contains(browser, word) {
			*flag = true
			mutex.Lock()
			{
				seenBrowsers[browser] = true
			}
			mutex.Unlock()
		}
	}
}