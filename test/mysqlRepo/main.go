package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	boil "github.com/volatiletech/sqlboiler/v4/boil"
	models "grpcTest/my_models"
	"log"
)

func main() {
	ctx := context.Background()
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	//err := http2.ConfigureTransport(tr)
	//if err != nil {
	//	return
	//}
	//client := &http.Client{Transport: tr}
	//bt := []byte(`{"name": "wangq"}`)
	//req, err := http.NewRequest("POST", "https://localhost:5001/v1/example/echo", bytes.NewReader(bt))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//req.Header.Add("Content-type", "application/grpc")
	//resp, err := client.Do(req)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(body))

	//return
	//pri, err := sm2.GenerateKey(rand.Reader)
	//aa, err := x509.WritePrivateKeyToPem(pri, nil)
	//err = ioutil.WriteFile("pri.pem", aa, 755)
	//bt, err := x509.WritePublicKeyToPem(&pri.PublicKey)
	//err = ioutil.WriteFile("pub.pem", bt, 755)
	//fmt.Println(pri)
	//bt, err = ioutil.ReadFile("./sm2PriKey.pem")
	//bb, err := smT.ReadPrivateKeyFromPem("./sm2PriKey.pem", nil)
	//fmt.Println(bb)
	//pri, err = x509.ReadPrivateKeyFromPem(bt, nil)
	//fmt.Println(pri)

	// Open handle to database like normal
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3307)/test?charset=utf8mb4&parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	//rsa.PrivateKey{}
	// If you don't want to pass in db to all generated methods
	// you can use boil.SetDB to set it globally, and then use
	// the G variant methods like so (--add-global-variants to enable)
	boil.SetDB(db)
	boil.DebugMode = true
	//var user models.User
	//user.Name = "李子明"
	//err = user.Insert(context.Background(), db, boil.Infer())
	//fmt.Println(err)
	////user, err := models.Users().One(context.Background(), db)
	//users, err := models.Users(qm.Where("name=?", "李子明")).All(context.Background(), db)
	//fmt.Printf("用户信息 %#v  错误信息 %v ", users, err)

	//for _, v := range users {
	//	fmt.Println(v.Name)
	//}

	//// Query all users
	users, err := models.Users().Count(ctx, db)
	fmt.Println(users, err)

	//// Panic-able if you like to code that way (--add-panic-variants to enable)
	//users := models.Users().AllP(db)
	//
	//// More complex query
	//users, err := models.Users(Where("age > ?", 30), Limit(5), Offset(6)).All(ctx, db)
	//
	//// Ultra complex query
	//users, err := models.Users(
	//	Select("id", "name"),
	//	InnerJoin("credit_cards c on c.user_id = users.id"),
	//	Where("age > ?", 30),
	//	AndIn("c.kind in ?", "visa", "mastercard"),
	//	Or("email like ?", `%aol.com%`),
	//	GroupBy("id", "name"),
	//	Having("count(c.id) > ?", 2),
	//	Limit(5),
	//	Offset(6),
	//).All(ctx, db)
	//
	//// Use any "boil.Executor" implementation (*sql.DB, *sql.Tx, data-dog mock db)
	//// for any query.
	//tx, err := db.BeginTx(ctx, nil)
	//if err != nil {
	//	return err
	//}
	//users, err := models.Users().All(ctx, tx)
	//
	//// Relationships
	//user, err := models.Users().One(ctx, db)
	//if err != nil {
	//	return err
	//}
	//movies, err := user.FavoriteMovies().All(ctx, db)
	//
	//// Eager loading
	//users, err := models.Users(Load("FavoriteMovies")).All(ctx, db)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(len(users.R.FavoriteMovies))
}
