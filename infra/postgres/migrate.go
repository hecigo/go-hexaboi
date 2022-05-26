package postgres

import "hoangphuc.tech/hercules/infra/orm"

func AutoMigrate() error {

	// Migrate
	err := DB().AutoMigrate(&orm.Category{}, &orm.Brand{}, &orm.Item{})
	if err != nil {
		return err
	}

	return nil
}
