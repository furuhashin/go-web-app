package backup

import (
	"path/filepath"
	"time"
)

type Monitor struct {
	Paths       map[string]string
	Archiver    Archiver
	Destination string
}

func (m *Monitor) Now() (int, error) {
	var counter int
	//m.Pathsはbackup/cmds/backupd/main.go 57行目で追加される(jsonがそのままマップに変換される)
	for path, lastHash := range m.Paths {
		newHash, err := DirHash(path)
		if err != nil {
			return 0, err
		}
		if newHash != lastHash {
			// アーカイブを行う
			err := m.act(path)
			if err != nil {
				return counter, err
			}
			m.Paths[path] = newHash // ハッシュ値を更新します
			counter++
		}
	}
	return counter, nil
}

func (m *Monitor) act(path string) error {
	dirname := filepath.Base(path)
	filename := m.Archiver.DestFmt()(time.Now().UnixNano())
	//m.Destinationはbackup/cmds/backupd/main.go 31行目で追加される("archive"が入る)
	return m.Archiver.Archive(path, filepath.Join(m.Destination, dirname, filename))
}
