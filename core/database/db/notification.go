package dbcore

import (
	dbcoretablenames "authentication_service/core/database/table_names"
	errm "authentication_service/core/errmodule"
	"authentication_service/core/typescore"
	dbutils "authentication_service/core/utilscore/db"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type NotificationDB struct {
	pool *pgxpool.Pool
}

func NewNotificationDB(pool *pgxpool.Pool) *NotificationDB {
	return &NotificationDB{pool: pool}
}

type NotificationDBI interface {
	GetNotificationListDB(ctx context.Context, options ...typescore.ListDbOptions) ([]*typescore.Notification, uint64, *errm.Error)
	CreateNotificationDB(ctx context.Context, tx pgx.Tx, notificationObj *typescore.Notification, returnObj ...bool) (*typescore.Notification, pgx.Tx, *errm.Error)
	UpdateNotificationDB(ctx context.Context, tx pgx.Tx, paramsUpdate *typescore.Notification, returnObj ...bool) (*typescore.Notification, pgx.Tx, *errm.Error)
	DeleteNotificationDB(ctx context.Context, params *typescore.Notification) *errm.Error
}

// GetNotificationListDB Получение уведомлений
func (u *NotificationDB) GetNotificationListDB(ctx context.Context, options ...typescore.ListDbOptions) ([]*typescore.Notification, uint64, *errm.Error) {
	// logrus.Info("🩵 GetNotificationListDB")
	fields := dbutils.GetStructFieldsDB(&typescore.Notification{}, nil)
	// Добавляем поле total_count
	selectFields := append(fields, "COUNT(*) OVER() AS total_count")

	opts, filter, err := dbutils.GetOptionsDB[typescore.Notification](options...)
	if err != nil {
		return nil, 0, errm.NewError(
			"error_get",
			fmt.Errorf("failed to create SELECT %s SQL: %w", dbcoretablenames.TableNameNotification.ToString(), err),
		)
	}

	query := dbutils.BuildSelectQuery(dbcoretablenames.TableNameNotification.ToString(), selectFields)
	query = dbutils.SetterLimitAndOffsetQuery(query, opts.Offset, opts.Limit)
	query = dbutils.ApplyFilters(query, filter, opts.LikeFields)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, errm.NewError(
			"error_select",
			fmt.Errorf("failed to create SELECT %s SQL: %v", dbcoretablenames.TableNameNotification.ToString(), err),
		)
	}

	rows, err := u.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, errm.NewError(
			"error_select",
			fmt.Errorf("failed to create SELECT %s SQL: %v", dbcoretablenames.TableNameNotification.ToString(), err),
		)
	}
	defer rows.Close()

	var notifications []*typescore.Notification
	var totalCount uint64
	for rows.Next() {
		notification := &typescore.Notification{}
		if err := dbutils.ScanRowsToStructRows(rows, notification, &totalCount); err != nil {
			logrus.Errorf("🔴 error: %s: %+v", "GetNotificationListDB-ScanRowsToStructRows", err)
			continue
		}

		notifications = append(notifications, notification)
	}

	return notifications, totalCount, nil
}

