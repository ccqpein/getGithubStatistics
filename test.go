package main

import (
	. "fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	//"net/http"
)

func main() {
	fi, err := ioutil.ReadFile("./token")
	if err != nil {
		panic(err)
	}

	Print(string(fi))

	//resp, _ := http.Get("https://api.github.com/user/ccqpein")
	//Print(resp)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(fi)},
	)

	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		Type:        "public",
		ListOptions: github.ListOptions{PerPage: 10, Page: 2},
	}
	_, resp, err := client.Repositories.ListByOrg("github", opt)

	Print(resp.NextPage)
}
