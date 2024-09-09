// GoPlus: https://github.com/goplus/gop
// https://play.goplus.org/p/opjtnHUWQ3d

// Ported from: https://raw.githubusercontent.com/systembugtj/ArtDay/master/main.cpp?token=AARLJP5E75YOS5D67VSRGQLBNSX6Q
// Conceptually, each class runs 3 times per day

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Maps names from ICS Art Day.csv --> ICS Student List.csv
var studentNameConversions = map[string]string{
	//"haji, aahil":               "",	// Doesn't exist at all
	//"darwazeh, ryan":            "",
	//"hautakangas, hugo":         "",
	//"lam, justin":               "",
	//"mahapatra, aditya":         "",
	//"honavalli, bhavana":        "",
	//"apfel, owen":               "",
	//"natarajan, kaushik":        "",
	//"dai, emily":                "",

	"parvathysumesh, shivany":   "Parvathy Sumesh, shivany",
	"uypeckcuat, cara":          "uypeckcuat, cara maxine",
	"mandava, krishna":          "mandava, BALASRI krishna",
	"li, izzy":                  "li, isabella",
	"payton, knauss":            "knauss, payton",
	"witriol, daniel":           "witriol, maya",
	"sabari, vairavan":          "vairavan, sabari",
	"gupta, si":                 "gupta, sia",
	"suhani, kakkar":            "kakkar, suhani",
	"addanki, nitya addanki":    "ADDANKI, NITYA SOMINI",
	"hanna, mohamed":            "mohamed, hanna",
	"lekutai, faye":             "lekutai, PATNISHA",
	"du, bella (jingxi)":        "du, JINGXI",
	"abram-profeta, theo":       "ABRAM-PROFETA, THEOPHILE",
	"yash, bhandari":            "bhandari, yash",
	"botello valadares, alina":  "BOTELLO-VALADARES, ALINA",
	"aarush, kasliwal":          "kasliwal, aarush",
	"liebsack ferrer, eleanor":  "LIEBSACK-FERRER, ELEANOR",
	"syed, rda ayesha":          "SYED, RIDA AYESHA",
	"baddipaidge, samyukth":     "BADDIPADIGE, SAMYUKTH REDDY",
	"massee, oliveir":           "MASSEE, OLIVIER",
	"yousef, mansour":           "MANSOUR, YOUSEF",
	"gowda, hima":               "GOWDA, HIMADYUTHI",
	"abid, yahya":               "ABID, YAHYA MUHAMMAD",
	"mahee, chandrasekhar":      "CHANDRASEKHAR, MAHEE",
	"satishkuma, aadhav":        "SATISHKUMAR, aadhav",
	"harnamhe, ahluwalia":       "AHLUWALIA, HARNAMHE",
	"alexei, chachkov":          "CHACHKOV, ALEXEI",
	"d, emily":                  "du, emily",
	"zheng, xiaorui (rick)":     "ZHENG, XIAORUI",
	"virenn, gupta":             "GUPTA, VIRENN",
	"saenz-badillos, divya":     "SAENZ BADILLOS, divya",
	"karthikeyan, hrish":        "KARTHIKEYAN, HRISHADH",
	"jesse, johnson":            "johnson, jesse",
	"deniz, caglar":             "CAGLAR, DENIZ",
	"nanduri, abhinav":          "NANDURI, ABHINAV VARMA",
	"richburg martinez, amalia": "RICHBURG-MARTINEZ, AMALIA",
}

func toStudentName(first, last string) StudentName {
	sn := strings.Trim(last, " \"") + ", " + strings.Trim(first, " \"") // Remove leading and trailing spaces & quotes
	//sn = strings.Replace(sn, ", ", " ", 1) // replace ", " with " "
	if newName, ok := studentNameConversions[strings.ToLower(sn)]; ok {
		//fmt.Printf("%s --> %s\n", sn, newName)
		sn = newName
	}
	return StudentName(strings.ToUpper(sn))
}

