package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

import md "github.com/nao1215/markdown"

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please enter the GitHub author name to check their stars")
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		command := strings.TrimSpace(input)

		switch command {
		case "exit":
			fmt.Println("Bye!")
			return
		default:
			stars, err := getStars(command)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				buildMarkdown(stars)
			}
		}
	}
}

func getStars(author string) (int, error) {
	page := 1
	stars := 0
	for {
		url := "https://api.github.com/users/" + author + "/starred?page=" + strconv.Itoa(page) + "&per_page=100"
		resp, err := http.Get(url)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("GitHub API returned error status code: %d", resp.StatusCode)
		}

		var repos []struct {
			StargazersCount int `json:"stargazers_count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			return 0, err
		}

		if len(repos) == 0 {
			break
		}

		for _, repo := range repos {
			stars += repo.StargazersCount
		}

		page++
	}

	return stars, nil
}

func buildMarkdown(stars int) {
	starsStr := strings.Repeat("*", stars)
	md.NewMarkdown(os.Stdout).
		H1("GitHub Stars:" + strconv.Itoa(stars)).
		PlainText(starsStr).Build()
}
