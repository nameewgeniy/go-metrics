package conf

type MemoryStorageConf struct {
	fileStoragePath string
	isRestore       bool
}

func NewMemoryStorageConf(fileStoragePath string, isRestore bool) *MemoryStorageConf {
	return &MemoryStorageConf{
		fileStoragePath: fileStoragePath,
		isRestore:       isRestore,
	}
}

func (c MemoryStorageConf) FileStoragePath() string {
	return c.fileStoragePath
}

func (c MemoryStorageConf) IsRestore() bool {
	return c.isRestore
}
