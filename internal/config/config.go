package config

import (
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"regexp"
)

const (
	ConfigPackageName = "Config"
)

type (
	LogMsg struct {
		Format      string `yaml:"formatLog"`
		FormatErr   string `yaml:"formatErr"`
		Init        string `yaml:"init"`
		InitOk      string `yaml:"initOk"`
		InitNoOk    string `yaml:"initNoOk"`
		Create      string `yaml:"create"`
		Read        string `yaml:"read"`
		Update      string `yaml:"update"`
		Delete      string `yaml:"delete"`
		CreateToken string `yaml:"createToken"`
		ReadToken   string `yaml:"readToken"`
		Marshal     string `yaml:"marshal"`
		Unmarshal   string `yaml:"unmarshal"`
	}
	MsgToUser struct {
		CharStatus struct {
			CharCreate string `yaml:"charCreate"`
			CharGet    string `yaml:"charGet"`
			CharUpdate string `yaml:"charUpdate"`
			CharDelete string `yaml:"charDelete"`
		} `yaml:"charStatus"`
		AuthStatus struct {
			BusyUsername    string `yaml:"busyUsername"`
			UserNotFound    string `yaml:"userNotFound"`
			InvalidUsername string `yaml:"invalidUsername"`
			InvalidPassword string `yaml:"invalidPassword"`
			SignIn          string `yaml:"signIn"`
			SignUp          string `yaml:"signUp"`
			SignOut         string `yaml:"signOut"`
			CookieName      string `yaml:"cookieName"`
		} `yaml:"authStatus"`
	}
)

type (
	CfgServer struct {
		Port string `env:"SERVER_PORT"`
	}
	CfgServices struct {
		HMACSecret   string `env:"SECRET_HMAC"`
		PasswordSalt string `env:"PASSWORD_SALT"`
	}
	CfgPostgres struct {
		DriverName    string `yaml:"driverName"`
		ConnectFormat string `yaml:"postgresConnectFormat"`
		Username      string `env:"POSTGRES_USER"`
		Password      string `env:"POSTGRES_PASSWORD"`
		DBName        string `env:"POSTGRES_DBNAME"`
		SSLMode       string `env:"POSTGRES_SSLMODE"`
	}
)

func InitStrSet(filename string, log logger.Logger) (LogMsg, MsgToUser) {
	var (
		logMsg    LogMsg
		msgToUser MsgToUser
	)

	unmarshalYaml(filename, &logMsg, log)
	unmarshalYaml(filename, &msgToUser, log)

	log.Print(msgToUser.AuthStatus.SignOut)
	return logMsg, msgToUser
}

func unmarshalYaml(filename string, out interface{}, log logger.Logger) {
	ymlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(ymlFile, out)
	if err != nil {
		log.Errorf("error: %v", err)
	}
}

func Init(filename string, log logger.Logger, logMsg LogMsg) (*CfgServer, *CfgServices, *CfgPostgres, error) {

	log.Infof(logMsg.Format, ConfigPackageName, logMsg.Init)

	var (
		cfgServer   CfgServer
		cfgServices CfgServices
		cfgPostgres CfgPostgres
	)

	re := regexp.MustCompile(`^(.*` + "ArchaicReverie" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	if err := godotenv.Load(string(rootPath) + `/.env`); err != nil {
		log.Errorf(logMsg.FormatErr, ConfigPackageName, logMsg.Init, err.Error())
	}

	readEnv(&cfgServer, log, logMsg)
	readEnv(&cfgServices, log, logMsg)
	readEnv(&cfgPostgres, log, logMsg)
	unmarshalYaml(filename, &cfgPostgres, log)

	log.Infof(logMsg.Format, ConfigPackageName, logMsg.InitOk)

	return &cfgServer, &cfgServices, &cfgPostgres, nil
}

func readEnv(cfg interface{}, log logger.Logger, logMsg LogMsg) {
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Errorf(logMsg.FormatErr, ConfigPackageName, logMsg.InitNoOk)
	}
}
