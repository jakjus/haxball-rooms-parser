package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
        "io/ioutil"
)

var dbCmd = &cobra.Command{
      Use:   "db",
      Short: "Database interaction",
      Args: cobra.MinimumNArgs(1),
}

var dbMysqlCmd = &cobra.Command{
  Use:   "mysql",
  Short: "MySQL database options",
}

var mysqlInitCmd = &cobra.Command{
  Use:   "init",
  Short: "Initialize database",
  Long: "Create database, tables and users for storing and reading data",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Opening MySQL Connection...")
    pool, err := sql.Open("mysql", getDsn()+"?multiStatements=true")
    defer pool.Close()
    if err != nil {
      panic(err)
    }
    fmt.Println("Reading init.sql file...")
    initSql, err := ioutil.ReadFile("./cmd/sql/init.sql")
    if err != nil {
      panic(err)
    }
    _, err = pool.Exec(string(initSql))
    fmt.Println("Running MySQL init query...")
    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Println("Success.")
    }
  },
}

var mysqlUploadCmd = &cobra.Command{
  Use:   "upload",
  Short: "Upload data to database",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Getting data...")
    body, _ := GetData()
    serverList := Parse(body)
    fmt.Println("Opening MySQL Connection...")
    pool, err := sql.Open("mysql", getDsn()+"haxball")
    if err != nil {
      panic(err)
    }
    defer pool.Close()
    fmt.Println("Uploading to MySQL...")
    for _, v := range serverList {
      insertSql := "INSERT INTO server (link, name, flag, private, playersNow, playersMax) VALUES (?, ?, ?, ?, ?, ?)"
      _, err := pool.Exec(insertSql, v.Link, v.Name, v.Flag, v.Private, v.PlayersNow, v.PlayersMax)
      if err != nil {
        fmt.Println(err)
      }
    }
    fmt.Println("Done")
  },
}

var user string
var pass string
var host string
var port string

func getDsn() string {
  return fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, pass, host, port)
}

func init() {
  user = "root"
  pass = "example"
  host = "localhost"
  port = "3306"
  dbCmd.AddCommand(dbMysqlCmd)
  dbMysqlCmd.AddCommand(mysqlInitCmd)
  dbMysqlCmd.AddCommand(mysqlUploadCmd)
  dbMysqlCmd.PersistentFlags().BoolP("help", "", false, "help for this command")
  dbMysqlCmd.PersistentFlags().StringVarP(&user, "user", "u", "root", "Database username")
  dbMysqlCmd.PersistentFlags().StringVarP(&pass, "pass", "p", "example", "Database password")
  dbMysqlCmd.PersistentFlags().StringVarP(&host, "host", "h", "localhost", "Database host (IP address/domain name/...)")
  dbMysqlCmd.PersistentFlags().StringVarP(&port, "port", "P", "3306", "Database port")
}
