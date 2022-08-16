package yoursql

import (
	"errors"
	"log"
	"strings"
	"sync"
)

var (
	ERR_EMPTY_RESPS = errors.New("error: need expect resps")
	ERR_BAD_RESPS   = errors.New("error: bad resps")
)

type Faker_resps struct {
	expect_resps *sync.Map // query : *Faker_resp
}

func new_resps() *Faker_resps {
	return &Faker_resps{expect_resps: &sync.Map{}}
}

func (this *Faker_resps) attach(resps []*Faker_resp) error {
	if len(resps) == 0 {
		return ERR_EMPTY_RESPS
	}

	for _, resp := range resps {

		if len(resp.Resp) == 0 {
			return ERR_BAD_RESPS
		}

		if resp == nil {
			return ERR_BAD_RESPS
		}

		if len(resp.Pattern) == 0 {
			return ERR_BAD_RESPS
		}

		ptr_rows := &Faker_rows{
			is_closed:              false,
			current_row_index:      -1,
			expect_error_row_index: -1,
		}

		get_cols := false
		tmp := []*real_row{}
		for _, vals := range resp.Resp {

			tmp_row := &real_row{}
			for col_name, val := range vals {

				if get_cols == false {
					ptr_rows.cols = append(ptr_rows.cols, col_name)
				}

				tmp_row.cols = append(tmp_row.cols, val)
			}

			get_cols = true
			tmp = append(tmp, tmp_row)
		}

		ptr_rows.rows = append(ptr_rows.rows, tmp)
		resp.result_rows = ptr_rows
		this.expect_resps.Store(resp.Pattern, resp)
	}
	return nil
}

func (this *Faker_resps) detach(query string) {
	this.expect_resps.Delete(query)
}

func (this *Faker_resps) find_resp(query string) *Faker_resp {
	if len(query) == 0 {
		err_msg := "cann't match empty query"
		log.Println(err_msg)
		panic(err_msg)
	}

	resp, ok := this.expect_resps.Load(query)
	if ok == false {
		err_msg := "cann't match such query: " + query
		log.Println(err_msg)
		panic(err_msg)
	}

	result_resp, ok := resp.(*Faker_resp)
	if !ok || result_resp == nil {
		err_msg := "bad resp of such query"
		log.Println(err_msg)
		panic(err_msg)

	}

	return result_resp

}

type Faker_resp struct {
	Pattern               string
	Resp                  [](map[string]interface{})
	Rows_affected         int64
	Last_insert_id        int64
	Expect_error          error
	Expect_query_bad_conn func() bool
	Expect_exec_bad_conn  func() bool
	result_rows           *Faker_rows
}

func (this *Faker_resp) is_query_match(query string) bool {
	if strings.Contains(query, this.Pattern) {
		return true
	}

	return false
}
