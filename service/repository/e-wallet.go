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

// TransferBalanceRepository function inisiasi
func TransferBalanceRepository(username, tujuan string, jumlah int) (err error) {
	db, err := config.Connect()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()

	// deklarasi query yang akan digunakan, dengan "db.Prepare" memungkinkan untuk query dipakai berulang kali
	stmtGetUserInit, _ := db.Prepare("select id from users where username = ? order by created_at desc limit 1")
	stmtGetBalance, _ := db.Prepare("select id, user_id, balance from user_balance where user_id = ? order by created_at desc limit 1")
	stmtUpdateBalance, _ := db.Prepare("update user_balance set balance = ? where user_id = ?")
	stmtRecordBalanceHistory, _ := db.Prepare("insert into user_balance_history(user_balance_id, balance_before, balance_after, amount, activity, location, user_agent, author, type, created_at) values(?,?,?,?,?,?,?,?,?,?)")
	stmtCreateNewUserBalance, _ := db.Prepare("insert into user_balance(user_id, balance, created_at) values (?,?,?)")
	now := time.Now()

	// deklarasi variable untuk menampung data dalam pemanggilan query selanjutnya
	var getUserMe = models.GetUserID{}
	var getUserTujuan = models.GetUserID{}
	var getBalanceMe = models.GetBalance{}
	var getBalanceTujuan = models.GetBalance{}

	// eksekusi query untuk mengambil data ID user sesuai dengan parameter "username", kemudian data disimpan di variable "getUserMe" yang sudah di deklarasi diatas
	stmtGetUserInit.QueryRow(username).Scan(&getUserMe.ID)
	// eksekusi query untuk mengambil data ID user yang menjadi tujuan "transfer" sesuai dengan parameter "username", kemudian data disimpan di variable "getUserTujuan" yang sudah di deklarasi diatas
	stmtGetUserInit.QueryRow(tujuan).Scan(&getUserTujuan.ID)

	// statment untuk melakukan cek apakah data user tujuan sudah ada atau belum di table "user_balance", jika nilainya != 0 berarti datanya sudah ada, jika == 0  maka datanya blm ada
	// jika datanya sudah ada / != 0, maka akan masuk ke statment dibawah
	if getUserTujuan.ID != 0 {
		// eksekusi query untuk mengambil data detail dari data balance yang ada pada "user_balance"
		// dalam perintah ini akan emngambil data user pengirim dan disimpan ke variable "getBalanceMe" yang sudah di deklarasi diatas dengan properti yang tertera
		stmtGetBalance.QueryRow(getUserMe.ID).Scan(&getBalanceMe.ID, &getBalanceMe.UserID, &getBalanceMe.Balance)
		// dalam perintah ini akan emngambil data user penerima dan disimpan ke variable "getBalanceTujuan" yang sudah di deklarasi diatas dengan properti yang tertera
		stmtGetBalance.QueryRow(getUserTujuan.ID).Scan(&getBalanceTujuan.ID, &getBalanceTujuan.UserID, &getBalanceTujuan.Balance)

		// statment pengecekan untuk mengecek saldo dari pengirim apakah saldo nya lebih banyak dari jumlah yang akan dikirimkan atau sebaliknya
		// syarat untuk bisa mengirim saldo atau balance adalah saldo / balance si pengirim harus lebih banyak atau sama dari jumlah yang akan dikirimkan (>=)
		// jika saldo >= jumlah yang akan dikirimkan maka akan masuk pada statment dibawah, jika syarat tidak terpenuhi maka proses akan diterminate atau mengeksekusi else kemudian selesai
		if getBalanceMe.Balance >= jumlah {
			// statment untuk mengetahui data yang dibawa merupakan data baru atau tidak
			// jika data user baru maka insertID akan terdapat isi berupa id data user, jika kosong maka data yang dibawa merupakan data yang sudah memiliki entitas di table "user_balance"
			var insertID int64
			// statment untuk mengecek apakah userID dari penerima transfer sudah ada datanya pada table "user_balance" atau belum, jika datanya sudah ada maka "getBalanceTujuan.UserID" berisi userID dari penerima, jika datanya belum ada di database maka "getBalanceTujuan.UserID" = 0
			// jika "getBalanceTujuan.UserID" kosong atau =0, maka artinya akan membuat data baru untuk si penerima transfer pada table "user_balance"
			if getBalanceTujuan.UserID == 0 {
				// perintah untuk membuat data balance baru, dalam hal ini membuat data balance untuk penerima transfer dengan parameter dibawah
				res, err := stmtCreateNewUserBalance.Exec(getUserTujuan.ID, jumlah, now)
				if err != nil {
					println("Error while exec")
					println(err.Error())
					return err
				}

				// pada saat insert atau membuat data baru maka akan ada nilai bailkan yang dapat ditangkap yaitu "ID" dari data yang dibuat
				// untuk menangkap data tersebut akan diproses dengan peritah dibawah, dan dimasukan ke variable "insertID"
				insertID, err = res.LastInsertId()
				if err != nil {
					println("Error while get insert ID")
					println(err.Error())
					return err
				}
			}
			// jika data "getBalanceTujuan.UserID" != 0 berarti merupakan user baru
			// jika data "getBalanceTujuan.UserID" == 0 berarti merupakan user yang datanya sudah ada ditabase
			if insertID != 0 {
				// pada scope ini data yang diproses pasti data user baru sebagai penerima dan pasti baru saja melakukan insert data baru ke table "user_balance", karena variable "insertID" != 0
				// perintah untuk menghitung balance setelah menerima transfer, maka rumusnya "saldo sebelumnya + jumlah transfer"
				balanceAfter := getBalanceTujuan.Balance + jumlah
				// melakukan pencatatan table "user_balance_history" untuk dicatat aktivitas transaksinya
				_, err = stmtRecordBalanceHistory.Exec(insertID, getBalanceTujuan.Balance, balanceAfter, jumlah, "transfer in", "Yogyakarta", username, tujuan, "debit", now)
				if err != nil {
					println("Error record history")
					println(err.Error())
					return err
				}
			} else {
				// jika "inserID" == 0 yang berarti data pada table "user_balance" sudah ada sebelumnya, dan pada scope ini akan melakukan update pada balance dimana balance dari penerima transfer akan ditambahkan dengan nominal yang di transfer
				balanceAdd := getBalanceTujuan.Balance + jumlah
				_, err = stmtUpdateBalance.Exec(balanceAdd, getUserTujuan.ID)
				if err != nil {
					println("Error while excute set balance")
					println(err.Error())
					return err
				}

				// perintah ini digunakan untuk mencatat transaksi pada table "user_balance_history" dimana pada proses sebelumnya telah melakukan update balance di table "user_balance"
				balanceAfter := getBalanceTujuan.Balance + jumlah
				_, err = stmtRecordBalanceHistory.Exec(getBalanceTujuan.ID, getBalanceTujuan.Balance, balanceAfter, jumlah, "transfer in", "Yogyakarta", username, tujuan, "debit", now)

				if err != nil {
					println("Error record history")
					println(err.Error())
					return err
				}
			}
			// pada scope ini akan mencatat transaksi untuk data pengirim
			// karena pengirim melakukan pengiriman saldo maka balancenya akan dikurangi dengan nominal yang dikirimkan, sehingga menghasilkan nilai yang disimpan pada "balanceAfterTransfer" seperti pada code dibawah
			balanceAfterTransfer := getBalanceMe.Balance - jumlah
			// melakukan eksekusi query untuk mencatat transaksi, transaksi dalam scope ini adalah pengirim mengirimkan saldo kepada penerima, dengan properti seperti yang tertera dibawah
			_, err = stmtRecordBalanceHistory.Exec(getBalanceMe.ID, getBalanceMe.Balance, balanceAfterTransfer, jumlah, "transfer out", "Yogyakarta", username, tujuan, "credit", now)

			if err != nil {
				println("Error record history")
				println(err.Error())
				return err
			}

			// setelah record transaksi berhasil dicatat, maka perintah dibawah ini berguna untuk melakukan update balance pada data "user_balance" milik pengirim
			balanceOff := getBalanceMe.Balance - jumlah
			_, err = stmtUpdateBalance.Exec(balanceOff, getUserMe.ID)
			if err != nil {
				println("Error while excute set balance")
				println(err.Error())
				return err
			}
		} else {
			println("Not enough balance")
		}
	} else {
		println("Username tujuan tidak dapat ditemukan")
	}

	fmt.Println("Insert success!")

	return err
}

// GetUser function
func GetUserData(username string) (data models.GetUserData, err error) {
	db, err := config.Connect()
	if err != nil {
		fmt.Println(err.Error())
		return data, err
	}
	defer db.Close()

	// deklarasi query dengan menggunakan "dbPrepare", bertujuan agar query ini bisa dipakai berulang kali (re-usable)
	stmtGetUserData, err := db.Prepare("select username, password from users where username = ?")
	if err != nil {
		println(err.Error())
		return data, err
	}
	stmtGetUserData.QueryRow(username).Scan(&data.Username, &data.Password)
	return data, err
}
