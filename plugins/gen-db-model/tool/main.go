package main

import (
	"flag"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
)

var (
	dbType   = "mysql"
	host     = "127.0.0.1"
	port     = 3306
	username = "root"
	password = "123456"
	schema   = ""
	table    = ""
	outPath  = ""
)

type DBType int

const (
	MySQL      DBType = 1
	PostgreSQL DBType = 2
)

func main() {
	flag.StringVar(&dbType, "db", dbType, "数据库类型（mysql: (default) MySQL; pg: PostgreSQL）")
	flag.StringVar(&host, "h", host, "主机")
	flag.IntVar(&port, "p", port, "端口号")
	flag.StringVar(&username, "u", username, "用户名")
	flag.StringVar(&password, "pw", password, "密码")
	flag.StringVar(&schema, "s", schema, "数据库名")
	flag.StringVar(&table, "t", table, "表名")
	flag.StringVar(&outPath, "o", outPath, "输出路径")
	flag.Parse()

	if dbType == "" || schema == "" || table == "" || outPath == "" {
		panic(fmt.Sprintf("missing params"))
	}

	dsn, databaseType := getDSN()
	switch databaseType {
	case MySQL:
		genMySQLModel(dsn)
		return
	case PostgreSQL:
		genPostgreSQLModel(dsn)
		return
	}
}

func getDSN() (string, DBType) {
	if strings.ToLower(dbType) == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, schema), MySQL
	} else if dbType == "pg" || strings.ToLower(dbType) == "postgresql" {
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable", host, username, password, schema, port), PostgreSQL
	} else {
		panic("unsupported database type")
	}
}

// PostgreSQL

func genPostgreSQLModel(dsn string) {
	conf := gen.Config{
		ModelPkgPath:   outPath,
		FieldNullable:  false,
		FieldCoverable: false,
		FieldSignable:  false,
	}
	g := gen.NewGenerator(conf)
	db, _ := gorm.Open(postgres.Open(dsn))
	g.UseDB(db)
	_ = g.GenerateModel(table, gen.FieldGORMTagReg("^.*$", func(tag field.GormTag) field.GormTag {
		tag.Remove("default")
		tag.Remove("comment")
		return tag
	}))
	g.Execute()
}

// MySQL

func genMySQLModel(dsn string) {
	mdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	tableStatus := &TableStatus{}
	if err = mdb.Raw(fmt.Sprintf("show table status like '%s'", table)).First(&tableStatus).Error; err != nil {
		panic(err)
	}

	ddlms := make([]DDLM, 0)
	if err = mdb.Raw(fmt.Sprintf("show full fields from %s", table)).Find(&ddlms).Error; err != nil {
		panic(err)
	}

	var fields []string
	var primaryKey string
	structStr := fmt.Sprintf(`package dbmodel

// %s %s mysql database.table: %s.%s
type %s struct {
`, FirstUpCase(CamelCase(table)), tableStatus.Comment, schema, table, FirstUpCase(CamelCase(table)))

	for _, v := range ddlms {
		fields = append(fields, fmt.Sprintf("`%s`", v.Field))
		if strings.TrimSpace(v.Comment) != "" {
			structStr += fmt.Sprintf("	// %s\n", strings.ReplaceAll(strings.TrimSpace(v.Comment), "\n", "\n  // "))
		}
		structStr += fmt.Sprintf("	%s %s ", FirstUpCase(CamelCase(v.Field)), getStructType(v))
		if strings.ToUpper(v.Key) == "PRI" {
			primaryKey = v.Field
			structStr += fmt.Sprintf("`gorm:\"column:%s;->\"`\n", v.Field)
		} else if strings.ToUpper(v.Type) == "DATETIME" && strings.ToUpper(v.Default) == "CURRENT_TIMESTAMP" {
			structStr += fmt.Sprintf("`gorm:\"column:%s;->\"`\n", v.Field)
		} else {
			structStr += fmt.Sprintf("`gorm:\"column:%s\"`\n", v.Field)
		}

	}

	structStr += "}\n\n"

	structStr += fmt.Sprintf("// \"select %s from %s where %s=?\"", strings.Join(fields, ","), table, primaryKey)

	structFileName := fmt.Sprintf("./%s.go", table)
	_ = ioutil.WriteFile(structFileName, []byte(structStr), 0666)
	fmt.Println("struct:", structFileName)
}

type (
	DDLM struct {
		Field   string
		Type    string
		Null    string
		Key     string
		Comment string
		Extra   string
		Default string
	}

	TableStatus struct {
		Name    string
		Comment string
	}
)

// TODO 完善更多
func getProtoType(s string) string {
	if strings.Contains(s, "int") {
		return "int64"
	}
	return "string"
}

func getStructType(s DDLM) string {
	types := strings.ToLower(s.Type)
	if strings.Contains(types, "int") {
		return "int64"
	}
	if strings.Contains(strings.ToLower(types), "datetime") {
		return "time.Time"
	}
	if strings.Contains(strings.ToLower(types), "timestamp") {
		return "time.Time"
	}
	if strings.Contains(strings.ToLower(types), "float") {
		return "float32"
	}
	return "string"
}

func CamelCase(s string) string {
	var b []byte
	var wasUnderscore bool
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '_' {
			if wasUnderscore && isASCIILower(c) {
				c -= 'a' - 'A'
			}
			b = append(b, c)
		}
		wasUnderscore = c == '_'
	}
	return string(b)
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

func FirstUpCase(str string) string {
	if len(str) == 0 {
		return str
	}

	if !isASCIILower(str[0]) {
		return str
	}
	c := str[0]
	c -= 'a' - 'A'
	b := []byte(str)
	b[0] = c
	return string(b)
}
