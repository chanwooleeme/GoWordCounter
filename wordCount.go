package main

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
	"sync"
	"regexp"
)

var wg sync.WaitGroup

// WordCount는 단어와 빈도수를 저장합니다.
type WordCount struct {
	Word  string
	Count int
}

func inputSplits(filePath string) []string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("파일 읽어오기 실패:", err)
		return nil
	}

	lines := strings.Split(string(content), "\n")
	return lines
}

func extractWords(sentence string) []string {
	sentence = strings.ToLower(sentence)
	re := regexp.MustCompile(`\b[a-zA-Z]+\b`)
	return re.FindAllString(sentence, -1)
}

// mapShuffleReduce 함수는 mapping, shuffling, reducing 과정을 포함합니다.
func mapShuffleReduce(line string, ch chan<- WordCount) {
	defer wg.Done()

	words := extractWords(line)
	db := make(map[string]int)
	for _, word := range words {
		db[word]++
	}
	for word, count := range db {
		ch <- WordCount{Word: word, Count: count}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("사용법: wordcount <파일경로>")
		return
	}
	inputFile := os.Args[1]
	lines := inputSplits(inputFile)
	wordChannel := make(chan WordCount)
	shuffled := make(map[string]int)

	if lines != nil {
		// 고루틴 시작
		for _, line := range lines {
			wg.Add(1)
			go mapShuffleReduce(line, wordChannel)
		}

		// 고루틴이 끝날 때까지 기다린 후 채널을 닫습니다.
		go func() {
			wg.Wait()
			close(wordChannel)
		}()

		// 채널에서 데이터를 읽어 키별로 그룹화합니다.
		for wc := range wordChannel {
			shuffled[wc.Word] += wc.Count
		}

		// 결과 출력
		for word, count := range shuffled {
			fmt.Printf("word: %s, count: %d\n", word, count)
		}
	}
}
