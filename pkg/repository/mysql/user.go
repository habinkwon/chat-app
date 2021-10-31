package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/habinkwon/chat-app/graph/model"
	"github.com/habinkwon/chat-app/pkg/util"
)

type User struct {
	DB *sql.DB
}

func (r *User) Get(ctx context.Context, id int64) (user *model.User, err error) {
	var (
		name             string
		nickname         string
		email            sql.NullString
		picture          sql.NullString
		livingPlace      sql.NullString
		preference1      sql.NullString
		preference2      sql.NullString
		preference3      sql.NullString
		selfIntroduction sql.NullString
		role             string
	)
	err = r.DB.QueryRowContext(ctx, `
	SELECT name, nickname, email, picture, livingPlace, preference1, preference2, preference3, selfIntroduction, role
	FROM user
	WHERE id = ?
	`, id).Scan(&name, &nickname, &email, &picture, &livingPlace, &preference1, &preference2, &preference3, &selfIntroduction, &role)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	user = &model.User{
		ID:               id,
		Name:             name,
		Nickname:         nickname,
		Email:            util.NullString(email),
		Picture:          util.NullString(picture),
		LivingPlace:      util.NullString(livingPlace),
		Preference1:      util.NullString(preference1),
		Preference2:      util.NullString(preference2),
		Preference3:      util.NullString(preference3),
		SelfIntroduction: util.NullString(selfIntroduction),
		Role:             model.UserRole(role),
	}
	return
}

func (r *User) GetAll(ctx context.Context, ids []int64) (users []*model.User, err error) {
	in := strings.Repeat(", ?", len(ids)-1)
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := r.DB.QueryContext(ctx, fmt.Sprintf(`
	SELECT id, name, nickname, email, picture, livingPlace, preference1, preference2, preference3, selfIntroduction, role
	FROM user
	WHERE id IN (?%s)
	ORDER BY id
	`, in), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id               int64
			name             string
			nickname         string
			email            sql.NullString
			picture          sql.NullString
			livingPlace      sql.NullString
			preference1      sql.NullString
			preference2      sql.NullString
			preference3      sql.NullString
			selfIntroduction sql.NullString
			role             string
		)
		if err = rows.Scan(&id, &name, &nickname, &email, &picture, &livingPlace, &preference1, &preference2, &preference3, &selfIntroduction, &role); err != nil {
			return
		}
		user := &model.User{
			ID:               id,
			Name:             name,
			Nickname:         nickname,
			Email:            util.NullString(email),
			Picture:          util.NullString(picture),
			LivingPlace:      util.NullString(livingPlace),
			Preference1:      util.NullString(preference1),
			Preference2:      util.NullString(preference2),
			Preference3:      util.NullString(preference3),
			SelfIntroduction: util.NullString(selfIntroduction),
			Role:             model.UserRole(role),
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (r *User) List(ctx context.Context, first int, after int64) (users []*model.User, err error) {
	rows, err := r.DB.QueryContext(ctx, `
	SELECT id, name, nickname, email, picture, livingPlace, preference1, preference2, preference3, selfIntroduction, role
	FROM user
	WHERE id > ?
	ORDER BY id
	LIMIT ?
	`, after, first)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id               int64
			name             string
			nickname         string
			email            sql.NullString
			picture          sql.NullString
			livingPlace      sql.NullString
			preference1      sql.NullString
			preference2      sql.NullString
			preference3      sql.NullString
			selfIntroduction sql.NullString
			role             string
		)
		if err = rows.Scan(&id, &name, &nickname, &email, &picture, &livingPlace, &preference1, &preference2, &preference3, &selfIntroduction, &role); err != nil {
			return
		}
		user := &model.User{
			ID:               id,
			Name:             name,
			Nickname:         nickname,
			Email:            util.NullString(email),
			Picture:          util.NullString(picture),
			LivingPlace:      util.NullString(livingPlace),
			Preference1:      util.NullString(preference1),
			Preference2:      util.NullString(preference2),
			Preference3:      util.NullString(preference3),
			SelfIntroduction: util.NullString(selfIntroduction),
			Role:             model.UserRole(role),
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
