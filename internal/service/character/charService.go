package character

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
)

type characterService struct {
	db         *repository.DB
	msgToUser  config.MsgToUser
	charConfig config.CharacterConfig
}

func NewCharacterService(db *repository.DB, msgToUser config.MsgToUser, charConfig config.CharacterConfig) *characterService {
	return &characterService{
		db:         db,
		msgToUser:  msgToUser,
		charConfig: charConfig,
	}
}
