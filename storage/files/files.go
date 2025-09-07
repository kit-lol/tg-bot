package files

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"tg-bot/lib/errorr"
	"tg-bot/storage"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{
		basePath: basePath,
	}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = errorr.WrapIfErr("can't save page", err) }()

	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = errorr.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	// Проверяем существует ли директория
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, storage.ErrNoSavedPages
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Фильтруем только файлы (не директории)
	var validFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() {
			validFiles = append(validFiles, file)
		}
	}

	if len(validFiles) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rng.Intn(len(validFiles))

	return s.decodePage(filepath.Join(path, validFiles[n].Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	fileName, err := fileName(page)
	if err != nil {
		return errorr.Wrap("can't remove page", err)
	}

	path := filepath.Join(s.basePath, page.UserName, fileName)

	msg := fmt.Sprintf("can't remove file: %s", path)

	if err := os.Remove(path); err != nil {
		return errorr.Wrap(fmt.Sprintf(msg, path), err)
	}

	return nil
}

func (s Storage) IsExists(page *storage.Page) (bool, error) {
	// Создаем директорию пользователя если её нет
	userDir := filepath.Join(s.basePath, page.UserName)
	if err := os.MkdirAll(userDir, defaultPerm); err != nil {
		return false, errorr.Wrap("can't create user directory", err)
	}

	fileName, err := fileName(page)
	if err != nil {
		return false, errorr.Wrap("can't get file name", err)
	}

	path := filepath.Join(userDir, fileName)

	// Правильная проверка существования файла
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // Файл не существует
		}
		return false, errorr.Wrap("can't check file existence", err)
	}

	return true, nil // Файл существует
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errorr.Wrap("can't decode page", err)
	}
	defer func() { _ = file.Close() }()

	var page storage.Page

	if err := gob.NewDecoder(file).Decode(&page); err != nil {
		return nil, errorr.Wrap("can't decode page", err)
	}

	return &page, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
