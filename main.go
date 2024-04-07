package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

type questionAnswer struct {
	question string
	answer   string
}

func main() {
	filePath := flag.String("file", "questions.csv", "path of file which contains questions and answers.")
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

	correctAnswerCount := play(questionAnswers)

	fmt.Printf("You scored %d out of %d\n", correctAnswerCount, len(questionAnswers))
}

func getQuestionAnswers(lines [][]string) []questionAnswer {
	questions := make([]questionAnswer, 0)
	for _, line := range lines {
		questions = append(questions, questionAnswer{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		})
	}
	return questions
}

func play(questions []questionAnswer) int {
	correctAnswerCount := 0
	for quesNum, ques := range questions {
		fmt.Printf("%d. %s \n", quesNum+1, ques.question)
		var userAnswer string
		fmt.Scanf("%s", &userAnswer)
		if userAnswer == ques.answer {
			correctAnswerCount++
		}
	}
	return correctAnswerCount
}
