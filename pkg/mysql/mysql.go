package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"net/url"
)

func NewMysql(options *Options) (*gorm.DB, error) {
	writeDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		options.Username, options.Password, options.Host, options.Port, options.Database, options.Charset, url.QueryEscape(options.Location))
	writeDB, err := gorm.Open(mysql.Open(writeDsn), &gorm.Config{
		Logger: NewGormCustomLogger(options.LogLevel),
	})
	if err != nil {
		return nil, err
	}
	replicaDsn := writeDsn
	if options.ReplicaHost != "" {
		replicaDsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
			options.Username, options.Password, options.ReplicaHost, options.Port, options.Database, options.Charset, url.QueryEscape(options.Location))
	}
	if err := writeDB.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{mysql.Open(replicaDsn)},
	}).SetConnMaxIdleTime(options.MaxConnectionLifeTime).
		SetConnMaxLifetime(options.MaxConnectionLifeTime).
		SetMaxIdleConns(options.MaxIdleConnections).
		SetMaxOpenConns(options.MaxOpenConnections),
	); err != nil {
		return nil, err
	}
	return writeDB, nil
}
