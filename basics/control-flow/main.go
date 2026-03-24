package main

import (
	"fmt"
	"sandbox/pkg/helpers"
	//"strings"
)

type Student struct {
	name   string
	grades map[string]int // Fach -> Note
}

type Studies struct {
	name   string
	grades []int
}

type Class struct {
	name     string
	students []Studies
}

func main() {
	//fizzBuzz(100)
	//validatePassword("Alexander")

	gradeList := []int{1, 3, 2, 1, 6, 4}
	// Unterscheidung hier nur wegen dem Reset
	currentStreak := []int{}
	bestStreak := []int{}

	GradeStreak(gradeList, currentStreak, bestStreak)

	// Verschachtelungen
	classes := [][]Student{
		{
			{"Alice", map[string]int{"Mathe": 2, "Deutsch": 1, "Englisch": 3}},
			{"Bob", map[string]int{"Mathe": 1, "Deutsch": 3, "Englisch": 2}},
			{"Clara", map[string]int{"Mathe": 3, "Deutsch": 2, "Englisch": 1}},
		},
		{
			{"David", map[string]int{"Mathe": 2, "Deutsch": 2, "Englisch": 2}},
			{"Eva", map[string]int{"Mathe": 1, "Deutsch": 1, "Englisch": 3}},
		},
	}

	classess := [][]Studies{
		{
			{"Alice", []int{4, 3, 3, 2, 1}},
			{"Bob", []int{4, 3, 2, 1, 2}},
			{"Clara", []int{5, 4, 3, 4, 2}},
		},
		{
			{"David", []int{3, 3, 2, 1, 1}},
			{"Eva", []int{2, 4, 3, 2, 1}},
		},
	}

	AverageGrade(classes)
	list := LongestStreak(classess)
	for _, item := range list {
		fmt.Println(item)
	}
}

func fizzBuzz(rng int) {
	for i := 1; i <= rng; i++ {
		result := ""
		if i%3 == 0 {
			result += "Fizz"
		}
		if i%5 == 0 {
			result += "Buzz"
		}
		if i%7 == 0 {
			result += "Bang"
		}

		if result == "" {
			fmt.Println(i)
		} else {
			fmt.Println(result)
		}
	}
}

func validatePassword(password string) {
	checkCounter := 0
	if len(password) >= 8 {
		checkCounter++
	}

	if helpers.HasUppercaseSlow(password) {
		checkCounter++
	}

}

/*
gradeList := []int{4, 3, 2, 5}
	newList := []int{}
	for index, grade := range gradeList {
		fmt.Print(grade)
		for _, betterGrade := range gradeList { // das hier könnte ein problem sein oder?
			fmt.Print(betterGrade)
			if grade < betterGrade || len(newList) == 0 {
				newList = append(newList, grade)
				fmt.Print(newList[index])
			} else {
				fmt.Print("kein Improvement")
				break
			}
		}
	}
*/

func GradeStreak(gradeList []int, currentStreak []int, bestStreak []int) {
	for _, grade := range gradeList {
		if len(currentStreak) == 0 || grade < currentStreak[len(currentStreak)-1] {
			currentStreak = append(currentStreak, grade)
		} else {
			if len(currentStreak) > len(bestStreak) {
				bestStreak = currentStreak
			}
			// reset für nächste iteration
			currentStreak = []int{grade}
		}
	}
	if len(bestStreak) >= 3 {
		fmt.Println("Hurray", bestStreak)
	} else {
		fmt.Println("Beste Streak:", bestStreak)
	}
}

// die Logik hiervon ist perfekt für den generellen flow von Verschachtelungen
func AverageGrade(classes [][]Student) []string {
	returnList := []string{}

	for index, class := range classes {
		bestName := ""
		bestAverage := 100.0
		for _, student := range class {
			sum := 0
			for _, grade := range student.grades {
				sum += grade
			}
			average := float64(sum) / float64(len(student.grades))

			if average < float64(bestAverage) {
				bestAverage = average
				bestName = student.name
			}

		}
		entry := fmt.Sprintf("Klasse %d: Bester Schüler = %s (Schnitt: %.2f)", index+1, bestName, bestAverage)
		returnList = append(returnList, entry)
	}
	return returnList
}

// ähnliche aufgabe mit Verschachelungen
func LongestStreak(classes [][]Studies) []string {
	returnList := []string{}

	for index, class := range classes {
		highestStreak := 0
		bestName := ""
		for _, student := range class {
			// reset for every student
			//currentBestGrade := 20
			streak := 0
			for i, grade := range student.grades {
				if i == 0 {
					continue
				}
				if grade < student.grades[i-1] {
					streak++
					if streak > highestStreak {
						highestStreak = streak
						bestName = student.name
					}
				} else {
					streak = 0
				}
			}
		}
		entry := fmt.Sprintf("Klasse %d und der Student: %s Längste Streak(%d Verbesserungen)", index+1, bestName, highestStreak)
		returnList = append(returnList, entry)
	}
	return returnList
}

// TODO: hier noch etwas zu switch cases, defer etc.
