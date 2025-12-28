package main

import(
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

/*
Given two files from the i command line, the first containing the base times of world aquatics, 
the second containing the list of athletes with the races and times achieved. 
The program calculates the FINA score of athletes' race times, adding them to the previous data given
*/


/*
Type struct that contains two float64 (times in seconds from the first file), we will uese that type with a map:
map string - type baseTime
*/
type baseTime struct {
	menTime float64
	womenTime float64
}


//This function converts times from minutes to seconds
/*
To do this we take the time separated by dots as a string, we convert dots to spaces and then we convert the string
to a slice of the single parts (minutes, seconds, centies). In this way we have the control on evrey different case.
*/
func time2seconds(time string) float64 {
	time = strings.ReplaceAll(time, ".", " ")
	parts := strings.Fields(time)
	
	if len(parts) == 3 {
		minutes, _ := strconv.ParseFloat(parts[0], 64)
		seconds, _ := strconv.ParseFloat(parts[1], 64)
		centies, _ := strconv.ParseFloat(parts[2], 64)

		if len(parts[2]) == 2 {
			return (minutes * 60) + seconds + (centies / 100)
		} else if len(parts[2]) == 1 {
			return (minutes * 60) + seconds + (centies / 10)
		}
	}else if len(parts) == 2 {
		seconds, _ := strconv.ParseFloat(parts[0], 64)
		centies, _ := strconv.ParseFloat(parts[1], 64)

		if len(parts[1]) == 2 {
			return seconds + (centies / 100)
		} else if len(parts[1]) == 1 {
			return seconds + (centies / 10)
		}
	} else if len(parts) == 1 {
		seconds, _ := strconv.ParseFloat(parts[0], 64)

		return seconds
	} else {
		return 0 
	}
	return 0 
}


//this function calculate the score points from the given time
func formula(b float64, t float64) int {
	points := 1000 * math.Pow(b/t, 3)
	return int(points)
}

func main() {
	baseTimeMap := make(map[string]baseTime)

//control the parameter of os.Args
	if len(os.Args) < 3 {
		fmt.Println("Errore: manca uno dei tre file")
		return
		}

//open the first file with os.Args
	baseTimePath := os.Args[1]

	baseTimeFile, err := os.Open(baseTimePath)
	if err != nil {
		fmt.Println("Errore: file non trovato")
		return
	}
	defer baseTimeFile.Close()

//open the second file whit os.Args
	athletesTimePath := os.Args[2]

	athletesTimeFile, err := os.Open(athletesTimePath)
	if err != nil {
		fmt.Println("Errore: file non trovato")
		return
	}
	defer athletesTimeFile.Close()

//a cycle to build a map containing all the information given in the first file
	myScanner := bufio.NewScanner(baseTimeFile)
	
	//to skip the first line I call the scan one time
	myScanner.Scan()

	for myScanner.Scan() {
		line := myScanner.Text()
		line = strings.ReplaceAll(line, ",", " ")
		sliceOfLine := strings.Fields(line)
		
		for i := 0; i < len(sliceOfLine); i++ {
			keyName := sliceOfLine[0]
			mt, _ := strconv.ParseFloat(sliceOfLine[1], 64)
			wt, _ := strconv.ParseFloat(sliceOfLine[2], 64)
			baseTimeMap[keyName] = baseTime{ menTime: mt, womenTime: wt }
		}

	}


//this second cycle is to rebuild another .csv file that contains all the previous data plus the score points calculated
	myScanner = bufio.NewScanner(athletesTimeFile)

	//to skip the first line I call the scan one time
	myScanner.Scan()

	for myScanner.Scan() {
		line := myScanner.Text()
		line = strings.ReplaceAll(line, ",", " ")
		sliceOfLine := strings.Fields(line)

		for i := 0; i < len(sliceOfLine); i++ {
			var gender bool
			//representing the first column containing the name
			if i == 0 {
				fmt.Print(sliceOfLine[i])
			//representing also the first column containing the surname
			} else if i == 1 {
				fmt.Print(sliceOfLine[i])
			//representing the third column containing the gender
			} else if i == 2 {
				if sliceOfLine[i] == "m" {
					gender = false
					fmt.Print("m")
				} else if sliceOfLine[i] == "f" {
					gender = true
					fmt.Print("f")
				} else {
					fmt.Println("Errore: genere non specificato")
				}
			//the odd columns containing the stroke of the races
			} else if i % 2 != 0 {
				fmt.Print(sliceOfLine[i])
			//the even columns containing the time of the race of the previous column
			} else {
				stroke := sliceOfLine[i - 1]
				athlethesTime := sliceOfLine[i]
				athlethesTimeSeconds := time2seconds(athlethesTime)
				var baseTime float64
				if gender == false {
					baseTime = baseTimeMap[stroke].menTime
				} else {
					baseTime = baseTimeMap[stroke].womenTime
				}
				fmt.Print(sliceOfLine[i], ",", formula(baseTime ,athlethesTimeSeconds)) 
			}
			
			//This part of the code is used to adjust the output formatting for a .csv file
			if i == 0 {
				fmt.Print(" ")
			} else if i != len(sliceOfLine) - 1 {
				fmt.Print(",")
			}
		}
		fmt.Println()
	}

}
