package model

import advert_api "github.com/s21platform/advert-service/pkg/advert"

type EditAdvert struct {
	ID          int        `db:"id"`
	Title       string     `db:"title"`
	TextContent string     `db:"text_content"`
	UserFilter  UserFilter `db:"filter"`
}

func (e *EditAdvert) ToDTO(in *advert_api.EditAdvertIn) {
	e.ID = int(in.Id)
	e.Title = in.Title
	e.TextContent = in.TextContent
	e.UserFilter = UserFilter{Os: in.UserFilter.Os}
}
