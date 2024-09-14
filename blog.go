package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Blog []BlogPost

func GetPostId(postUrl string) (int, error) {
	parsedUrl, err := url.Parse(postUrl) // https://user.livejournal.com/1699420.html
	if err != nil {
		fmt.Printf("Failed to parse URL: %s\n", err)
		return 0, err
	}
	path := parsedUrl.Path                            // /1699420.html
	htmlPath := path[strings.LastIndex(path, "/")+1:] // 1699420.html
	postID := strings.TrimSuffix(htmlPath, ".html")   // 1699420

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		fmt.Printf("Failed to convert post ID to int: %s\n", err)
		return 0, err
	}
	return postIDInt, nil
}