func main() {
	classes := readClasses()
	directory := readDirectory()
	studentPreferences := readStudentPreferences(classes, directory)

	// Manual assignments
	/*	if !assignStudentToClass(studentPreferences, classes, StudentName("WANG-ARYATTAWANICH, AMY"), ClassName("treehouse design"), 0) {
			panic("Couldn't assign student to class")
		}

		if !assignStudentToClass(studentPreferences, classes, StudentName("HANIF, LILA"), ClassName("treehouse design"), 0) {
			panic("Couldn't assign student to class")
		}
		if !assignStudentToClass(studentPreferences, classes, StudentName("SAKIYA, MELODY"), ClassName("treehouse design"), 0) {
			panic("Couldn't assign student to class")
		}
	*/

	// Assign using this grade order:
	for _, g := range []int{11, 10, 9, 11, 10, 9, 8, 7, 6, 11, 10, 9, 8, 7, 6, 8, 7, 6} {
	NextStudent:
		for _, si := range studentPreferences {
			if si.grade != g { // Skip students not in the grade we're currently processing
				continue
			}

			// For the student: Find 1st highest preference doable:
			// Get 1st desired preference; if available, this student gets it; process next student
			// If not available, check next preference; repeat
			for {
				// If no desired classes satisfied, postpone assigning random classes to this student
				// until all other students get their best desired classes.
				if len(si.classesDesired) == 0 { // No classes desired
					continue NextStudent
				}

				desired := si.classesDesired[0]

				// Shuffle the period order to increase "fairness"
				periods := []int{0, 1, 2} // Should use NumPeriods
				//rand.Shuffle(len(periods), func(i, j int) { periods[i], periods[j] = periods[j], periods[i] })
				for _, p := range periods {
					if si.classesAssigned[p] != NoClassName { // The student is already assigned a class this period, try another period
						continue
					}

					if ci := classes[desired][p]; ci.countRemaining > 0 {
						// This class can take the student; assign them to the class
						ci.countRemaining--
						ci.countAssigned++
						si.classesAssigned[p] = desired
						si.classesDesired = si.classesDesired[1:]
						continue NextStudent
					}
				}
				// Highest preference not available; remove it from student's list
				// and try next highest for this same student
				si.classesDesired = si.classesDesired[1:]
			}
			// We never get here
		}
		// Finished all students in this grade; do next grade
	} // Finished all grades

	// Add any students in the directory that didn't supply any desired class
	for sn, di := range directory {
		if di.grade == 12 {
			continue /* Skip 12th graders*/
		}
		if _, ok := studentPreferences[sn]; !ok {
			// This student is not in the list of students, add them
			studentPreferences[sn] = &StudentInfo{grade: di.grade}
		}
	}

	contains := func(s []ClassName, e ClassName) bool {
		for _, a := range s {
			if a == e {
				return true
			}
		}
		return false
	}

	// Randomly assign classes to students with unassigned periods
	for sn, si := range studentPreferences {
		for period := 0; period < NumPeriods; period++ {
			if si.classesAssigned[period] != NoClassName {
				continue // This student has a class assigned this period
			}
			// This student has no class assigned for this period
			findClass := func(period int, exclude []ClassName) (ClassName, *ClassInfo) {
				for cn, ci := range classes {
					// Class this period has openings & student not already assigned to this class
					if ci[period].countRemaining > 0 && !contains(exclude, cn) {
						return cn, ci[period] // We found a class to put this student it
					}
				}
				return NoClassName, nil // We could not find a class to put this student in
			}
			foundClassName, foundClassInfo := findClass(period, si.classesAssigned[:])
			if foundClassName == NoClassName {
				fmt.Printf("Couldn't find a class for %s in period %d\n", sn, period+1)
			} else {
				// put student in the found class
				foundClassInfo.countRemaining--
				foundClassInfo.countAssigned++
				si.classesAssigned[period] = foundClassName
			}
		}
	}
	writeClassSchedule(classes, studentPreferences, directory)
	writeStudentSchedule(studentPreferences, directory)
}

type ClassName string

const NoClassName = ClassName("")

type ClassInfo struct {
	countRemaining int
	countAssigned  int
}

const NumPeriods = 3

type Classes map[ClassName](*[NumPeriods]*ClassInfo)

type StudentName string // must be all Uppercase "last, first"

type DirectoryInfo struct {
	studentName StudentName // Last, First in all caps
	grade       int
	found       bool
}

type Directory map[StudentName]*DirectoryInfo

type StudentInfo struct {
	grade           int // this is copied from DirectoryInfo for easy lookup
	classesDesired  []ClassName
	classesAssigned [NumPeriods]ClassName // All elements default to NoClassName
}
type Students map[StudentName]*StudentInfo

func readClasses() Classes {
	records := readCSV("_Classes.csv", 4) // https://github.com/systembugtj/ArtDay/blob/master/cmake-build-debug/input-class.csv
	c := Classes{}
	for _, r := range records[1:] { // Skip header (1 row)
		className := ClassName(strings.ToLower(r[0]))
		c[className] = &([NumPeriods]*ClassInfo{})
		for p := 0; p < NumPeriods; p++ {
			c[className][p] = &ClassInfo{
				countRemaining: atoi(r[1+p]),
			}
		}
	}
	return c
}

func readDirectory() Directory { // "input-directory.csv"
	records := readCSV("_Directory.csv", 3) // https://github.com/systembugtj/ArtDay/blob/master/cmake-build-debug/input-directory.csv
	d := Directory{}
	for _, r := range records[1:] { // Skip header
		first, last, grade := r[1], r[0], atoi(r[2])
		sn := toStudentName(first, last)
		d[sn] = &DirectoryInfo{
			studentName: sn,
			grade:       grade,
			found:       false,
		}
	}
	return d
}

