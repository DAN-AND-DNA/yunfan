package yoursql

type Faker_db struct {
	owner_driver *Faker_driver
	dsn          string
}

func new_db(dsn string, owner_driver *Faker_driver) *Faker_db {
	return &Faker_db{dsn: dsn, owner_driver: owner_driver}
}
