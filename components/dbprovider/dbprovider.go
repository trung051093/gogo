package components

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBProvider interface {
	GetDBConnection() *gorm.DB
}

type DBConfig struct {
	Host     string
	Username string
	Password string
	Name     string
	Port     int
	SSLMode  string
	TimeZone string
}

type dbOptions struct {
	Debug        bool
	DstMigration []interface{}
}

type dbprovider struct {
	config *DBConfig
	db     *gorm.DB
}

func NewDBProvider(config *DBConfig, optionsFunc ...func(*dbOptions)) (DBProvider, error) {
	provider := &dbprovider{}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.Host,
		config.Username,
		config.Password,
		config.Name,
		config.Port,
		config.SSLMode,
		config.TimeZone)

	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		return nil, err
	} else {
		provider.config = config
		provider.db = db
	}
	provider.handleOptions(optionsFunc...)
	return provider, nil
}

func initOptions() *dbOptions {
	options := &dbOptions{}
	options.Debug = false
	options.DstMigration = []interface{}{}
	return options
}

func WithDebug(options *dbOptions) {
	options.Debug = true
}

func WithAutoMigration(dst ...interface{}) func(*dbOptions) {
	return func(options *dbOptions) {
		options.DstMigration = dst
	}
}

func (p *dbprovider) handleOptions(optionsFunc ...func(options *dbOptions)) *dbOptions {
	options := initOptions()
	for _, optionFunc := range optionsFunc {
		optionFunc(options)
	}

	if options.Debug {
		p.db.Debug()
	}

	if len(options.DstMigration) > 0 {
		p.db.AutoMigrate(options.DstMigration...)
	}

	return options
}

func (p *dbprovider) GetDBConnection() *gorm.DB {
	return p.db
}
