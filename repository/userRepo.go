package repository

import (
	"database/sql"
	"errors"
	"simple-product-api/models"
)

type UserRepo struct {
	DB *sql.DB
}

//constructors untuk decoupling
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (ur *UserRepo) GetAllUsers()([]*models.User, error){
	//Alur: Execute query, return domain struct
	var data []*models.User

	rows, err := ur.DB.Query("Select * from user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row = &models.User{}

		if err := rows.Scan(&row.Id, &row.Name, &row.Email, &row.IsAdmin); err != nil{
			return nil, err
		}
		data = append(data, row)
	}
	//udh aman
	return data, nil
}

func (ur *UserRepo) GetUserbyId(id int)(*models.User, error){
	//Alur: Execute query, return domain struct
	var data = &models.User{}

	rows := ur.DB.QueryRow("Select * from user where id = ?", id)

	if err := rows.Scan(&data.Id, &data.Name, &data.Email, &data.IsAdmin); err != nil{
		return nil, err
	}

	return data, nil
}

func (ur *UserRepo) CreateUser(model *models.User)(*models.User, error) {
	//Alur : buat object tampungan untuk simpan request ke domain struct
	

	query := "Insert into user (name, email, isAdmin) values (?,?,?)"
	result, err := ur.DB.Exec(query, model.Name, model.Email, model.IsAdmin)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.New("User Not Created")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	//inject auto generated id ke domain struct
	model.Id = int(id)

	return model, nil
}

func (ur *UserRepo) UpdateUser(id int, model *models.User)(*models.User, error){
	//Alur: Execute query, return domain struct
	//exec query
	query := "update user set name=?, email=? where id = ?"
	result, err := ur.DB.Exec(query, model.Name, model.Email, id)
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

	//ambil isAdmin
	query = "Select isAdmin from user where id = ?"
	if err := ur.DB.QueryRow(query).Scan(&model.IsAdmin); err != nil {
		return nil, err
	}

	return model, nil
}

func (ur *UserRepo) DeleteUser(id int)(*models.User, error){
	//Alur: Execute query, return domain struct
	var data = &models.User{}

	//buat response
	err := ur.DB.QueryRow("select name,email,isAdmin from user where id = ?", id).
		Scan(&data.Name, &data.Email, &data.IsAdmin)
	if err != nil {
		return nil, err
	}

	//exec query
	result, err := ur.DB.Exec("delete from user where id = ?", id)
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