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

type (
	LogMsg struct {
		Format                    string `yaml:"formatLog"`
		FormatErr                 string `yaml:"formatErr"`
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
			InfoCharFormat string `yaml:"infoCharFormat"`
			LocationFormat string `yaml:"locationFormat"`
			JumpOver       string `yaml:"jumpOver"`
			JumpFell       string `yaml:"jumpFell"`
		} `yaml:"actionMsg"`
	}
	Endpoints struct {
		AuthEndpoints struct {
			Auth          string `yaml:"auth" json:"Auth"`
			SignIn        string `yaml:"signIn" json:"SignIn"`
			SignUp        string `yaml:"signUp" json:"SignUp"`
			SignOut       string `yaml:"signOut" json:"SignOut"`
			NewPassword   string `yaml:"newPassword" json:"NewPassword"`
			DeleteAccount string `yaml:"deleteAccount" json:"DeleteAccount"`
		} `yaml:"authEndpoints" json:"AuthEndpoints"`
		ApiEndpoints struct {
			Api        string `yaml:"api" json:"Api"`
			CreateChar string `yaml:"createChar" json:"CreateChar"`
			GetAllChar string `yaml:"getAllChar" json:"GetAllChar"`
			SelectChar string `yaml:"selectChar" json:"SelectChar"`
			UpdChar    string `yaml:"updChar" json:"UpdChar"`
			DelChar    string `yaml:"delChar" json:"DelChar"`
			DelAllChar string `yaml:"delAllChar" json:"DelAllChar"`
		} `yaml:"apiEndpoints" json:"ApiEndpoints"`
		ActionEndpoints struct {
			Action                string `yaml:"api" json:"Action"`
			InfoAboutSelectedChar string `yaml:"infoAboutSelectedChar" json:"InfoAboutSelectedChar"`
			BeginActionScene      string `yaml:"beginActionScene" json:"BeginActionScene"`
			ActionScene           string `yaml:"actionScene" json:"ActionScene"`
		} `yaml:"actionEndpoints" json:"ActionEndpoints"`
	}
	UtilitiesStr struct {
		CookieName string `yaml:"cookieName"`
		BadRequest           string `yaml:"badRequest"`
		NumberCharacterLimit int    `yaml:"characterLimit"`
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

func InitStrSet(filename string, log logger.Logger) (LogMsg, MsgToUser, UtilitiesStr, Endpoints) {
	var (
		logMsg       LogMsg
		msgToUser    MsgToUser
		utilitiesStr UtilitiesStr
		endpoints    Endpoints
	)

	unmarshalYaml(filename, log,
		&logMsg,
		&msgToUser,
		&utilitiesStr,
		&endpoints)

	return logMsg, msgToUser, utilitiesStr, endpoints
}

func unmarshalYaml(filename string, log logger.Logger, outs ...interface{}) {
	ymlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("yamlFile.Get err   #%v ", err)
	}

	for _, out := range outs {
		err = yaml.Unmarshal(ymlFile, out)
		if err != nil {
			log.Errorf("error: %v", err)
		}
	}
}

func Init(filename string, log logger.Logger, logMsg LogMsg) (*CfgServer, *CfgServices, *CfgPostgres, error) {

	log.Infof(logMsg.Format, log.PackageAndFileNames(), logMsg.Init)

	var (
		cfgServer   CfgServer
		cfgServices CfgServices
		cfgPostgres CfgPostgres
	)

	re := regexp.MustCompile(`^(.*` + "ArchaicReverie" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	if err := godotenv.Load(string(rootPath) + `/.env`); err != nil {
		log.Errorf(logMsg.FormatErr, log.PackageAndFileNames(), logMsg.Init, err.Error())
	}

	readEnv(log, logMsg,
		&cfgServer,
		&cfgServices,
		&cfgPostgres)

	unmarshalYaml(filename, log, &cfgPostgres)

	log.Infof(logMsg.Format, log.PackageAndFileNames(), logMsg.InitOk)

	return &cfgServer, &cfgServices, &cfgPostgres, nil
}

func readEnv(log logger.Logger, logMsg LogMsg, cfgs ...interface{}) {
	for _, cfg := range cfgs {
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			log.Errorf(logMsg.FormatErr, log.PackageAndFileNames(), logMsg.InitNoOk)
		}
	}
}
