package backup

import (
	"flag"
	"github.com/pkg/errors"
	"log"
)

func main() {
	var fatalError error
	defer func() {
		if fatalErr != nil {
			flag.PrintDefaults()
			log.Fatalln(fatalError)
		}
	}()
	var (
		dppath = flag.String("db", "./backupdata", "データベースのディレクトリへのパス")
	)
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fatalError = errors.New("エラー；コマンドを指定してください")
	}
}
