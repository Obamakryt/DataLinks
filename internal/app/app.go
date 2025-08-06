package app

import (
	"DataLinks/internal/api/rest"
	"DataLinks/internal/api/rest/handler"
	"DataLinks/internal/service"
	"DataLinks/internal/service/jwt_hash"
	"DataLinks/internal/slogger"
	"DataLinks/internal/storages/postgreSQL/connect_pool"
	"DataLinks/internal/storages/postgreSQL/storage_crud"
	"DataLinks/pkg/config/godotenv"
	"DataLinks/pkg/config/parser"
	"github.com/go-playground/validator/v10"
	"log"
	"log/slog"
)

func RunApp() {
	err := godotenv.LoadEnv("./pkg/secrets.env")
	if err != nil {
		log.Fatal(err.Error())
	}
	DBDataStruct := connect_pool.NewParseStruct("./pkg/cfg.yaml", "./pkg/secrets.env")
	err = DBDataStruct.ParseConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	urlConect := connect_pool.CreatePgxUrl(DBDataStruct.PostgresUrl)
	//////// logger
	deflogger := slog.Logger{}
	slogSetup := slogger.Setup{Logger: &deflogger}
	slogSetup.SetupLogger(0)
	/////////
	pool, err := connect_pool.NewPool(urlConect, 3, slogSetup)
	if err != nil {
		log.Fatal(err.Error())
	}
	GlobalPool := storage_crud.NewPostgresPool(pool)
	/////////// service
	Registration := service.LogicReg{Logger: slogSetup, Storage: storage_crud.AuthReg{GlobalPool}}
	Authorization := service.LogicAuth{Logger: slogSetup, Storage: storage_crud.AuthReg{GlobalPool}}
	AddLink := service.LogicLinkAdd{Logger: slogSetup, Storage: storage_crud.NewLinks{GlobalPool, GlobalPool}}
	TakeChart := service.LogicTakeChart{Logger: slogSetup, Storage: storage_crud.UserLinks{GlobalPool}}
	Swap := service.LogicUpdateLink{Logger: slogSetup, Storage: storage_crud.UpdateLink{GlobalPool}}
	DeleteLink := service.LogicDeleteLink{Logger: slogSetup, Storage: storage_crud.DeleteLink{GlobalPool}}
	//// handler
	defValidator := validator.New()
	handler.PasswordValidator(defValidator)
	handler.UrlValidator(defValidator)
	handler.HTTPValidator(defValidator)

	ServiceStoreAuthReg := storage_crud.HandlerRegAuthStorage{Registration: &Registration, Login: &Authorization}

	HandlerLogAuth := handler.AuthRegHandler{Validator: defValidator, Logger: slogSetup, Storage: ServiceStoreAuthReg}

	ServiceStoreLinks := storage_crud.HandlerStorage{
		ServiceNewLinks:   &AddLink,
		ServiceUserLinks:  &TakeChart,
		ServiceDeleteLink: &DeleteLink,
		ServiceUpdateLink: &Swap}

	HandlerLinks := handler.LinkHandler{Validator: defValidator, Logger: slogSetup, Storage: ServiceStoreLinks}
	//////// server

	ServerSetings := parser.NewServerSettings("./pkg/cfg.yaml")
	err = ServerSetings.ParserServerSet()
	if err != nil {
		log.Fatal(err.Error())
	}
	sigh := jwt_hash.JWTSigh{}
	err = sigh.CreateSigh()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = rest.Router(ServerSetings, sigh, HandlerLogAuth, HandlerLinks)
	if err != nil {
		log.Fatal(err.Error())
	}
}
