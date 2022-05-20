package postgres

import "hoangphuc.tech/hercules/infra/orm"

func AutoMigrate() error {

	// Migrate `item`
	err := DB().AutoMigrate(&orm.Item{})
	if err != nil {
		return err
	}

	return nil
}
