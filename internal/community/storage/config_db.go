package storage

import (
	"another_node/conf"
	"another_node/internal/web_server/pkg"
	"encoding/json"
	"errors"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"golang.org/x/xerrors"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	configDb *gorm.DB
	onlyOnce = sync.Once{}
)

func Init() {
	onlyOnce.Do(func() {
		configDBDsn := conf.GetConfigDbDSN()
		configDBVar, err := gorm.Open(postgres.Open(configDBDsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		configDb = configDBVar
	})
}

type BaseData struct {
	// ID
	ID        int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"softDelete:flag" json:"deleted_at"`
}

type ApiKeyDbModel struct {
	BaseData
	UserId  int64          `gorm:"column:user_id;type:integer" json:"user_id"`
	Disable bool           `gorm:"column:disable;type:bool" json:"disable"`
	ApiKey  string         `gorm:"column:api_key;type:varchar(255)" json:"api_key"`
	KeyName string         `gorm:"column:key_name;type:varchar(255)" json:"key_name"`
	Extra   datatypes.JSON `gorm:"column:extra" json:"extra"`
}

func (*ApiKeyDbModel) TableName() string {
	return "aastar_api_key"
}
func GetApiInfoByApiKey(apiKey string) (*pkg.ApiKeyModel, error) {
	apikeyModel := &ApiKeyDbModel{}
	tx := configDb.Where("api_key = ?", apiKey).First(&apikeyModel)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, tx.Error
		}
		return nil, xerrors.Errorf("error when finding apikey: %w", tx.Error)
	}
	apikeyRes := convertApiKeyDbModelToApiKeyModel(apikeyModel)
	return apikeyRes, nil
}

type ApiModelExtra struct {
	NetWorkLimitEnable bool     `json:"network_limit_enable"`
	DomainWhitelist    []string `json:"domain_whitelist"`
	IPWhiteList        []string `json:"ip_white_list"`
	AirAccountEnable   bool     `json:"airaccount_enable"`
}

func convertApiKeyDbModelToApiKeyModel(apiKeyDbModel *ApiKeyDbModel) *pkg.ApiKeyModel {
	apiKeyModel := &pkg.ApiKeyModel{
		Disable: apiKeyDbModel.Disable,
		ApiKey:  apiKeyDbModel.ApiKey,
		UserId:  apiKeyDbModel.UserId,
	}
	if apiKeyDbModel.Extra != nil {
		// convert To map
		eJson, _ := apiKeyDbModel.Extra.MarshalJSON()
		apiKeyExtra := &ApiModelExtra{}
		err := json.Unmarshal(eJson, apiKeyExtra)
		if err != nil {
			return nil
		}
		apiKeyModel.NetWorkLimitEnable = apiKeyExtra.NetWorkLimitEnable
		if apiKeyExtra.IPWhiteList != nil {
			apiKeyModel.IPWhiteList = mapset.NewSetWithSize[string](len(apiKeyExtra.IPWhiteList))
			for _, v := range apiKeyExtra.IPWhiteList {
				apiKeyModel.IPWhiteList.Add(v)
			}
		}
		if apiKeyExtra.DomainWhitelist != nil {
			apiKeyModel.DomainWhitelist = mapset.NewSetWithSize[string](len(apiKeyExtra.DomainWhitelist))
			for _, v := range apiKeyExtra.DomainWhitelist {
				apiKeyModel.DomainWhitelist.Add(v)
			}
		}
		apiKeyModel.AirAccountEnable = apiKeyExtra.AirAccountEnable

	}
	return apiKeyModel
}
