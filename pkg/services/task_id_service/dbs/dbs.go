package dbs

import (
	"errors"
	"time"
	pkg_dbs "yunfan/pkg/dbs"
	pkg_errcode "yunfan/pkg/errcode"
	"yunfan/pkg/services/task_id_service/error_code"
	sdk_dbs "yunfan/sdk/dbs/task_id_service"
	sdk_errcode "yunfan/sdk/errcode"

	"gorm.io/gorm"
)

func Get_all_segments() (res []sdk_dbs.Segments, err error) {

	db, ok := pkg_dbs.Postgres()
	if !ok {
		log_err := error_code.Err_get_all_segments_disconnect_from_db
		error_code.Print(log_err)
		err = log_err
		return
	}

	result := db.Find(&res)
	if result.Error != nil {
		log_err := pkg_errcode.New("get_all_segement: "+result.Error.Error(), error_code.Me, sdk_errcode.Code_db_internal_error)
		error_code.Print(log_err)
		err = log_err
	}
	return
}

func Create_segments(s *sdk_dbs.Segments) (data *sdk_dbs.Segments, err error) {
	data = new(sdk_dbs.Segments)

	db, ok := pkg_dbs.Postgres()
	if !ok {
		log_err := error_code.Err_get_all_segments_disconnect_from_db
		error_code.Print(log_err)
		err = log_err
		return
	}
	result := db.Where("biz_tag = ?", s.BizTag).First(data)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {

		log_err := pkg_errcode.New("create_segement: has existed "+result.Error.Error(), error_code.Me, sdk_errcode.Code_db_internal_error)
		error_code.Print(log_err)
		err = log_err
		return
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {

		s.CreateTime = time.Now().Unix()
		s.UpdateTime = time.Now().Unix()

		result := db.Create(s)
		if result.Error != nil {

			log_err := pkg_errcode.New("create_segement: Create "+result.Error.Error(), error_code.Me, sdk_errcode.Code_db_internal_error)
			error_code.Print(log_err)
			err = log_err
			return
		}
		data = s
	}
	return
}

func Segments_next_id(tag string) (id *sdk_dbs.Segments, err error) {
	db, ok := pkg_dbs.Postgres()
	if !ok {
		log_err := error_code.Err_get_all_segments_disconnect_from_db
		error_code.Print(log_err)
		err = log_err
		return
	}
	tx := db.Begin()
	id = &sdk_dbs.Segments{}
	if db_err := tx.Exec("update segments set max_id=max_id+step,update_time = ? where biz_tag = ?", time.Now().Unix(), tag).Error; db_err != nil {

		log_err := pkg_errcode.New("Segments_next_id: Update "+db_err.Error(), error_code.Me, sdk_errcode.Code_db_internal_error)
		error_code.Print(log_err)
		err = log_err
		_ = tx.Rollback()

		return
	}
	if db_err := tx.Where("biz_tag = ?", tag).First(id).Error; db_err != nil {

		log_err := pkg_errcode.New("Segments_next_id: Get "+db_err.Error(), error_code.Me, sdk_errcode.Code_db_internal_error)
		error_code.Print(log_err)
		err = log_err
		return
	}
	if db_err := tx.Commit().Error; db_err != nil {

		log_err := pkg_errcode.New("Segments_next_id: Get "+db_err.Error(), error_code.Me, sdk_errcode.Code_db_internal_error)
		error_code.Print(log_err)
		err = log_err

	}

	return
}
