package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"microservices/internal/pkg/options"
	"microservices/internal/store"
	"net/url"
)

type mysqlStore struct {
	db *gorm.DB
}

// Users .
func (s *mysqlStore) Users() store.UserStore {
	return newUsers(s)
}

// GetMysqlInstance .
func GetMysqlInstance(opts *options.MySQLOptions) store.DataFactory {
	writeDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		opts.Username, opts.Password, opts.Host, opts.Port, opts.Database, opts.Charset, url.QueryEscape(opts.Location))
	writeDB, err := gorm.Open(mysql.Open(writeDsn), &gorm.Config{
		Logger: logger.Default.LogMode(opts.LogLevel),
	})
	if err != nil {
		panic(err)
	}
	replicaDsn := writeDsn
	if opts.ReplicaHost != "" {
		replicaDsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
			opts.Username, opts.Password, opts.ReplicaHost, opts.Port, opts.Database, opts.Charset, url.QueryEscape(opts.Location))
	}
	writeDB.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{mysql.Open(replicaDsn)},
	}).SetConnMaxIdleTime(opts.MaxConnectionLifeTime).
		SetConnMaxLifetime(opts.MaxConnectionLifeTime).
		SetMaxIdleConns(opts.MaxIdleConnections).
		SetMaxOpenConns(opts.MaxOpenConnections),
	)
	return &mysqlStore{db: writeDB}
}
