package auth

import (
	"fmt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestSubQuery(t *testing.T) {
	// 连接数据库
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // 禁用彩色打印
		},
	)
	dsn := "root:rootPassword@tcp(192.168.150.40:23306)/core"
	dsn = fmt.Sprintf("%s?charset=utf8&parseTime=True&loc=Local&time_zone=%s", dsn, url.QueryEscape("'UTC'"))

	mysqlDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger, NamingStrategy: schema.NamingStrategy{
		TablePrefix:   "t_",  // 表前缀
		SingularTable: false, // 复数形式
	}},
	)
	if err != nil {
		t.Error("init mysql db failed, err: ", err.Error())
		return
	}

	var p []entity.Permission
	//err = mysqlDb.Model(entity.ProjectXServices{}).Where(entity.ProjectXServiceFields.ProjectId, 204).Find(&ps).Error
	//if err != nil {
	//	t.Fatal(err)
	//}
	subquery := mysqlDb.Select(entity.ProjectXServiceFields.ServiceId).Where(entity.ProjectXServiceFields.ProjectId, 134).Table(constant.MysqlProjectXServicesTable)

	//subquery2 := mysqlDb.Select(entity.ServiceFields.ID).Where(entity.ServiceFields.ID, subquery).Table(constant.MysqlServicesTable)
	//
	//subquery3 := mysqlDb.Debug().Where(entity.ServiceXPermissoinFields.PermissionId, subquery2).Table(constant.MysqlServiceXPermissoinTable).Find(&p)

	subquery2 := mysqlDb.Select(entity.ServiceXPermissoinFields.PermissionId).Where("service_id in (?)", subquery).Table(constant.MysqlServiceXPermissoinTable)

	subquery3 := mysqlDb.Debug().Where("id in (?)", subquery2).Table(constant.MysqlPermissoinTable).Find(&p)

	if subquery3.Error != nil {
		t.Fatal(subquery3.Error.Error())
	}
	t.Log(p)

}

func TestReg(t *testing.T) {
	p := "/v2/account,/v2/nft,/v2/contract,/v2/orders,/v2/tx"
	path := strings.ReplaceAll(p, ",", "|")
	t.Log(path)
	matched, err := regexp.MatchString(path, "/v2/account/history")
	if err != nil {
		t.Log(err)
	}
	t.Log(matched)
}
