package main

import (
	. "fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var client *github.Client
var userName string

type repoDetail struct {
	Name   string
	Detail []github.WeeklyStats
}

func GetAllRepos(userName string) []github.Repository {
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

	repos, _, _ := client.Repositories.List(userName, &ReOption)
	return repos
}

func GetWeeklyStats(userName string, repos []github.Repository, rD chan repoDetail) {
	for _, repo := range repos {
		var A repoDetail
		name := repo.Name
		reposs, _, _ := client.Repositories.ListCodeFrequency(userName, *name)
		A.Name = *name
		A.Detail = reposs
		rD <- A
	}
}

type repoWeekDetail struct {
	Name       string
	weeklyData [][]int
}

func DoWeeklyStats(repoD chan repoDetail, repos []github.Repository) []repoWeekDetail {
	now := time.Now()
	OneYearAgo := now.AddDate(-1, 0, 0)
	var repoWeekDetailList []repoWeekDetail

	for i := 0; i < len(repos); i++ {
		var sumAdd, sumDel int
		var weeklyData [][]int

		A := <-repoD
		for _, codeStatues := range A.Detail {
			we := *codeStatues.Week
			if we.After(OneYearAgo) {
				ad := *codeStatues.Additions
				de := *codeStatues.Deletions
				var temp = []int{ad, de}
				weeklyData = append(weeklyData, temp)
				sumAdd += ad
				sumDel += de
			}
		}

		var tempDetail = repoWeekDetail{Name: A.Name, weeklyData: weeklyData}
		repoWeekDetailList = append(repoWeekDetailList, tempDetail)
		//Println(weeklyData, sumAdd, sumDel)
	}
	return repoWeekDetailList
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type ChartFile struct {
	Title, SubTitle, ValueSuffix, YAxisText string
	XAxisNumbers                            []int
	Data                                    []repoWeekDetail
}

func MakeChartFile(dataInput *[]repoWeekDetail) ChartFile {
	var chartTemp = ChartFile{
		Title:        "LineNumbers",
		SubTitle:     " ",
		ValueSuffix:  "",
		XAxisNumbers: []int{5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70},
		YAxisText:    "Line ",
	}

	for _, i := range *dataInput {
		chartTemp.Data = append(chartTemp.Data, i)

	}
	Println(chartTemp)
	return chartTemp
}

type intArray1 []int
type intArray2 [][]int

func (dd intArray1) changeToString() string {
	ss := ""
	for _, num := range dd {
		ss = ss + strconv.Itoa(num) + ", "
	}
	return ss
}

func (dd intArray2) changeToString2(index int) string {
	ss := ""
	for _, num := range dd {
		ss = ss + strconv.Itoa(num[index]) + ", "
	}
	return ss
}

func WriteChartFileIn(dataInput ChartFile) error {
	var stringToWrite string

	stringToWrite = Sprintf("ChartType = bar \nTitle = %s \nSubTitle = %s \nValueSuffix = %s \nXAxisNumbers = %s \nYAxisText = %s \n \n# The data and the name of the lines \n",
		dataInput.Title,
		dataInput.SubTitle,
		dataInput.ValueSuffix,
		intArray1(dataInput.XAxisNumbers).changeToString(),
		dataInput.YAxisText)

	stringToWrite = stringToWrite +
		func(d []repoWeekDetail) string {
			stringTemp := ""
			for _, i := range d {
				stringTemp = stringTemp + Sprintf("Data|%s = %s \n",
					i.Name, intArray2(i.weeklyData).changeToString2(0))
			}
			return stringTemp
		}(dataInput.Data)

	if _, err := os.Stat("./tmp"); err != nil {
		if os.IsNotExist(err) {
			Print("Create new folder store data")
			os.MkdirAll("./tmp", 0777)
		}
	}

	f, err := os.Create("./tmp/data.chart")
	check(err)
	defer f.Close()

	_, err = f.WriteString(stringToWrite)
	check(err)
	return err
}

func main() {

	//Scanf("Input your name %s \n", &userName)
	//Need make username can be changed from cli
	userName = "ccqpein"
	allRepos := GetAllRepos(userName)
	rD := make(chan repoDetail)

	go GetWeeklyStats(userName, allRepos, rD)
	tempFileDat := DoWeeklyStats(rD, allRepos)
	fileData := MakeChartFile(&tempFileDat)
	WriteChartFileIn(fileData)
}
