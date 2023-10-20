package conf

type StorageConf struct {
	fileStoragePath string
}

func NewStorageConf(fileStoragePath string) *StorageConf {
	return &StorageConf{
		fileStoragePath: fileStoragePath,
	}
}

func (c StorageConf) FileStoragePath() string {
	return c.fileStoragePath
}
