package data

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/kalougata/go-take-out/internal/model"
	"github.com/kalougata/go-take-out/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Data struct {
	DB *gorm.DB
}

func NewData(conf *config.Config) (*Data, func(), error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci",
		conf.DB.User,
		conf.DB.Passwd,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.DbName,
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Errorf("failed to connect database: %s \n", err)
		return nil, nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)  //设置连接池，空闲
	sqlDB.SetMaxOpenConns(100) //打开
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	log.Info("Succeed to connect database \n")

	if err := db.AutoMigrate(&model.Employee{}); err != nil {
		log.Errorf("failed to migrate database: %s", err)
		return nil, nil, err
	}

	return &Data{
			DB: db,
		}, func() {
			if err := sqlDB.Close(); err != nil {
				log.Warnf("failed to close database: %s \n", err)
			}
		}, nil
}
