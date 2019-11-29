package room

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/error2215/go-convert"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"

	model "onesound/compiled/onesound_models"
	"onesound/server/elastic/client"
)

func (r *Request) CheckIfRoomExist(ctx context.Context) (bool, error) {
	count, err := client.GetClient().Count().Index(r.index).Query(r.buildSearchQuery()).Do(ctx)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (r *Request) CreateRoom(ctx context.Context) (*model.Room, error) {
	room := &model.Room{
		Id:   r.getNewId(ctx),
		Name: r.names[0],
	}
	if r.password != "" {
		room.Password = r.password
	}
	jsonData, err := json.Marshal(room)
	if err != nil {
		return nil, err
	}
	resp, err := client.GetClient().Index().Index(r.index).BodyJson(jsonData).Id(convert.String(room.Id)).Do(ctx)
	if err != nil {
		return nil, err
	}
	if convert.Int32(resp.Id) == room.Id {
		return room, nil
	}
	return nil, errors.New("Room was not created due to some problems ")
}

func (r *Request) getNewId(ctx context.Context) int32 {
	hits, err := client.GetClient().Search().
		Index(r.index).Query(elastic.NewBoolQuery()).Sort(idField, false).Size(1).Do(ctx)
	if err != nil {
		logrus.WithField("func", "room.getNewId").Error(err)
	}
	if hits.TotalHits() == 0 {
		return 1
	}
	return convert.Int32(hits.Hits.Hits[0].Id)
}

func (r *Request) DeleteRoom(ctx context.Context) error {

}
