package repository

import (
	"e-wallet/config"
	"e-wallet/models"
	"fmt"
	"time"
)

// TopupBalanceRepository function inisiasi
func TopupBalanceRepository(username string, jumlah int) (err error) {
	db, err := config.Connect()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()

	// deklarasi query dengan menggunakan "dbPrepare", bertujuan agar query ini bisa dipakai berulang kali (re-usable)
	stmtGetUserInit, _ := db.Prepare("select id from users where username = ? order by created_at desc limit 1")
	stmtGetBalance, _ := db.Prepare("select id, user_id, balance from user_balance where user_id = ? order by created_at desc limit 1")
	stmtUpdateBalance, _ := db.Prepare("update user_balance set balance = ? where user_id = ?")
	stmtRecordBalanceHistory, _ := db.Prepare("insert into user_balance_history(user_balance_id, balance_before, balance_after, amount, activity, location, user_agent, author, type, created_at) values(?,?,?,?,?,?,?,?,?,?)")
	stmtCreateNewUserBalance, _ := db.Prepare("insert into user_balance(user_id, balance, created_at) values (?,?,?)")
	now := time.Now()

	// deklarasi variable untuk menampung data dari operasi mengambil id user, dan mengambil detail balance user
	var getUser = models.GetUserID{}
	var getBalance = models.GetBalance{}

	// operasi untuk mengeksekusi query yang sudah di deklarasi diatas dengan parameter yang tertera dibawah ini
	// kenapa menggunakan "QueryRow" ?, hal ini bertujuan agar hasil yang didapatkan setelah query hanya 1 record data saja, karena memang yang dibutuhkan hanya 1 data saja dan data yang diambil adalah data yang paling baru
	// data yang diambil berupa ID dari table "users" dan disimpan pada variable "getUser"
	stmtGetUserInit.QueryRow(username).Scan(&getUser.ID)
	// data yang diambil berupa properti yang tertera dari table "user_balance" dan disimpan pada variable "getBalance"
	stmtGetBalance.QueryRow(getUser.ID).Scan(&getBalance.ID, &getBalance.UserID, &getBalance.Balance)

	// deklarasi insertID untuk menampung ID yang direturn pada saat melakukan insert
	var insertID int64

	// statment ini berguna untuk melakukan checkpoint kepada value yang ada, untuk mencegah operasi yang tidak diperlukan
	// statment ini juga digunakan untuk mendeteksi apakah data user baru atau tidak pada table _user_balance
	if getBalance.UserID == 0 {
		// eksekusi query "stmtCreateNewUserBalance" untuk membuat data user baru jika data user pada table "user_balance" blm tersedia
		res, err := stmtCreateNewUserBalance.Exec(getUser.ID, jumlah, now)
		if err != nil {
			println("Error while exec")
			println(err.Error())
			return err
		}

		// operasi untuk mendapatkan balikan id setelah proses insert berhasil, kenapa diambil ?, karena saya membutuhkan data tersebut daripada harus melakukan query lagi hanya untuk mengambil data id
		insertID, err = res.LastInsertId()
		if err != nil {
			println("Error while get insert ID")
			println(err.Error())
			return err
		}
	}

	// statment untuk mengetahui data yang dibawa merupakan data baru atau tidak
	// jika data user baru maka insertID akan terdapat isi berupa id data user, jika kosong maka data yang dibawa merupakan data yang sudah memiliki entitas di table "user_balance"
	if insertID != 0 {
		// operasi penjumlahan untuk mendapatkan hasil berupa angka setelah penjumlahan, yang dijumlahkan adalah balance asli + nominal yang di top-up
		balanceAfter := getBalance.Balance + jumlah
		// eksesuksi query untuk membuat record atau jejak transaksi
		_, err = stmtRecordBalanceHistory.Exec(insertID, getBalance.Balance, balanceAfter, jumlah, "top-up", "Yogyakarta", username, username, "debit", now)

		if err != nil {
			println("Error record history")
			println(err.Error())
			return err
		}
	} else {
		// prinsipnya sama dengan operasi sebelumnya, hanya saja operasi ini dilakukan dengan syarat data user sudah ada di table "user_balance" sebelumnya
		topupBalance := getBalance.Balance + jumlah
		_, err = stmtUpdateBalance.Exec(topupBalance, getUser.ID)
		if err != nil {
			println("Error while excute set balance")
			println(err.Error())
			return err
		}

		// operasi untuk melakukan pencatatan transaksi
		balanceAfter := getBalance.Balance + jumlah
		_, err = stmtRecordBalanceHistory.Exec(getBalance.ID, getBalance.Balance, balanceAfter, jumlah, "top-up", "Yogyakarta", username, username, "debit", now)

		if err != nil {
			println("Error record history")
			println(err.Error())
			return err
		}
	}
	fmt.Println("Insert success!")
	return err
}
