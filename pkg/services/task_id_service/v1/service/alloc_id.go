package service

import (
	"yunfan/pkg/services/task_id_service/dbs"
	sdk_dbs "yunfan/sdk/dbs/task_id_service"
)

func (s *Service) GetId(tag string) (id int64, err error) {
	s.alloc.Mu.Lock()
	defer s.alloc.Mu.Unlock()
	val, ok := s.alloc.BizTagMap[tag]
	if !ok {
		if err = s.CreateTag(&sdk_dbs.Segments{
			BizTag: tag,
			MaxId:  1,
			Step:   1000,
		}); err != nil {
			return 0, err
		}
		val = s.alloc.BizTagMap[tag]
	}
	return val.GetId(s)
}

func (s *Service) CreateTag(e *sdk_dbs.Segments) error {
	data, err := dbs.Create_segments(e)
	if err != nil {
		return err
	}
	b := &Biz_alloc{
		BazTag:  e.BizTag,
		GetDb:   false,
		IdArray: make([]*IdArray, 0),
	}
	b.IdArray = append(b.IdArray, &IdArray{Start: data.MaxId, End: data.MaxId + data.Step})
	s.alloc.BizTagMap[e.BizTag] = b
	return nil
}
