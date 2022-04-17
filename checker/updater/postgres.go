package updater

import (
	"context"
	"fmt"
	"log"

	"github.com/Abbas-gheydi/hanoproxy/model"
	"github.com/jackc/pgx/v4"
)

func (rec *recordTable) checkPostgresHealth(r model.Record) (updatedRecord model.Record) {

	for i := range r.Ip {
		go func(i int) {

			isSalve, err := psqlConnect(r.Ip[i].UserName,
				r.Ip[i].Password,
				r.Options.DbName,
				r.Ip[i].Addr,
				r.Ip[i].Port)
			if err != nil {

				rec.update(&r.Ip[i], false)
			} else {

				rec.update(&r.Ip[i], true)
			}
			switch isSalve {
			case true:
				if r.Options.MasterOnly {

					rec.update(&r.Ip[i], false)

				}
			case false:
				if r.Options.SlaveOnly {

					rec.update(&r.Ip[i], false)

				}

			}
		}(i)

	}

	return r

}

func psqlConnect(username string, password string, dbname string, ip string, port string) (isSlave bool, err error) {

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	//urlExample := "postgres://user:pass@localhost:5432/postgres"
	urlExample := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", username, password, ip, port, dbname)
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return

	}
	defer conn.Close(context.Background())

	err = conn.QueryRow(context.Background(), "select pg_is_in_recovery()").Scan(&isSlave)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		return

	}

	log.Println("postgres", ip, "connected successfully. isSlave:", isSlave)

	return

}
