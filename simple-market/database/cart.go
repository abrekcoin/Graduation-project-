package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"market/models"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": productID})

	if err != nil {
		log.Println(err)
		fmt.Println("eror")
		return errors.New("can't find product")
		
	}
	var productcart []models.ProductUser
	err = searchfromdb.All(ctx, &productcart)
	if err != nil {
		log.Println(err)
		return errors.New("can't find product")
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("user is not valid")
	}
// Tax Calculation
	var taxRatio = productcart[0].Price
	var productPrice = productcart[0].Tax_Ratio
	var productTaxValue = (taxRatio * productPrice) / 100
	productcart[0].Tax_Value = productTaxValue

//Updating cart

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productcart}}}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return errors.New("cannot add product to cart")
	}
	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("user is not valid")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return errors.New("cannot add product to cart")
	}
	return nil

}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("cannot remove item from cart")
	}

	var getcartitems models.User
	var ordercart models.Order
	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Orderered_At = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"},
		{Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

	grouping2 := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"},
		{Key: "total_tax", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.tax_value"}}}}}}

	currentresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	taxresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping2})
	ctx.Done()
	if err != nil {
		panic(err)
	}
	var getusercart []bson.M
	var gettaxcart []bson.M

	if err = currentresults.All(ctx, &getusercart); err != nil {
		panic(err)
	}

	if err = taxresults.All(ctx, &gettaxcart); err != nil {
		panic(err)
	}

	var total_price int32
	var total_tax int32

	for _, user_item := range getusercart {
		price := user_item["total"]
		total_price = price.(int32)
	}

	for _, tax_item := range gettaxcart {
		fmt.Println(tax_item["total_tax"])
		total_tax_json := tax_item["total_tax"]
		total_tax = total_tax_json.(int32)
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getcartitems)
	if err != nil {
		log.Println(err)
	}

	ordercart.Price = int(total_price)
	ordercart.Total_Tax = int(total_tax)

	// For 4th orders based on tax ratio
	discountByOrder := 0

	if (len(getcartitems.Order_Status)+1)%4 == 0 {
		for _, item := range getcartitems.UserCart {
			if item.Tax_Ratio == 8 {
				fmt.Println("%8")
				discountByOrder += item.Price * 10 / 100
			} else if item.Tax_Ratio == 18 {
				fmt.Println("%18")
				discountByOrder += item.Price * 15 / 100
			} else {
				fmt.Println("%1")
			}
		}
	}

	// For more than 3 items
	discountByCount := 0

	cartMap := make(map[string]int)
	priceMap := make(map[string]int)

	for _, item := range getcartitems.UserCart {
		if cartMap[item.Product_ID.String()] < 1 {
			cartMap[item.Product_ID.String()] = 1
			priceMap[item.Product_ID.String()] = item.Price
		} else {
			cartMap[item.Product_ID.String()] += 1
		}
	}

	for key, value := range cartMap {
		if value > 3 {
			discountByCount += (priceMap[key] * (value - 3)) * 8 / 100
		}
	}

	totalDiscount := 0

	if discountByCount > discountByOrder {
		totalDiscount = discountByCount
	} else {
		totalDiscount = discountByOrder
	}

	ordercart.Discount = totalDiscount
	ordercart.Price -= ordercart.Discount

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordercart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getcartitems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	usercart_empty := make([]models.ProductUser, 0)
	filtered := bson.D{primitive.E{Key: "_id", Value: id}}
	updated := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}
	_, err = userCollection.UpdateOne(ctx, filtered, updated)
	if err != nil {
		return errors.New("cannot update the purchase")

	}
	return nil
}
