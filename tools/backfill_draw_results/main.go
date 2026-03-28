package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"go-admin/config"
	bizsvc "go-admin/internal/services/biz"

	gormx "github.com/wangyahua6688-maker/tk-common/utils/dbx/gormx"
)

func main() {
	configPath := flag.String("config", "config.yaml", "config file path")
	specialLotteryID := flag.Uint("special-lottery-id", 0, "only backfill the specified special lottery id")
	force := flag.Bool("force", false, "rebuild all records, not only incomplete ones")
	flag.Parse()

	config.Init(*configPath)
	cfg := config.GetConfig()

	dbCfg := gormx.DefaultDBConfig()
	dbCfg.DSN = cfg.Database.DSN
	dbCfg.LogLevel = gormx.GormLogLevelFromString(cfg.Database.LogLevel)

	db, err := gormx.NewMySQLDB(dbCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect database failed: %v\n", err)
		os.Exit(1)
	}

	service := bizsvc.NewLotteryService(db)
	count, err := service.BackfillDrawResultData(context.Background(), uint(*specialLotteryID), *force)
	if err != nil {
		fmt.Fprintf(os.Stderr, "backfill draw results failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("backfilled %d draw records\n", count)
}
