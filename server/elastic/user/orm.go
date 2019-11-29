package user

import (
	"context"
	"encoding/json"

	"github.com/error2215/go-convert"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"onesound/server/elastic/client"

	model "onesound/compiled/onesound_models"
)

func (r *Request) CreateUser(ctx context.Context) (bool, error) {
	user := &model.User{
		Name:  r.names[0],
		Email: r.emails[0],
	}
	if id, err := r.getNewId(ctx); err != nil {
		user.Id = id
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(r.password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	user.Password = convert.String(pass)

	jsonStr, err := json.Marshal(user)
	if err != nil {
		return false, err
	}
	_, err = client.GetClient().Index().
		Id(convert.String(user.Id)).BodyJson(jsonStr).Index(r.index).Do(ctx)

	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) CheckIfUserExist(ctx context.Context) (bool, error) {
	hits, err := client.GetClient().Search().Index(r.index).
		Size(1).Query(r.buildSearchQuery()).Do(ctx)
	if err != nil {
		return false, err
	}
	if hits.TotalHits() == 0 {
		return false, nil
	}
	return true, nil
}

func (r *Request) FindUsers(ctx context.Context) ([]*model.User, error) {
	hits, err := client.GetClient().Search().Index(r.index).
		Size(convert.Int(r.size)).Query(r.buildSearchQuery()).Do(ctx)
	if err != nil {
		return nil, err
	}
	var res []*model.User
	var oneUser *model.User
	for _, usr := range hits.Hits.Hits {
		oneUser = &model.User{}
		err = json.Unmarshal(usr.Source, oneUser)
		if err != nil {
			return res, err
		}
		res = append(res, &model.User{
			Id:       oneUser.Id,
			Name:     oneUser.Name,
			Email:    oneUser.Email,
			Password: oneUser.Password,
		})
	}
	return res, nil
}

func (r *Request) getNewId(ctx context.Context) (int32, error) {
	hits, err := client.GetClient().Search().
		Index(r.index).Query(elastic.NewBoolQuery()).Sort(idField, false).Size(1).Do(ctx)
	if err != nil {
		logrus.WithField("func", "user.getNewId").Error(err)
		return 0, err
	}
	if hits.TotalHits() == 0 {
		return 1, nil
	}
	return convert.Int32(hits.Hits.Hits[0].Id), nil
}
