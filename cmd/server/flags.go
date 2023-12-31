package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

type flags struct {
	addr            string
	logLevel        string
	fileStoragePath string
	storeInterval   int
	restore         bool
	databaseDsn     string
	downMigrations  bool
	hashKey         string
}

func (f *flags) validate() error {
	if _, _, err := net.SplitHostPort(f.addr); err != nil {
		return fmt.Errorf("address is not valid: %s", err)
	}

	if _, err := zap.ParseAtomicLevel(f.logLevel); err != nil {
		return fmt.Errorf("log level is not valid: %s", err)
	}

	if err := writeFileTest(filepath.Dir(f.fileStoragePath)); err != nil {
		return fmt.Errorf("can`t write file: %s", err)
	}

	return nil
}

func parseFlags() (*flags, error) {
	var f flags

	flag.StringVar(&f.addr, "a", "localhost:8080", "address")
	flag.StringVar(&f.logLevel, "l", "info", "log level")
	flag.StringVar(&f.fileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.StringVar(&f.databaseDsn, "d", "", "database dsn")
	flag.StringVar(&f.hashKey, "k", "", "hash key")
	flag.IntVar(&f.storeInterval, "i", 300, "store interval")
	flag.BoolVar(&f.restore, "r", true, "restore")
	flag.BoolVar(&f.downMigrations, "dm", true, "down migrations after stop")
	flag.Parse()

	envAddr := os.Getenv("ADDRESS")
	if envAddr != "" {
		f.addr = envAddr
	}

	envHashKey := os.Getenv("KEY")
	if envHashKey != "" {
		f.hashKey = envHashKey
	}

	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel != "" {
		f.logLevel = envLogLevel
	}

	envFileStoragePah := os.Getenv("FILE_STORAGE_PATH")
	if envFileStoragePah != "" {
		f.fileStoragePath = envFileStoragePah
	}

	envDatabaseDsn := os.Getenv("DATABASE_DSN")
	if envDatabaseDsn != "" {
		f.databaseDsn = envDatabaseDsn
	}

	envStoreInterval := os.Getenv("STORE_INTERVAL")
	if envStoreInterval != "" {
		v, err := strconv.Atoi(envStoreInterval)

		if err != nil {
			return nil, err
		}

		f.storeInterval = v
	}

	envRestore := os.Getenv("RESTORE")
	if envRestore != "" {
		v, err := strconv.ParseBool(envRestore)

		if err != nil {
			return nil, err
		}

		f.restore = v
	}

	if err := f.validate(); err != nil {
		return nil, err
	}

	return &f, nil
}

func writeFileTest(dirPath string) error {
	// Создаем временный файл
	tmpFilePath := dirPath + "/tmp.txt"
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Удаляем временный файл
	err = os.Remove(tmpFilePath)
	if err != nil {
		return err
	}

	return nil
}
