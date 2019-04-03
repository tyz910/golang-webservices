package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var (
		i         int
		email     string
		isAndroid bool
		isMSIE    bool
	)

	user := User{}
	browsers := map[string]bool{}

	fmt.Fprintln(out, "found users:")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err = user.UnmarshalJSON(scanner.Bytes()); err != nil {
			panic(err)
		}

		isAndroid = false
		isMSIE = false
		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
			} else if strings.Contains(browser, "MSIE") {
				isMSIE = true
			} else {
				continue
			}

			browsers[browser] = true
		}

		if isAndroid && isMSIE {
			email = strings.Replace(user.Email, "@", " [at] ", -1)
			fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", i, user.Name, email))
		}

		i++
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(browsers))
}
