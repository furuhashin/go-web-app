package main

import (
	"../../../backup"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/matryer/filedb"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type path struct {
	Path string
	Hash string
}

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			log.Fatalln(fatalErr)
		}
	}()
	var (
		interval = flag.Int("interval", 10, "チェック間隔（秒単位）")
		archive  = flag.String("archive", "archive", "アーカイブの保存先")
		dbpath   = flag.String("db", "./db", "filedbデータベースのパス")
	)
	flag.Parse()
	m := &backup.Monitor{
		Destination: *archive,
		Archiver:    backup.ZIP,
		Paths:       make(map[string]string),
	}
	db, err := filedb.Dial(*dbpath)
	if err != nil {
		fatalErr = err
		return
	}
	defer db.Close()
	col, err := db.C("paths")
	if err != nil {
		fatalErr = err
		return
	}
	var path path
	col.ForEach(func(_ int, data []byte) bool {
		if err := json.Unmarshal(data, &path); err != nil {
			fatalErr = err
			return true
		}
		m.Paths[path.Path] = path.Hash
		return false //処理を中止します
	})
	if fatalErr != nil {
		return
	}
	if len(m.Paths) < 1 {
		fatalErr = errors.New("パスがありません。backupツールを作って追加してください")
		return
	}
	check(m, col)
	signalChan := make(chan os.Signal, 1)
	//チャネルが終了シグナルを受け取れるようにする
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//以下無限ループ
Loop:
	for {
		select {
		//指定した時間(コマンドライン引数で設定)が経過するとチャネルを返す
		case <-time.After(time.Duration(*interval) * time.Second):
			check(m, col)
		case <-signalChan:
			//終了
			fmt.Println()
			log.Printf("終了します...")
			break Loop
		}
	}
}

func check(m *backup.Monitor, col *filedb.C) {
	log.Println("チェックします...")
	counter, err := m.Now()
	if err != nil {
		log.Panicln("バックアップに失敗しました:", err)
	}
	//この時点ですべてのパスに対してのチェックが完了し、その総カウントが設定されている
	if counter > 0 {
		log.Printf(" %d個のディレクトリをアーカイブしました\n", counter)
		// ハッシュ値を更新します
		var path path
		col.SelectEach(func(_ int, data []byte) (bool, []byte, bool) {
			//json形式のdataをpath構造体に書き込む
			if err := json.Unmarshal(data, &path); err != nil {
				log.Println("JSONデータの読み込みに失敗しました。"+
					"次の項目に進みます:", err)
				return true, data, false
			}
			//変更の有無にかかわらずハッシュがセットされる
			path.Hash, _ = m.Paths[path.Path]
			//構造体を元にjsonが返される
			newdata, err := json.Marshal(&path)
			if err != nil {
				log.Println("JSONデータの書き出しに失敗しました。"+
					"次の項目に進みます:", err)
				return true, data, false
			}
			//第一引数がtrueなら書き込みを行う。この場合のみ新しい値に書き換わる
			return true, newdata, false
		})
	} else {
		log.Println(" 変更はありません")
	}
}
