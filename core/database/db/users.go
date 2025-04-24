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

type UserDB struct {
	pool *pgxpool.Pool
}

func NewUserDB(pool *pgxpool.Pool) *UserDB {
	return &UserDB{pool: pool}
}

type UserDBI interface {
	GetUsersListDB(ctx context.Context, options ...typescore.ListDbOptions) ([]*typescore.User, uint64, *errm.Error)
	CreateUserDB(ctx context.Context, tx pgx.Tx, userObj *typescore.User, returnObj ...bool) (*typescore.User, pgx.Tx, *errm.Error)
	UpdateUserDB(ctx context.Context, tx pgx.Tx, paramsUpdate *typescore.User, returnObj ...bool) (*typescore.User, pgx.Tx, *errm.Error)
	DeleteUserDB(ctx context.Context, params *typescore.User) *errm.Error
}

// GetUsersListDB Получение пользователей
func (u *UserDB) GetUsersListDB(ctx context.Context, options ...typescore.ListDbOptions) ([]*typescore.User, uint64, *errm.Error) {
	// logrus.Info("🩵 GetUsersListDB")
	fields := dbutils.GetStructFieldsDB(&typescore.User{}, nil)
	// Добавляем поле total_count
	selectFields := append(fields, "COUNT(*) OVER() AS total_count")

	opts, filter, err := dbutils.GetOptionsDB[typescore.User](options...)
	if err != nil {
		return nil, 0, errm.NewError(
			"error_get",
			fmt.Errorf("get_options: failed to create SELECT %s SQL: %w", dbcoretablenames.TableNameUsers.ToString(), err),
		)
	}

	query := dbutils.BuildSelectQuery(dbcoretablenames.TableNameUsers.ToString(), selectFields)
	query = dbutils.SetterLimitAndOffsetQuery(query, opts.Offset, opts.Limit)
	query = dbutils.ApplyFilters(query, filter, opts.LikeFields)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, errm.NewError(
			"error_select",
			fmt.Errorf("sql: failed to create SELECT %s SQL: %v", dbcoretablenames.TableNameUsers.ToString(), err),
		)
	}

	rows, err := u.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, errm.NewError(
			"error_select",
			fmt.Errorf("query: failed to create SELECT %s SQL: %v", dbcoretablenames.TableNameUsers.ToString(), err),
		)
	}
	defer rows.Close()

	var users []*typescore.User
	var totalCount uint64
	for rows.Next() {
		user := &typescore.User{}
		if err := dbutils.ScanRowsToStructRows(rows, user, &totalCount); err != nil {
			logrus.Errorf("🔴 error: %s: %+v", "GetUsersListDB-ScanRowsToStructRows", err)
			continue
		}

		users = append(users, user)
	}

	return users, totalCount, nil
}

