package storage

type Option func(storage *Storage) error

func WithTextile(config *TextileConfig) Option {
	return func(s *Storage) error {
		s.textileConfig = config
		return nil
	}
}

func WithNftStorage(config *NftStorageConfig) Option {
	return func(s *Storage) error {
		s.nftStorageConfig = config
		return nil
	}
}

