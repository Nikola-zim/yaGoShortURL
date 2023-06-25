package cache

import "yaGoShortURL/internal/entity"

type URLsAllInfo struct {
	IDKey   map[string]entity.JSONAllInfo
	URLKey  map[string]string
	Counter int
}
