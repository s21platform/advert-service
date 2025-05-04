package model

import "github.com/s21platform/advert-service/pkg/advert"

type EditAdvert struct {
	ID          int        `db:"id"`
	TextContent string     `db:"text_content"`
	UserFilter  UserFilter `db:"filter"`
}

func (e *EditAdvert) ToDTO(in *advert.EditAdvertIn) {
	e.ID = int(in.Id)
	e.TextContent = in.TextContent
	e.UserFilter = UserFilter{Os: in.UserFilter.Os}
}
