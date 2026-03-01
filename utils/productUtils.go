package utils

import (
	"errors"
	"strings"
)

type Category string

const (
	CategoryClothes   Category = "Baju"
	CategoryAccessory Category = "Aksesoris"
	CategoryFootwear  Category = "Alas Kaki"
	CategoryInner     Category = "Dalaman"
)

//validasi namaprod, kategori, price, stock

func ValidateProduct(namaprod string, kategori string, price float64, stock int) []error {
	var listErr []error
	validCat := []Category{CategoryClothes, CategoryAccessory, CategoryFootwear, CategoryInner}

	//validasi nama
	if namaprod == "" || len(strings.TrimSpace(namaprod)) == 0 {
		listErr = append(listErr, errors.New("Name must not be empty!"))
	}

	//validasi kategori
	if kategori == "" || len(strings.TrimSpace(kategori)) == 0 {
		listErr = append(listErr, errors.New("Category must not be empty!"))
	} else {
		//ngecek kategori dengan setiap valid Category, valid jd true kalo kategori == string(each)
		valid := false

		for _, each := range validCat {
			//ketemu yg cocok, lgsg jd true
			if kategori == string(each) {
				valid = true
				break
			}
		}
		if !valid {
			listErr = append(listErr, errors.New("Category not in list"))
		}
	}

	//validasi price
	if price <= 0 {
		listErr = append(listErr, errors.New("Price may not be negative nor empty"))
	}

	//validasi stock
	if stock < 0 {
		listErr = append(listErr, errors.New("Stock may not be negative"))
	}
	//berarti aman
	return listErr
}
