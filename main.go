package main

import(
	"fmt"
	"github.com/gagiopapinni/balance-api-exercise/exchange"
	"github.com/gagiopapinni/balance-api-exercise/models"
	"github.com/gagiopapinni/balance-api-exercise/helper"
	"github.com/gin-gonic/gin"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
 )

var config helper.Config

func main() { 
	config.Load()

	db, err := sql.Open("mysql", config.DATA_SOURCE_NAME)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()


	exchangeRates, err := exchange.GetRates()
	if err!=nil {
		fmt.Println(err.Error())
		return
	}

	r := gin.Default()

	r.POST("/create-user", func(c *gin.Context) {

		obj := struct{Name string}{}

		err_bind := c.ShouldBindJSON(&obj)
		if err_bind != nil {
			c.JSON(400, gin.H{"error": err_bind.Error()})
			return
		}

		id, err := models.InsertUser(obj.Name, db)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
				
		c.JSON(200, gin.H{"result": gin.H{ "uid": id  } })

	})


	r.POST("/balance-operation", func(c *gin.Context) {
		obj := struct{
			Uid int
			Amount float32
			Note string
		}{}

		err_bind := c.ShouldBindJSON(&obj)
		if err_bind != nil {
			c.JSON(400, gin.H{"error": err_bind.Error()})
			return
		}
		fmt.Println(obj)
		if obj.Note == "" {
			c.JSON(400, gin.H{"error": "argument note missing" })
			return
		}

		err := models.BalanceOperation(obj.Uid, obj.Amount, db)
		if err != nil { 
			c.JSON(400, gin.H{"error": err.Error() })
			return 
		}
		
		models.InsertNote(obj.Uid, obj.Note, db)
				
		c.JSON(200, gin.H{"result": "ok" })
	})

	
	r.POST("/transaction", func(c *gin.Context) {
		obj := struct{
			From_uid int
			To_uid int
			Amount float32
			Note string
		}{}

		err_bind := c.ShouldBindJSON(&obj)
		if err_bind != nil {
			c.JSON(400, gin.H{"error": err_bind.Error()})
			return
		}
		
		if obj.Amount < 0 {
			c.JSON(400, gin.H{"error": "negative amount to transfer" })
			return
	
		}

		if obj.Note == "" {
			c.JSON(400, gin.H{"error": "argument note missing" })
			return
		}
		
		if !models.DoesUserExist(obj.To_uid, db) {
			c.JSON(400, gin.H{"error": "no such receiver user" })
			return
		}		

		err_withdraw := models.BalanceOperation(obj.From_uid, -obj.Amount, db)
		if err_withdraw != nil { 
			c.JSON(400, gin.H{"error": err_withdraw.Error() })
			return 
		}


		err_deposit := models.BalanceOperation(obj.To_uid, obj.Amount, db)
		if err_deposit != nil { 
			c.JSON(400, gin.H{"error": err_deposit.Error() })
			return 
		}

		
		models.InsertNote(obj.From_uid, obj.Note, db)
		models.InsertNote(obj.To_uid, obj.Note, db)
				
		c.JSON(200, gin.H{"result": "ok" })
	})

	r.GET("/notes", func(c *gin.Context) {
		obj := struct {
			Uid int  `Form:"uid"`
		}{}

		err_bind := c.ShouldBindQuery(&obj)
		if err_bind != nil {
			c.JSON(400, gin.H{"error": "bad request"})
			return 
		}

		res, err := models.BalanceNotes(obj.Uid, db)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"result": res})

	})

	r.GET("/balance", func(c *gin.Context) {
		obj := struct {
			Uid int  
			Currency string 
		}{}

		err_bind := c.ShouldBindQuery(&obj)
		if err_bind != nil {
			c.JSON(400, gin.H{"error": "bad request"})
			return 
		}

		if obj.Currency == "" { obj.Currency = "RUB" }
		rate, ok := exchangeRates[obj.Currency]	
		if !ok {
			c.JSON(400, gin.H{"error": "invalid currency"})
			return 
		}

		balance, err := models.Balance(obj.Uid, db)
		if err!=nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"result": balance*rate})	
	})
	r.Run(fmt.Sprintf(":%d",config.PORT)) 
}




























