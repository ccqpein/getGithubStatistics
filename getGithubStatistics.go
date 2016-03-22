package main

import (
	. "fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
)

type repoDetail struct {
	Name   string
	Detail []github.WeeklyStats
}

var client *github.Client

func GetAllRepos() []github.Repository {
	fi, err := ioutil.ReadFile("./token")
	if err != nil {
		panic(err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(fi)},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client = github.NewClient(tc)

	repos, _, _ := client.Repositories.List("ccqpein", nil)
	return repos
}

func GetWeeklyStats(repos []github.Repository, rD chan repoDetail) {
	for _, repo := range repos {
		name := repo.Name
		reposs, _, _ := client.Repositories.ListCodeFrequency("ccqpein", *name)
		var A repoDetail
		A.Name = *name
		A.Detail = reposs
		//Println(A.Name)
		rD <- A
	}
}

func DoWeeklyStats(repoD chan repoDetail, repos []github.Repository) {
	for i := 0; i < len(repos); i++ {
		A := <-repoD
		Println(A.Name)
		for _, codeStatues := range A.Detail {
			we := *codeStatues.Week
			ad := *codeStatues.Additions
			de := *codeStatues.Deletions
			Println(we, ad, de)
		}
	}
}

func main() {
	allRepos := GetAllRepos()
	rD := make(chan repoDetail)
	go GetWeeklyStats(allRepos, rD)

	DoWeeklyStats(rD, allRepos)
}
