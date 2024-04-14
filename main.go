package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	defaultTimeLimit = 30
	defaultFilePath  = "questionAnswers.csv"
)

type questionAnswer struct {
	question string
	answer   string
}

func main() {
	filePath := flag.String("file", defaultFilePath, "path of file which contains questions and answers.")
	timeLimit := flag.Int("timeLimit", defaultTimeLimit, "optional, time limit in seconds for the quiz. Default is 30s")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %s : %s", *filePath, err)
		os.Exit(1)
	}

	csvReader := csv.NewReader(file)
	lines, err := csvReader.ReadAll()
	if err != nil {
		fmt.Printf("Failed to read from file: %s : %s", *filePath, err)
	}

	questionAnswers := getQuestionAnswers(lines)
	rand.New(rand.NewSource(time.Now().Unix()))
	rand.Shuffle(len(questionAnswers), func(i, j int) {
		questionAnswers[i], questionAnswers[j] = questionAnswers[j], questionAnswers[i]
	})

	timer := time.NewTimer(time.Second * time.Duration(*timeLimit))

	timeUp, correctAnswerCount := play(questionAnswers, timer)
	if timeUp {
		fmt.Println("Time up!")
	}
	fmt.Printf("You scored %d out of %d\n", correctAnswerCount, len(questionAnswers))
}

func getQuestionAnswers(lines [][]string) []questionAnswer {
	questions := make([]questionAnswer, 0)
	for _, line := range lines {
		questions = append(questions, questionAnswer{
			question: line[0],
			answer:   strings.ToLower(strings.TrimSpace(line[1])),
		})
	}
	return questions
}

func play(questionAnswers []questionAnswer, timer *time.Timer) (timeUp bool, correct int) {
	correctAnswerCount := 0
	answerCh := make(chan string)
	for quesNum, questionAnswer := range questionAnswers {
		fmt.Printf("%d. %s \n", quesNum+1, questionAnswer.question)
		go getAnswer(answerCh)
		select {
		case <-timer.C:
			return true, correctAnswerCount
		case userAnswer := <-answerCh:
			if userAnswer == questionAnswer.answer {
				correctAnswerCount++
			}
		}
	}
	return false, correctAnswerCount
}

func getAnswer(answerCh chan string) {
	reader := bufio.NewReader(os.Stdin)
	userAnswer, _ := reader.ReadString('\n')
	userAnswer = strings.Trim(userAnswer, " ")
	userAnswer = strings.Trim(userAnswer, "\n")
	userAnswer = strings.ToLower(userAnswer)
	answerCh <- userAnswer
}
