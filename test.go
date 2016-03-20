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

	//fi := string(fi)

	//resp, _ := http.Get("https://api.github.com/user/ccqpein")
	//Print(resp)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(fi)},
	)

	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	repos, _, err := client.Repositories.ListContributorsStats("ccqpein", "what_to_eat")

	Println(repos)
	//Println(repp)
}
