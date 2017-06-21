# README #

Get the contribution statistics from github.


## Dependency ##

Use golang/github 
[package](https://godoc.org/github.com/google/go-github) 
[Document](https://godoc.org/github.com/google/go-github/github)

Use zieckey/gochart 
[package](https://github.com/zieckey/gochart)

You need to get the authentication from github, [Get Token](https://github.com/google/go-github#authentication) and create token file in code folder.

## Usage ##

Code will read *Token* file and get the repos information from github. Then code will collect weekly contribution data and make the chartfile which could be loaded by gochart in `/tmp`. You can use `gochart` under `tmp` folder.

~~Please remember change the `userName` in code~~
Now you can input username in console.
