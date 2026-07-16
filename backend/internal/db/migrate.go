package db

func Migrate() error {
	return DB.AutoMigrate(AllModels()...)
}