// CreateNotificationDB Создание нового уведомления
func (u *NotificationDB) CreateNotificationDB(ctx context.Context, tx pgx.Tx, notificationObj *typescore.Notification, returnObj ...bool) (*typescore.Notification, pgx.Tx, *errm.Error) {
	// logrus.Info("🩵 CreateNotificationDB")
	if notificationObj == nil {
		return nil, nil, errm.NewError(
			"error_insert",
			fmt.Errorf("failed to create INSERT %s SQL: %v", dbcoretablenames.TableNameNotification.ToString(), errors.New("NotificationObj is nil")),
		)
	}

	err := dbutils.ExecuteTx(ctx, u.pool, tx, func(tx pgx.Tx) error {
		query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert(dbcoretablenames.TableNameNotification.ToString())

		sqlV, args, errW := dbutils.GenerateInsertRequest(query, notificationObj, typescore.InsertOptions{
			IgnoreConflict: true,
		})
		if errW != nil {
			return errW
		}

		_, err := tx.Exec(ctx, *sqlV, args...)
		if err != nil {
			logrus.Errorf("🔴 error: %s: %+v", "CreateNotificationDB-Exec", err)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, tx, err
	}

	returnObjBool := false
	if len(returnObj) > 0 {
		returnObjBool = returnObj[0]
	}

	if returnObjBool {
		options := typescore.ListDbOptions{Filtering: notificationObj}
		getInfoUp, _, errW := u.GetNotificationListDB(ctx, options)
		if errW != nil {
			return nil, tx, errW
		}
		if len(getInfoUp) > 0 {
			return getInfoUp[0], tx, nil
		}
	}
	return nil, tx, nil
}

// UpdateNotificationDB Обновление уведомления
func (u *NotificationDB) UpdateNotificationDB(ctx context.Context, tx pgx.Tx, paramsUpdate *typescore.Notification, returnObj ...bool) (*typescore.Notification, pgx.Tx, *errm.Error) {
	// logrus.Info("🩵 UpdateNotificationDB")
	if paramsUpdate.UniqUUID == nil {
		logrus.Errorf("❌ UpdateNotificationDB error: %s", errors.New("uniq_uuid is nil"))
		return nil, nil, errm.NewError(
			"error_update",
			fmt.Errorf("failed to create UPDATE %s SQL: %v", dbcoretablenames.TableNameNotification.ToString(), errors.New("uniq_uuid is nil")),
		)
	}

	err := dbutils.ExecuteTx(ctx, u.pool, tx, func(tx pgx.Tx) error {
		query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update(dbcoretablenames.TableNameNotification.ToString())

		// Используем функцию для добавления ненулевых полей в запрос
		query = dbutils.AddNonNullFieldsToQueryUpdate(query, *paramsUpdate)

		// Добавляем условие WHERE
		query = query.Where(squirrel.Eq{"uniq_uuid": paramsUpdate.UniqUUID})

		// Генерируем SQL и аргументы
		sql, args, err := query.ToSql()
		if err != nil {
			logrus.Errorf("🔴 error: %s: %+v", "UpdateNotificationDB-ToSql", err)
			return err
		}

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			logrus.Errorf("🔴 error: %s: %+v", "UpdateNotificationDB-Exec", err)
			if strings.Contains(err.Error(), "violates foreign key constraint") {
				return err
			}
			return err
		}

		return nil
	})

	if err != nil {
		return nil, tx, err
	}

	returnObjBool := false
	if len(returnObj) > 0 {
		returnObjBool = returnObj[0]
	}

	if returnObjBool {
		options := typescore.ListDbOptions{Filtering: &typescore.Notification{UniqUUID: paramsUpdate.UniqUUID}}
		getInfoUp, _, errW := u.GetNotificationListDB(ctx, options)
		if errW != nil {
			logrus.Errorf("🔴 error: %s: %+v", "UpdateNotificationDB-GetNotificationListDB", errW)
			return nil, tx, errW
		}
		if len(getInfoUp) > 0 {
			return getInfoUp[0], tx, nil
		}
	}
	return nil, tx, nil
}

func (u *NotificationDB) DeleteNotificationDB(ctx context.Context, params *typescore.Notification) *errm.Error {
	// logrus.Info("🩵 DeleteNotificationDB")
	if params != nil && params.UniqUUID != nil {
		query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
			Delete(dbcoretablenames.TableNameNotification.ToString()).
			Where(squirrel.Eq{"uniq_uuid": *params.UniqUUID}).
			ToSql()

		if err != nil {
			return errm.NewError(
				"error_delete",
				fmt.Errorf("failed to create DELETE %s SQL: %v", dbcoretablenames.TableNameNotification.ToString(), err),
			)
		}

		rows, err := u.pool.Query(ctx, query, args...)
		if err != nil {
			return errm.NewError(
				"error_delete",
				fmt.Errorf("failed to create DELETE %s SQL: %v", dbcoretablenames.TableNameNotification.ToString(), err),
			)
		}
		defer rows.Close()

		return nil
	} else {
		return errm.NewError(
			"error_delete",
			fmt.Errorf("failed to create DELETE %s SQL: %v", dbcoretablenames.TableNameNotification.ToString(), errors.New("uniq_uuid is nil")),
		)
	}
}
