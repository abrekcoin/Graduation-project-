package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name    *string            `json:"first_name" validate:"required,min=2,max=30"`
	Last_Name     *string            `json:"last_name"  validate:"required,min=2,max=30"`
	Password      *string            `json:"password"   validate:"required,min=6"`
	Email         *string            `json:"email"      validate:"email,required"`
	Phone         *string            `json:"phone"      validate:"required"`
	Created_At    time.Time          `json:"created_at"`
	Updated_At    time.Time          `json:"updtaed_at"`
	User_ID       string             `json:"user_id"`
	UserCart      []ProductUser      `json:"usercart" bson:"usercart"`
	Order_Status  []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *uint64            `json:"price"`
	Tax_Ratio    int                `json:"tax_ratio"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name" bson:"product_name"`
	Price        int                `json:"price"  bson:"price"`
	Tax_Ratio    int                `json:"tax_ratio" bson:"tax_ratio"`
	Tax_Value    int                `json:"tax_value" bson:"tax_value"`
}

type Order struct {
	Order_ID     primitive.ObjectID `bson:"_id"`
	Order_Cart   []ProductUser      `json:"order_list"  bson:"order_list"`
	Orderered_At time.Time          `json:"ordered_on"  bson:"ordered_on"`
	Price        int                `json:"total_price" bson:"total_price"`
	Discount     int                `json:"discount"    bson:"discount"`
	Total_Tax    int                `json:"total_tax" bson:"total_tax"`
}
