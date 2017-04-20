package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"solarwinds.com/goswisclient/swis"
)

func main() {

	var host = flag.String("host", "localhost", "The address of the swis service")
	var user = flag.String("user", "", "The orion username")
	var password = flag.String("password", "", "The password of the account")
	var port = flag.Int("port", 17778, "Port that SWIS service is listening to")

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	client := swis.NewClient(*host, *user, *password, *port)

	var text string
	for text != "quit" { // break the loop if text == "q"
		scanner.Scan()
		text = scanner.Text()
		if text != "quit" {
			if strings.HasPrefix(text, "query>") {
				query := strings.TrimPrefix(text, "query>")
				result, err := client.Query(query)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error running query: %v", err)
					continue
				}

				display(result)
			}
		}
	}
}

func display(results []map[string]interface{}) {
	for row, value := range results {
		fmt.Println(row+1, ":")
		print(value)
	}
}

func print(m map[string]interface{}) {
	for key, value := range m {
		fmt.Println("\t", key, ":", value)
	}
}

func convert(m []interface{}) [][]string {

	result := make([][]string, len(m)+1)

	value := m[0].(map[string]interface{})
	// caption
	result[0] = make([]string, len(value))

	i := 0
	for k := range value {
		result[0][i] = k
		i++
	}

	for row, value := range m {
		rowIndex := row + 1
		rowData := value.(map[string]interface{})
		result[rowIndex] = make([]string, len(rowData))

		col := 0
		for _, value := range rowData {
			result[rowIndex][col] = fmt.Sprintf("%v", value)
			col++
		}
	}

	return result
}
