package backup

import (
	"flag"
	"github.com/matryer/filedb"
	"github.com/pkg/errors"
	"log"
	"strings"
)

type path struct {
	Path string
	Hash string
}

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			flag.PrintDefaults()
			log.Fatalln(fatalErr)
		}
	}()
	var (
		dppath = flag.String("db", "./backupdata", "データベースのディレクトリへのパス")
	)
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fatalErr = errors.New("エラー；コマンドを指定してください")
	}
	db, err := filedb.Dial(*dppath)
	if err != nil {
		fatalErr = err
		return
	}
	defer db.Close()
	col, err := db.C("path")
	if err != nil {
		fatalErr = err
		return
	}
	switch strings.ToLower(args[0]) {
	case "list":
	case "add":
	case "remove":
	}
}
