package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"strconv"
	"tax-aggregator-service-demo/pkg/dbconn"
	"tax-aggregator-service-demo/tax/domain"
	"time"

	taxHandler "tax-aggregator-service-demo/tax/handler"
	taxRepository "tax-aggregator-service-demo/tax/repository"
	taxUsecase "tax-aggregator-service-demo/tax/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/urfave/cli/v2"

	"tax-aggregator-service-demo/config"
)

type Server *echo.Echo

var commands = []*cli.Command{
	{
		Name: "start",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "-p port will be used for application eg: -p 3000",
			},

			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "-c path will be used for config eg: -c ./config/config.json",
			},
		},
		Action: func(ctx *cli.Context) error {
			config := ctx.String("config")
			port := ctx.Int("port")
			return App(config, port)
		},
	},
}

func main() {
	app := &cli.App{
		Name:     "tax-aggregator-service",
		Version:  "1.0.1",
		Commands: commands,
	}
	if err := app.Run(os.Args); err != nil {
		log.Println("[app.main]:: error running application.")
		os.Exit(1)
	}
}

func App(cfg string, port int) error {
	e := echo.New()
	e.Use(middleware.Recover())
	e.HideBanner = true
	config, err := config.LoadConfig(cfg)
	if err != nil {
		return err
	}
	sourceDBConn, err := dbconn.NewMySQLDBConn(&config.SourceDatabase)
	if err != nil {
		return err
	}

	serviceDBConn, err := dbconn.NewPostgreSQLDBConn(&config.ServiceDatabase)
	if err != nil {
		return err
	}

	TaxRegistry(e, sourceDBConn, serviceDBConn, &config.PpnConfig)

	serverPort := ":" + strconv.Itoa(port)
	go func(){
		log.Println("starting tax-aggregator-service")
		if err := e.Start(serverPort); err != nil {
			log.Printf("[main.App]:: error starting server on port %s.\n", serverPort)
			e.Logger.Fatal("[main.App]:: shutting down application.")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 120 * time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	defer func() {
		log.Println("[main.App]:: closing source database connection...")
		if err := sourceDBConn.Close(); err != nil {
			log.Fatal("[main.App]:: error closing source database connection.")
		}

		log.Println("[main.App]:: closing service database conection...")
		if err := serviceDBConn.Close(); err != nil {
			log.Fatal("[main.App]:: error closing service database connection.")
		}
	}()
	return nil
}

func TaxRegistry(e Server, sourceDBConn, serviceDBConn *sql.DB, ppnConfig *config.PpnConfig) {
	taxRepository := taxRepository.NewTaxRepository(sourceDBConn, serviceDBConn)
	taxUsecase := taxUsecase.NewTaxUsecase(taxRepository, &domain.TaxConfig{
		TimeStartPpn: ppnConfig.TimeStartPpn,
		TimeStartPpnNew: ppnConfig.TimeStartPpnNew,
		TarifPpn: ppnConfig.TarifPpn,
		TarifPpnNew: ppnConfig.TarifPpnNew,
	})
	taxHandler := taxHandler.NewTaxHandler(taxUsecase)
	taxHandler.Routes(e)
}
