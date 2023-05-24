package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const MysqlConfig = "root:zhax040214@(localhost:3306)/faas_project_data"

func main() {
	db, err := gorm.Open(mysql.Open(MysqlConfig))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}

	// 生成实例
	g := gen.NewGenerator(gen.Config{
		OutPath:           "../model/query",
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		FieldNullable:     true,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  true,
	})
	g.UseDB(db)

	User := g.GenerateModel("user_user")
	Function := g.GenerateModel("functions")
	Template := g.GenerateModel("templates")
	Triggers := g.GenerateModel("triggers")
	g.ApplyBasic(User)
	g.ApplyBasic(Function)
	g.ApplyBasic(Template)
	g.ApplyBasic(Triggers)
	g.Execute()
}
