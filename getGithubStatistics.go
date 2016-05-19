package main

import (
	. "fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"time"
)

var client *github.Client

type repoDetail struct {
	Name   string
	Detail []github.WeeklyStats
}

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
	ReOption := github.RepositoryListOptions{Type: "owner"}

	repos, _, _ := client.Repositories.List("ccqpein", &ReOption)
	return repos
}

func GetWeeklyStats(repos []github.Repository, rD chan repoDetail) {
	for _, repo := range repos {
		var A repoDetail
		name := repo.Name
		reposs, _, _ := client.Repositories.ListCodeFrequency("ccqpein", *name)
		A.Name = *name
		A.Detail = reposs
		//Println(A.Name)
		rD <- A
	}
}

func DoWeeklyStats(repoD chan repoDetail, repos []github.Repository) {
	now := time.Now()
	OneYearAgo := now.AddDate(-1, 0, 0)

	for i := 0; i < len(repos); i++ {
		var sumAdd, sumDel int
		A := <-repoD
		Println(A.Name)
		//Println(now, OneYearAgo)
		for _, codeStatues := range A.Detail {
			we := *codeStatues.Week
			if we.After(OneYearAgo) {
				ad := *codeStatues.Additions
				de := *codeStatues.Deletions
				sumAdd += ad
				sumDel += de
			}
		}
		Println(sumAdd, sumDel)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type ChartFile struct {
	Title, SubTitle, ValueSuffix, YAxisText string
	XAxisNumbers                            []int
	Data                                    map[string][]int
}

func MakeChartFile() {
	os.MkdirAll("~/Desktop/tmp", 0777)
	f, err := os.Create("~/Desktop/tmp/data.chart")
	check(err)
	defer f.Close()

	_, err = f.WriteString("test")
	check(err)
}

func main() {
	allRepos := GetAllRepos()
	rD := make(chan repoDetail)

	go GetWeeklyStats(allRepos, rD)
	DoWeeklyStats(rD, allRepos)

	//make a folder to collect all chart files, for gochart to use
	/*testC := ChartFile{
		Title:        "tt",
		SubTitle:     "ttt",
		ValueSuffix:  "tet",
		YAxisText:    "re",
		XAxisNumbers: []int{1, 2, 3, 4},
		Data:         map[string][]int{"tt": []int{2, 2, 3, 4, 5}},
	}*/
	MakeChartFile()
}
