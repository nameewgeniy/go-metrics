package conf

type MemoryStorageConf struct {
	fileStoragePath string
}

func NewMemoryStorageConf(fileStoragePath string) *MemoryStorageConf {
	return &MemoryStorageConf{
		fileStoragePath: fileStoragePath,
	}
}

func (c MemoryStorageConf) FileStoragePath() string {
	return c.fileStoragePath
}
