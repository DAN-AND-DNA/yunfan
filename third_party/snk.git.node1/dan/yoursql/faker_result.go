package yoursql

type Faker_result struct {
	last_insert_id int64
	rows_affected  int64
}

func (this *Faker_result) Set(id, affected int64) {
	this.last_insert_id = id
	this.rows_affected = affected
}

func (this *Faker_result) LastInsertId() (int64, error) {
	return this.last_insert_id, nil
}

func (this *Faker_result) RowsAffected() (int64, error) {
	return this.rows_affected, nil
}
