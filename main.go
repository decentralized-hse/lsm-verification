package main

import (
	"log"
	"lsm-verification/calculations"
	"lsm-verification/config"
	"lsm-verification/db"
	"lsm-verification/orchestrator"
	"path"
	"time"
)

func signLoop(orch orchestrator.Orchestrator, cfg config.Config) error {
	for {
		err := orch.SignNew()
		if err == nil {
			continue
		}
		if cfg.SkipErrors {
			log.Println("Skipping error error: ", err)
			continue
		}
		if err == orchestrator.ErrNoNewEntities {
			time.Sleep(time.Duration(cfg.SignTimeout) * time.Second)
		} else {
			log.Fatalln(err)
		}
	}
}

func validateDb(orch orchestrator.Orchestrator, cfg config.Config) (bool, *string, error) {
	lastValidated, lastHash, err := orch.ValidateFromLseq(nil, nil)
	for {
		if err == orchestrator.ErrNoNewEntities {
			return true, lastValidated, nil
		}
		if err == orchestrator.ErrValidationFailed {
			return false, lastValidated, err
		}
		if err != nil {
			if !cfg.SkipErrors {
				return false, lastValidated, err
			}
			log.Println("Skipping error error: ", err)
		}
		lastValidated, lastHash, err = orch.ValidateFromLseq(lastValidated, lastHash)
	}
}

func main() {
	cfg := config.LoadConfig(path.Join("config", "config.yaml"))
	dbState, err := db.CreateDbState(cfg)
	if err != nil {
		log.Fatalln("Failed to load db: ", err)
	}
	defer dbState.CloseConnection()

	hashCalculator := calculations.CreateHashCalculator()
	orch := orchestrator.CreateOrchestrator(dbState, hashCalculator)
	log.Println("Running in mode: ", cfg.RunMode)
	if cfg.RunMode == config.RunModeValidation {
		valid, lastLseq, err := validateDb(orch, cfg)
		if err != nil {
			log.Fatalln(err)
		}
		if valid {
			log.Println("Database is valid to lseq: ", lastLseq)
		} else {
			log.Println("Database is not valid on lseq: ", lastLseq)
		}
	} else if cfg.RunMode == config.RunModeSign {
		err = signLoop(orch, cfg)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln("Run mode unsupported")
	}
	log.Println("Task is done")
}
