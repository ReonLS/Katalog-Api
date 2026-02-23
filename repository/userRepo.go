package repository

import (
	"database/sql"
	"simple-product-api/models"
)

type UserRepo struct {
	DB *sql.DB
}

func (ur *UserRepo) GetAllUsers()([]models.User, error){
	//Alur: Execute query, return domain struct
	var data []models.User

	rows, err := ur.DB.Query("Select * from user")
	if err != nil {
		return []models.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var row models.User

		if err := rows.Scan(&row.Id, &row.Name, &row.Email, &row.IsAdmin); err != nil{
			return []models.User{}, err
		}
		data = append(data, row)
	}
	//udh aman
	return data, nil
}

func (ur *UserRepo) GetUserbyId(id int)(models.User, error){
	//Alur: Execute query, return domain struct
	var data models.User

	rows := ur.DB.QueryRow("Select * from user where id = ?", id)

	if err := rows.Scan(&data.Id, &data.Name, &data.Email, &data.IsAdmin); err != nil{
		return models.User{}, err
	}

	return data, nil
}

func (ur *UserRepo) CreateUsers(model models.UserRequest)(models.User, error) {
	//Alur : buat object tampungan untuk simpan request ke domain struct
	//autofill isAdmin == 0 (false) karna admin pasti inject akun dari belakang
	//return domain struct

	data := models.User{
		Name: model.Name,
		Email: model.Email,
		IsAdmin: false,
	}

	query := "Insert into user (name, email, isAdmin) values (?,?,?)"
	result, err := ur.DB.Exec(query, data.Name, data.Email, data.IsAdmin)
	if err != nil {
		return models.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.User{}, err
	}
	//inject auto generated id ke domain struct
	data.Id = int(id)

	return data, nil
}

func (ur *UserRepo) UpdateUser(id int, req models.UserRequest)(models.User, error){
	//Alur: Execute query, return domain struct
	data := models.User{
		Id: id,
		Name: req.Name,
		Email: req.Email,
	}
	query := "update product set name=?, email=? where id = ?"
	result, err := ur.DB.Exec(query, data.Name, data.Email, data.Id)

	if err := rows.Scan(&data.Id, &data.Name, &data.Email, &data.IsAdmin); err != nil{
		return models.User{}, err
	}

	return data, nil
}