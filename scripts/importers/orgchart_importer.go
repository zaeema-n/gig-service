package main

import (
	"GIG/scripts/crawlers/pdf_crawler/parsers"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

/**
config before running
 */

var category = []string{"Gazette", "OrgChart"}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("file path not specified")
		os.Exit(1)
	}
	filePath := args[0]
	textContent := parsers.ParsePdf(filePath) //parse pdf

	fmt.Println("processing org chart info...")
	splitArray := strings.Split(textContent, "\n")

	var ministryTitle1 = regexp.MustCompile(`^\(\d+\) Minister of`)
	var ministryTitle2 = regexp.MustCompile(`^Minister of`)
	for lineNumber, line := range splitArray {
		ministryMatch1 := ministryTitle1.FindStringSubmatch(line)
		ministryMatch2 := ministryTitle2.FindStringSubmatch(line)
		if len(ministryMatch1) > 0 || len(ministryMatch2) > 0 {
			fmt.Println(line)
			i := lineNumber
			startDepartments := false
			for {
				i++
				if i == len(splitArray) {
					break
				}
				subline := splitArray[i]
				subMinistryMatch1 := ministryTitle1.FindStringSubmatch(subline)
				subMinistryMatch2 := ministryTitle2.FindStringSubmatch(subline)
				if len(subMinistryMatch1) > 0 || len(subMinistryMatch2) > 0 {
					fmt.Println("**************")
					break
				}
				if len(subline) > 2 && (subline[0:2] == "* " || subline[0:2] == " (") { // where department list ends
					startDepartments=false
				}
				if startDepartments {
					fmt.Println("	",subline)
				}
				if subline == "Corporations" && splitArray[i+1][0:1]!="(" { // where department list is assumed to start
					startDepartments = true
				}

			}

		}
	}

	////NER extraction
	//entityTitles, err := utils.ExtractEntityNames(textContent)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//if err := create_entity.CreateEntityFromText(filePath, "Gazette 2017", category, entityTitles); err != nil {
	//	fmt.Println(err.Error(), filePath)
	//}

}
