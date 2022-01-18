package config

import (
	"github.com/Aserose/ArchaicReverie/pkg/logger"
)

type msgConfig struct {
	filename string
	log      logger.Logger
}

type (
	LogMsg struct {
		Format                    string `yaml:"formatLog"`
		Init                      string `yaml:"init"`
		InitOk                    string `yaml:"initOk"`
		InitNoOk                  string `yaml:"initNoOk"`
		Create                    string `yaml:"create"`
		Read                      string `yaml:"read"`
		Update                    string `yaml:"update"`
		Delete                    string `yaml:"delete"`
		CreateToken               string `yaml:"createToken"`
		ReadToken                 string `yaml:"readToken"`
		Marshal                   string `yaml:"marshal"`
		Unmarshal                 string `yaml:"unmarshal"`
		CharWeightOutErr          string `yaml:"charWeightOutErr"`
		CharGrowthOutErr          string `yaml:"charGrowthOutErr"`
		CharGrowthAndWeightOutErr string `yaml:"charGrowthAndWeightOutErr"`
		CharLimitOutErr           string `yaml:"charLimitOutErr"`
		WriterResponse            string `yaml:"writerResponse"`
	}
	MsgToUser struct {
		CharStatus struct {
			CharCreate         string `yaml:"charCreate"`
			CharCreateLimit    string `yaml:"charCreateLimit"`
			CharGet            string `yaml:"charGet"`
			CharUpdate         string `yaml:"charUpdate"`
			CharDelete         string `yaml:"charDelete"`
			CharAllDelete      string `yaml:"charAllDelete"`
			CharNotSelect      string `yaml:"charNotSelect"`
			CharWeightRange    string `yaml:"charWeightRange"`
			CharGrowthRange    string `yaml:"charGrowthRange"`
			CharHeadListFormat string `yaml:"charHeadListFormat"`
			CharListFormat     string `yaml:"charListFormat"`
		} `yaml:"charStatus"`
		AuthStatus struct {
			BusyUsername    string `yaml:"busyUsername"`
			UserNotFound    string `yaml:"userNotFound"`
			InvalidUsername string `yaml:"invalidUsername"`
			InvalidPassword string `yaml:"invalidPassword"`
			SignIn          string `yaml:"signIn"`
			SignUp          string `yaml:"signUp"`
			SignOut         string `yaml:"signOut"`
			SignAlready     string `yaml:"signAlready"`
			PasswordUpdated string `yaml:"passwordUpdated"`
			AccountDeleted  string `yaml:"accountDeleted"`
			NoCharacter     string `yaml:"noCharacter"`
			Error           string `yaml:"error"`
		} `yaml:"authStatus"`
		ActionMsg struct {
			InfoCharFormat         string `yaml:"infoCharFormat"`
			LocationFormat         string `yaml:"locationFormat"`
			JumpOver               string `yaml:"jumpOver"`
			JumpFell               string `yaml:"jumpFell"`
			RemainHealth           string `yaml:"remainHealth"`
			InvalidSum             string `yaml:"invalidSum"`
			NoNeedToRecover        string `yaml:"noNeedToRecover"`
			InvalidFood            string `yaml:"invalidFood"`
			InvalidWeapon          string `yaml:"invalidWeapon"`
			InvalidNumberOfWeapons string `yaml:"invalidNumberOfWeapons"`
			LowHP                  string `yaml:"lowHP"`
		} `yaml:"actionMsg"`
	}
	UtilitiesStr struct {
		CookieName           string `yaml:"cookieName"`
		BadRequest           string `yaml:"badRequest"`
		NumberCharacterLimit int    `yaml:"characterLimit"`
	}
)

func NewMsgConfig(msgFilename string, log logger.Logger) *msgConfig {
	return &msgConfig{
		filename: msgFilename,
		log:      log,
	}
}

func (m msgConfig) InitMsg() (LogMsg, MsgToUser, UtilitiesStr) {
	var (
		logMsg       LogMsg
		msgToUser    MsgToUser
		utilitiesStr UtilitiesStr
	)

	unmarshalYaml(m.filename, m.log,
		&logMsg,
		&msgToUser,
		&utilitiesStr,
	)

	return logMsg, msgToUser, utilitiesStr
}
