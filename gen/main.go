package main

import (
	"CodeQuizAI/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("CodeQuizAI.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 创建生成器实例
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./dao",
		ModelPkgPath: "models",
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	// 关联数据库
	g.UseDB(db)

	// 生成指定模型的 DAO
	g.ApplyBasic(
		models.User{},
		models.Question{},
		models.Paper{},
		models.PaperQuestion{},
		models.TempQuestion{},
	)

	// 执行生成
	g.Execute()
}
