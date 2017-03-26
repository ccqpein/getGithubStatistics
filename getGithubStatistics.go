package main

import (
	"bufio"
	"context"
	. "fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
	"strconv"
	"time"
)

var userName = "ccqpein"

// Define types
type repoDetail struct {
	Name   string
	Detail []*github.WeeklyStats
}

type repoWeekDetail struct {
	Name       string
	weeklyData [][]int
}

type ChartFile struct {
	ChartType, Title, SubTitle, ValueSuffix, YAxisText string
	XAxisNumbers                                       []int
	Data                                               []repoWeekDetail
}

type intArray1 []int
type intArray2 [][]int

// Authentication and collect repos information
// Codes come from Go official document
func Authentication(userName string) *github.Client {
	fi, err := os.Open("./token")
	check(err)
	ffi := bufio.NewReader(fi)

	str, _, _ := ffi.ReadLine()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(str)},
	)

	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	return client
}

func GetAllRepos(userName string, client *github.Client) []*github.Repository {
	ctx := context.Background()
	ReOption := &github.RepositoryListOptions{Type: "owner"}
	repos, _, err2 := client.Repositories.List(ctx, userName, ReOption)
	if err2 != nil {
		Println(err2)
	}

	//Println(repos)
	return repos
}

func GetWeeklyStats(userName string, repos []*github.Repository, rD chan repoDetail, client *github.Client) {
	ctx := context.Background()
	for _, repo := range repos {
		var A repoDetail
		name := repo.Name
		reposs, _, _ := client.Repositories.ListCodeFrequency(ctx, userName, *name)
		A.Name = *name
		A.Detail = reposs
		rD <- A
	}
}

// Handle the information
func DoWeeklyStats(repoD chan repoDetail, repos []*github.Repository) []repoWeekDetail {
	now := time.Now()
	OneYearAgo := now.AddDate(-1, 0, 0)
	//Println((now.Sub(OneYearAgo).Hours() / 24))
	var repoWeekDetailList []repoWeekDetail

	for i := 0; i < len(repos); i++ {
		var sumAdd, sumDel int
		var weeklyData [][]int

		A := <-repoD
		for _, codeStatues := range A.Detail {
			we := *codeStatues.Week
			//Println(A.Name, we)
			if we == *A.Detail[0].Week && we.After(OneYearAgo) {
				da := int(we.Sub(OneYearAgo).Hours() / (24 * 7))
				for daa := 0; daa < da; daa++ {
					weeklyData = append(weeklyData, []int{0, 0})
				}
			}
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
		//Println(len(weeklyData))
		Println(A.Name, sumAdd, sumDel)
	}
	return repoWeekDetailList
}

//// Make chart file below
//------------------------------------------------------------------------------------
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func MakeChartFile(dataInput *[]repoWeekDetail) ChartFile {
	var chartTemp = ChartFile{
		ChartType:    "column",
		Title:        "LineNumbers",
		SubTitle:     " ",
		ValueSuffix:  "",
		XAxisNumbers: []int{5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70},
		YAxisText:    "Line ",
	}

	for _, i := range *dataInput {
		chartTemp.Data = append(chartTemp.Data, i)

	}
	//Println(chartTemp)
	return chartTemp
}

func (dd intArray1) changeToString() string {
	ss := ""
	for _, num := range dd {
		ss = ss + strconv.Itoa(num) + ", "
	}
	return ss
}

func (dd intArray2) changeToString(index int) string {
	ss := ""
	for _, num := range dd {
		ss = ss + strconv.Itoa(num[index]) + ", "
	}
	return ss
}

func WriteChartFileIn(dataInput ChartFile) error {
	var stringToWrite string

	// Write gochart file
	stringToWrite = Sprintf("ChartType = %s \nTitle = %s \nSubTitle = %s \nValueSuffix = %s \nXAxisNumbers = %s \nYAxisText = %s \n \n# The data and the name of the lines \n",
		dataInput.ChartType,
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
					i.Name, intArray2(i.weeklyData).changeToString(0))
			}
			return stringTemp
		}(dataInput.Data)

	// Save file in folder
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
	client := Authentication(userName)

	allRepos := GetAllRepos(userName, client)
	rD := make(chan repoDetail)
	go GetWeeklyStats(userName, allRepos, rD, client)
	tempFileDat := DoWeeklyStats(rD, allRepos)

	fileData := MakeChartFile(&tempFileDat)
	WriteChartFileIn(fileData)

}
