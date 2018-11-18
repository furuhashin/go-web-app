package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Archiver interface {
	DestFmt() func(int64) string
	Archive(src, dest string) error
}

type zipper struct{}

//ZIPには*zipperが格納される。*zipperはArchive()とDestFmt()を実装しているのでArchive型となる
var ZIP Archiver = (*zipper)(nil)

func (z *zipper) Archive(src, dest string) error {
	//archive/testフォルダが作成される
	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return err
	}
	//時刻.zipが作成される
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	w := zip.NewWriter(out)
	defer w.Close()
	//srcはファイルorディレクトリ(./test) backupdの実行ファイルと同じ場所に同名のフォルダを設置する必要がある(44行目)
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		//./test/aaa.textを見つけに行く
		if info.IsDir() {
			return nil //スキップします
		}
		if err != nil {
			return err
		}
		//pathはbackupd配下のフォルダ（適当に書き込むとそれをバックアップしてくれる）
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		//wにdest(archive/test)が設定されているので、そこに時刻.zipを作成していく
		f, err := w.Create(path)
		if err != nil {
			return err
		}
		//f(アーカイブ内のファイル)にin(backupd配下のフォルダのファイルの内容)を書き込んでいる
		io.Copy(f, in)
		return nil
	})
}

func (*zipper) DestFmt() func(int64) string {
	return func(i int64) string {
		return fmt.Sprintf("%d.zip", i)
	}
}