// CreateUserDB Создает нового пользователя
func (u *UserDB) CreateUserDB(ctx context.Context, tx pgx.Tx, userObj *typescore.User, returnObj ...bool) (*typescore.User, pgx.Tx, *errm.Error) {
	// logrus.Info("🩵 CreateUserDB")
	// Проверяем, что UserProviderW не является nil
	if userObj == nil {
		return nil, nil, errm.NewError(
			"error_insert",
			fmt.Errorf("failed to create INSERT %s SQL: %v", dbcoretablenames.TableNameUsers.ToString(), errors.New("userObj is nil")),
		)
	}

	err := dbutils.ExecuteTx(ctx, u.pool, tx, func(tx pgx.Tx) error {
		query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert(dbcoretablenames.TableNameUsers.ToString())

		// Если LastIPLogin задан (не пустой), используем суффикс для обновления этого поля при конфликте.
		var customSuffix string

		sqlV, args, errW := dbutils.GenerateInsertRequest(query, userObj, typescore.InsertOptions{
			Suffix:         customSuffix,
			IgnoreConflict: false,
		})
		if errW != nil {
			return errW
		}

		_, err := tx.Exec(ctx, *sqlV, args...)
		if err != nil {
			logrus.Errorf("🔴 error: %s: %+v", "CreateUserDB-Exec", err)
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
		users, _, err := u.GetUsersListDB(ctx, typescore.ListDbOptions{Filtering: &typescore.User{
			SystemID: userObj.SystemID,
		}})
		if err != nil {
			return nil, tx, err
		}
		if len(users) > 0 {
			return users[0], tx, nil
		}
	}

	return nil, tx, nil
}

// UpdateUserDB Обновление пользователя
func (u *UserDB) UpdateUserDB(ctx context.Context, tx pgx.Tx, paramsUpdate *typescore.User, returnObj ...bool) (*typescore.User, pgx.Tx, *errm.Error) {
	// logrus.Info("🩵 UpdateUserDB")
	if paramsUpdate.SystemID == nil {
		logrus.Errorf("❌ UpdateUserDB error: %s", errors.New("system_id is nil"))
		return nil, nil, errm.NewError(
			"error_update",
			fmt.Errorf("failed to create UPDATE %s SQL: %v", dbcoretablenames.TableNameUsers.ToString(), errors.New("raw_user_addr is nil")),
		)
	}

	err := dbutils.ExecuteTx(ctx, u.pool, tx, func(tx pgx.Tx) error {

		// Если обновляются другие поля, используем стандартный UPDATE
		query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update(dbcoretablenames.TableNameUsers.ToString())
		var err error

		// Используем функцию для добавления ненулевых полей в запрос
		query = dbutils.AddNonNullFieldsToQueryUpdate(query, *paramsUpdate)

		// Добавляем условие WHERE
		query = query.Where(squirrel.Eq{"system_id": paramsUpdate.SystemID})

		// Генерируем SQL и аргументы
		sql, args, err := query.ToSql()
		if err != nil {
			logrus.Errorf("🔴 error: %s: %+v", "UpdateUserDB-ToSql", err)
			return err
		}

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			logrus.Errorf("🔴 error: %s: %+v", "UpdateUserDB-Exec", err)
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
		options := typescore.ListDbOptions{Filtering: &typescore.User{
			SystemID: paramsUpdate.SystemID,
		}}
		getInfoUp, _, errW := u.GetUsersListDB(ctx, options)
		if errW != nil {
			logrus.Errorf("🔴 error: %s: %+v", "UpdateUserDB-GetUsersListDB", errW)
			return nil, nil, errW
		}
		if len(getInfoUp) > 0 {
			return getInfoUp[0], tx, nil
		}
	}

	return nil, tx, nil
}

func (u *UserDB) DeleteUserDB(ctx context.Context, params *typescore.User) *errm.Error {
	// logrus.Info("🩵 DeleteUserDB")
	if params != nil && params.SystemID != nil {
		query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
			Delete(dbcoretablenames.TableNameUsers.ToString()).
			Where(squirrel.Eq{"system_id": *params.SystemID}).
			ToSql()

		if err != nil {
			return errm.NewError(
				"error_delete",
				fmt.Errorf("failed to create DELETE %s SQL: %v", dbcoretablenames.TableNameUsers.ToString(), err),
			)
		}

		rows, err := u.pool.Query(ctx, query, args...)
		if err != nil {
			return errm.NewError(
				"error_delete",
				fmt.Errorf("failed to create DELETE %s SQL: %v", dbcoretablenames.TableNameUsers.ToString(), err),
			)
		}
		defer rows.Close()

		return nil
	} else {
		return errm.NewError(
			"error_delete",
			fmt.Errorf("failed to create DELETE %s SQL: %v", dbcoretablenames.TableNameUsers.ToString(), errors.New("raw_user_addr is nil")),
		)
	}
}

// GetUsersIDsDB получает список ID пользователей по параметрам фильтрации
func (u *UserDB) GetUsersAddressesDB(ctx context.Context, options ...typescore.ListDbOptions) ([]*string, *errm.Error) {
	// logrus.Info("🩵 GetUsersAddressesDB")
	var selectFields []string
	selectFields = append(selectFields, "raw_user_address")

	opts, filter, err := dbutils.GetOptionsDB[typescore.User](options...)
	if err != nil {
		return nil, errm.NewError(
			"error_get",
			fmt.Errorf("get_options: failed to create SELECT %s SQL: %w", dbcoretablenames.TableNameUsers.ToString(), err),
		)
	}

	query := dbutils.BuildSelectQuery(dbcoretablenames.TableNameUsers.ToString(), selectFields)
	query = dbutils.SetterLimitAndOffsetQuery(query, opts.Offset, opts.Limit)
	query = dbutils.ApplyFilters(query, filter, opts.LikeFields)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errm.NewError(
			"error_select",
			fmt.Errorf("sql: failed to create SELECT %s SQL: %v", dbcoretablenames.TableNameUsers.ToString(), err),
		)
	}

	rows, err := u.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, errm.NewError(
			"error_select",
			fmt.Errorf("query: failed to create SELECT %s SQL: %v", dbcoretablenames.TableNameUsers.ToString(), err),
		)
	}
	defer rows.Close()

	var addresses []*string
	for rows.Next() {
		var userAddr string
		err := rows.Scan(&userAddr)
		if err != nil {
			logrus.Errorf("🔴 error: %s: %+v", "GetUsersAddressesDB-ScanRowsToStructRows", err)
			continue
		}
		// Добавляем user_id в массив
		addresses = append(addresses, &userAddr)
	}

	return addresses, nil
}
