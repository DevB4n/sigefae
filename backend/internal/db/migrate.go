package db

func Migrate() error {
	if err := DB.AutoMigrate(BaseModels()...); err != nil {
		return err
	}
	return DB.AutoMigrate(CircularModels()...)
}
