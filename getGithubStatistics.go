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

	//sumVar := make(chan int)
	//delVar := make(chan int)
	type repoDetail struct {
		Name   string
		Detail []github.WeeklyStats
	}
	rD := make(chan repoDetail)

	go func() {
		for _, repo := range repos {
			name := repo.Name
			//Println(*name)
			reposs, _, _ := client.Repositories.ListCodeFrequency("ccqpein", *name)
			var A repoDetail
			A.Name = *name
			A.Detail = reposs
			Println(A.Name)
			rD <- A
		}
	}()

	for i := 0; i < len(repos); i++ {
		A := <-rD
		Println(A.Name)
		for _, codeStatues := range A.Detail {
			we := *codeStatues.Week
			ad := *codeStatues.Additions
			de := *codeStatues.Deletions
			Println(we, ad, de)
		}
	}
}
