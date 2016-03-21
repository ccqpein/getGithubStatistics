package main

import (
	. "fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
)

func main() {
	fi, err := ioutil.ReadFile("./token")
	if err != nil {
		panic(err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(fi)},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	repos, _, err := client.Repositories.List("ccqpein", nil)
	for _, repo := range repos {
		name := string(repo.Name)
		Printf(name)
	}
	//reposs, _, err := client.Repositories.ListCodeFrequency("ccqpein", "Arithmetic-Exercises")

	//Println(reposs)
}
