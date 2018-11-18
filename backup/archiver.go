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
	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	w := zip.NewWriter(out)
	defer w.Close()
	//srcはファイルorディレクトリ
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil //スキップします
		}
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		f, err := w.Create(path)
		if err != nil {
			return err
		}
		//fにinを書き込んでいる
		io.Copy(f, in)
		return nil
	})
}

func (*zipper) DestFmt() func(int64) string {
	return func(i int64) string {
		return fmt.Sprintf("%d.zip", i)
	}
}
