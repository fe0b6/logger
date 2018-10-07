package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	locker  sync.Mutex
	logHash map[string]*os.File
)

func init() {
	logHash = make(map[string]*os.File)
}

// OpenLog - открываем лог
func OpenLog(name, path string) (err error) {
	locker.Lock()
	defer locker.Unlock()

	// Перенаправляем вывод в файл
	logHash[name], err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// CloseLog - закрываем лог
func CloseLog(name, path string) (err error) {
	locker.Lock()
	defer locker.Unlock()

	wr, ok := logHash[name]
	if !ok {
		err = fmt.Errorf("Логгер %s не найден", name)
		return
	}

	wr.Close()

	delete(logHash, name)

	return
}

// WriteLog - Пишем лог
func WriteLog(name string, data ...interface{}) (err error) {
	locker.Lock()
	defer locker.Unlock()

	wr, ok := logHash[name]
	if !ok {
		err = fmt.Errorf("Логгер %s не найден", name)
		return
	}

	for _, v := range data {
		_, err = wr.WriteString(fmt.Sprintf("%x ", v))
		if err != nil {
			log.Println("[error]", err)
			return
		}
	}

	_, err = wr.WriteString("\n")
	if err != nil {
		log.Println("[error]", err)
		return
	}
	return
}
