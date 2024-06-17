package persistence

import (
	"encoding/json"
	"os"

	"CollabDoc/pkg/document"
)

type Persistence struct {
	filePath string
}

func NewPersistence(filePath string) *Persistence {
	return &Persistence{filePath: filePath}
}

func (p *Persistence) SaveState(ss *document.StateSynchronizer) error {
	data, err := json.Marshal(ss)
	if err != nil {
		return err
	}
	return os.WriteFile(p.filePath, data, 0644)
}

func (p *Persistence) LoadState() (*document.StateSynchronizer, error) {
	data, err := os.ReadFile(p.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return document.NewStateSynchronizer(), nil
		}
		return nil, err
	}
	var ss document.StateSynchronizer
	err = json.Unmarshal(data, &ss)
	if err != nil {
		return nil, err
	}
	return &ss, nil
}
