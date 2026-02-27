package repository

import (
	"context"
	"database/sql"
	"errors"
	"simple-product-api/models"
)

type UserRepo struct {
	DB *sql.DB
}

// constructors untuk decoupling
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (ur *UserRepo) Register(ctx context.Context, model *models.User) (*models.User, error) {
	//Alur : buat object tampungan untuk simpan request ke domain struct

	query := "Insert into user (id, name, password, email, role) values (?,?,?,?,?)"
	result, err := ur.DB.ExecContext(ctx, query, model.Id, model.Name, model.Password, model.Email, model.Role)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, err
	}

	return model, nil
}

func (ur *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	//Alur: get user by email, return domain struct
	var model = &models.User{}

	query := "select * from user where email = ?"
	err := ur.DB.QueryRowContext(ctx, query, email).
		Scan(&model.Id, &model.Name, &model.Password, &model.Email, &model.Role)

	if err != nil {
		return nil, err
	}
	return model, nil
}

func (ur *UserRepo) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	//Alur: Execute query, return domain struct
	var data []*models.User

	rows, err := ur.DB.QueryContext(ctx, "Select id, name, email, role from user")
	if err != nil {
		return nil, errors.New("Query Gagal")
	}
	defer rows.Close()

	for rows.Next() {
		var row = &models.User{}
		if err := rows.Scan(&row.Id, &row.Name, &row.Email, &row.Role); err != nil {
			return nil, errors.New("Scannya gagal")
		}
		data = append(data, row)
	}

	//udh aman
	return data, nil
}

func (ur *UserRepo) GetUserById(ctx context.Context, id string) (*models.User, error) {
	//Alur: Execute query, return domain struct
	var data = &models.User{}

	rows := ur.DB.QueryRowContext(ctx, "Select id,name,email,role from user where id = ?", id)

	if err := rows.Scan(&data.Id, &data.Name, &data.Email, &data.Role); err != nil {
		return nil, err
	}

	return data, nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, id string, model *models.User) (*models.User, error) {
	//Alur: Execute query, return domain struct
	//exec query
	query := "update user set name=?, password=? , email=? where id = ?"
	result, err := ur.DB.ExecContext(ctx, query, model.Name, model.Password, model.Email, id)
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.New("User Not Updated")
	}

	//ambil role
	query = "Select id,role from user where id = ?"
	if err := ur.DB.QueryRowContext(ctx, query, id).Scan(&model.Id, &model.Role); err != nil {
		return nil, err
	}
	return model, nil
}

func (ur *UserRepo) DeleteUser(ctx context.Context, id string) (*models.User, error) {
	//Alur: Execute query, return domain struct
	var data = &models.User{}

	//buat response
	err := ur.DB.QueryRowContext(ctx, "select id, name,email,role from user where id = ?", id).
		Scan(&data.Id, &data.Name, &data.Email, &data.Role)
	if err != nil {
		return nil, err
	}

	//exec query
	result, err := ur.DB.ExecContext(ctx, "delete from user where id = ?", id)
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.New("User Not Deleted")
	}
	return data, nil
}
