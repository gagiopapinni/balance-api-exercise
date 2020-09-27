package models 


import ( 
	//"fmt" 
	"time"
	"errors"
	"database/sql"
)


func DoesUserExist(id int, db *sql.DB) bool {
	if id<0 { return false }	
	var res int
	err := db.QueryRow("select 1 from user where id = ?", id).Scan(&res)
	if err != nil { return false }
	return true
}

func balanceId(uid int, db *sql.DB) int {	
	if uid<0 { return -1 }	
	var id int64
	err := db.QueryRow("select id from balance where uid = ?", uid).Scan(&id)
	if err != nil { return -1 }
	return int(id)
}

func createBalance(uid int, value float32, db *sql.DB) (int, error) {
	stmt, err := db.Prepare("INSERT INTO balance(uid, value) VALUES(?,?)")
	if err != nil {
		return 0, errors.New("unknown error")
	}

	res, err := stmt.Exec(uid,value)
	if err != nil {
		return 0, errors.New("unknown error")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New("unknown error")
	}

	return int(id), nil
}


func currentBalanceValue(uid int, db *sql.DB) (float32, error) {
	var res float32
	err := db.QueryRow("select value from balance where uid = ?", uid).Scan(&res)
	if err != nil { 
		return 0, errors.New("unknown error")
	}
	return res, nil
}

func updateBalanceValue(uid int, value float32, db *sql.DB) error {
	stmt, err := db.Prepare("UPDATE balance SET value=? WHERE uid = ?")
	if err != nil {
		return  errors.New("unknown error")
	}

	_, err = stmt.Exec(value, uid)
	if err != nil {
		return  errors.New("unknown error")
	}

	return nil
}

func Balance(uid int, db *sql.DB ) (float32, error){
	if !DoesUserExist(uid, db) { return 0, errors.New("no such user") }
	bal, err := currentBalanceValue(uid, db)
	if err != nil { return 0, err }
	return bal, nil
}

func BalanceOperation(uid int, amount float32, db *sql.DB) error {
	//if amount == 0 { return nil }
	if !DoesUserExist(uid, db) { return errors.New("no such user") }

	bid := balanceId(uid, db)
	if bid>0 {

		value, err := currentBalanceValue(uid, db)
		if err!=nil { return err }
		if value+amount<0 { return errors.New("Not enough money on balance")}

		err_update := updateBalanceValue(uid, value+amount, db)
		if err_update != nil { return err_update }

	} else {
		if amount < 0 { 
			return errors.New("Not enough money on balance") 
		} else {
			_, err := createBalance(uid, amount, db)
			if err!=nil { return err }
		}
	}

	return nil

}


func BalanceNotes(uid int, db *sql.DB) (interface{}, error) {
	
	if !DoesUserExist(uid, db) { return nil, errors.New("no such user") }

	result := make([]interface{}, 0)

	bid := balanceId(uid, db)
	if bid==-1 { return result, nil }

	rows, err := db.Query(`SELECT timestamp, text FROM note 
			       WHERE bid = ? 
			       ORDER BY timestamp`, bid)
	if err != nil {
		return nil, errors.New("unknown error")
	}
	defer rows.Close()

	for rows.Next() {
		var timestamp int
		var text string
		err := rows.Scan(&timestamp, &text)
		if err != nil {
			return nil, errors.New("unknown error")
		}
		json := map[string]interface{}{
			"timestamp": timestamp,
			"text": text,
		}
		result = append(result,json)
	}

	return result, nil

}

func InsertNote(uid int, text string, db *sql.DB) (int, error) {
	if len(text) == 0 {
		return 0, errors.New("missing note argument")
	}

	stmt, err := db.Prepare("INSERT INTO note(bid, text, timestamp) VALUES(?, ?, ?)")
	if err != nil {
		return 0, errors.New("unknown error")
	}

	timestamp :=  time.Now().Unix()
	bid := balanceId(uid, db)
	if bid==-1 { 
		return 0, errors.New("no balance account")
	}

	res, err := stmt.Exec(bid, text, timestamp)
	if err != nil {
		return 0, errors.New("unknown error")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New("unknown error")
	}

	return int(id), nil
}

func InsertUser(name string, db *sql.DB ) (int, error) {
	if len(name) == 0 {
		return 0, errors.New("missing name argument")
	}

	stmt, err := db.Prepare("INSERT INTO user(name) VALUES(?)")
	if err != nil {
		return 0, errors.New("unknown error")
	}

	res, err := stmt.Exec(name)
	if err != nil {
		return 0, errors.New("unknown error")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New("unknown error")
	}

	return int(id), nil
}