func readStudentPreferences(classes Classes, directory Directory) Students {
	records := readCSV("_Preferences.csv", 7) // https://github.com/systembugtj/ArtDay/blob/master/cmake-build-debug/input-student.csv
	si := Students{}
	// Find takes a slice and looks for an element in it. If found it returns its index; else -1.
	find := func(slice []ClassName, val ClassName) int {
		for i, item := range slice {
			if item == val {
				return i
			}
		}
		return -1
	}

	for _, r := range records[1:] { // Skip header
		first, last := r[1], r[0]
		sn := toStudentName(first, last)

		classesDesired := []ClassName{}
		for f := 1; f <= 5; f++ {
			cn := ClassName(strings.ToLower(r[1+f])) // Class choices start in column 2
			if _, ok := classes[cn]; !ok {           // Make sure class preference matches a valid class name
				fmt.Printf("Student %s selected invalid class: %s\n", sn, cn)
			} else if find(classesDesired, cn) != -1 { // Class Name repeated
				//fmt.Printf("Student %s selected class 2+ times: %s\n", sn, cn)
			} else { // Class Name is valid & NOT repeated; add it
				classesDesired = append(classesDesired, cn)
			}
		}

		// Add student if in directory, not already added, and grade is not 12
		d, ok := directory[sn]
		if !ok {
			if false { // Skip students not in directory
				fmt.Printf("Student not in directory: %s\n", sn)
				continue
			} else { // Forcably add student to the student directory
				fmt.Printf("Student not in directory: %s\n", sn)
				d = &DirectoryInfo{
					studentName: sn,
					grade:       6, // Assume lowest priority
				}
				directory[sn] = d
			}
		}
		if d.grade != 12 {
			d.found = true         // We found the student
			si[sn] = &StudentInfo{ // If student already added, overwrite them (do not add them multiple times)
				grade:          d.grade,
				classesDesired: classesDesired,
			}
			//fmt.Printf("Student added %s: %#v\n", sn, si[sn])
		}
	}

	// Add students from directory not already found
	/*for name, di := range directory {
		if !di.found  && di.grade != 12 {
			si[name] = &StudentInfo{
				grade:          6, // Assume lowest priority
				classesDesired: []ClassName{}, // No desired classes
			}
		}
	}*/
	return si
}

func assignStudentToClass(students Students, classes Classes, studentName StudentName, className ClassName, period int) bool {
	si := students[studentName]      // We accept the panic if the student doesn't exist
	ci := classes[className][period] // We accept the panic if the class doesn't exist
	if ci.countRemaining > 0 {
		// This class can take the student; assign them to the class
		ci.countRemaining--
		ci.countAssigned++
		si.classesAssigned[period] = className

		// Find index of className in desired classes & remove it
		for i, desiredClass := range si.classesDesired {
			if desiredClass == className { // Found
				si.classesDesired = append(si.classesDesired[:i], si.classesDesired[i+1:]...) // Remove from desired classes
				break
			}
		}
		return true
	}
	return false
}

func writeClassSchedule(classes Classes, students Students, directory Directory) {
	f, err := os.Create("ScheduleByClass.txt")
	PanicOnErr(err)
	defer f.Close()

	for c, ci := range classes {
		for p := 0; p < NumPeriods; p++ {
			fmt.Fprintf(f, "\nClass: %s, Period: %d, Students: %d (remaining=%d)\n", c, p+1, ci[p].countAssigned, ci[p].countRemaining)
			for s, si := range students {
				if si.classesAssigned[p] != c {
					continue // This student is not in this class, try next student
				}
				fmt.Fprintf(f, "   %s\n", directory[s].studentName)
			}
		}
	}
}

func writeStudentSchedule(students Students, directory Directory) {
	var studentSched []string
	for s, si := range students {
		const nca = "No Class Assigned"
		c := [NumPeriods]ClassName{nca, nca, nca}
		for p := 0; p < NumPeriods; p++ {
			if si.classesAssigned[p] != NoClassName { // No class assigned this period
				c[p] = si.classesAssigned[p]
			}
		}
		studentSched = append(studentSched,
			strings.Trim(strings.Join([]string{string(directory[s].studentName), string(c[0]), string(c[1]), string(c[2])}, ","), "\""))
	}
	sort.Strings(studentSched) // Sort students by name

	f, err := os.Create("ScheduleByStudent.csv")
	PanicOnErr(err)
	defer f.Close()
	_, err = f.WriteString("Last,First,Period-1,Period-2,Period-3\n")
	PanicOnErr(err)
	for _, ss := range studentSched {
		_, err = f.WriteString(ss + "\n")
		PanicOnErr(err)
	}
}
