package storage

import (
	"bufio"
	"log"
	"os"
	"path"
	"path/filepath"
)

type PlainFileDao struct {
	dataDir string
}

func MakePlainFileDao(dataDir string) *PlainFileDao {
	return &PlainFileDao{dataDir}
}

func (t *PlainFileDao) GetFilePath(table string) string {
	return path.Join(t.dataDir, table)
}

// table -- это путь к таблице относительно папки dataDir
func (t *PlainFileDao) Append(table string, data string) {
	filePath := t.GetFilePath(table)

	// Создаем родительскую папку, если она не существует
	err := os.MkdirAll(filepath.Dir(filePath), 0750)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.WriteString(data + "\n"); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func (t *PlainFileDao) GetRows(table string) chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		f, err := os.Open(t.GetFilePath(table))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()
	return ch
}

// prefix в данном случае это путь к папке
func (t *PlainFileDao) GetAllTables() []string {
	entries, err := os.ReadDir(t.dataDir)
	if err != nil {
		log.Fatal(err)
	}
	tables := make([]string, 0, len(entries))
	for _, entry := range entries {
		tables = append(tables, entry.Name())
	}
	return tables
}
