### hbparser

[![Go Reference](https://pkg.go.dev/badge/github.com/jakjus/hbparser.svg)](https://pkg.go.dev/github.com/jakjus/hbparser)

CLI tool that shows HaxBall room list. It may also upload room list data to database.

## Installation
Download binary for your OS from *Releases*

## Usage
```
$ ./hbparser

Usage:
  hbparser [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  db          Database interaction
  help        Help about any command
  table       Print out HaxBall rooms as table

Flags:
  -h, --help   help for hbparser

Use "hbparser [command] --help" for more information about a command.
```

## Examples
### Show table
```
./hbparser table
```
![Table output example](<images/example.png>)
### Init and upload to db
1. Run your own DB or use Docker and existing `docker-compose.yml` as follows:
```sh
docker-compose up
```
2. Execute commands
```sh
$ ./hbparser db mysql init
db mysql init
Opening MySQL Connection...
Reading init.sql file...
Running MySQL init query...
Success.

$ ./hbparser db mysql upload
Getting data...
Opening MySQL Connection...
Uploading to MySQL...
Success.
```
3. Login to http://localhost:8080 (Adminer) with user `root` and password `example` to view data.
