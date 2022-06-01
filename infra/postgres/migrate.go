package postgres

import "hoangphuc.tech/hercules/infra/orm"

func AutoMigrate(model string) error {

	// Migrate
	switch model {
	case "all":
		err := DB().AutoMigrate(&orm.Category{}, &orm.Brand{}, &orm.Item{})
		if err != nil {
			return err
		}
	case "category":
		err := DB().AutoMigrate(&orm.Category{})
		if err != nil {
			return err
		}
	case "brand":
		err := DB().AutoMigrate(&orm.Brand{})
		if err != nil {
			return err
		}
	case "item":
		err := DB().AutoMigrate(&orm.Item{})
		if err != nil {
			return err
		}
	}

	return nil
}
