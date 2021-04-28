package storage

type Option func(storage *Storage) error

func WithConfig(config *TextileConfig) Option {
	return func(s *Storage) error {
		s.config = config
		return nil
	}
}

