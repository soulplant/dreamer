package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

func OpenTestDb() *gorm.DB {
	os.Remove("test.db")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&Project{}, &User{}, &Task{}, &Stretch{}, &Category{})
	return db
}

func main() {
	db := OpenTestDb()
	test(db)
}

func test(db *gorm.DB) {
	// Read
	var project Project
	if e := db.First(&project, 1000); e.Error != nil {
		fmt.Println("Couldn't find 1000")
	}
	if e := db.First(&project, 1); e.Error != nil {
		fmt.Println("Couldn't find 1")
	}
	db.First(&project, "Name = ?", "Dreamer")

	// Delete - delete project
	// db.Delete(&project)
}
