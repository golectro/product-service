package command

import (
	"os"
	"strings"

	"slices"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CommandExecutor struct {
	DB *gorm.DB
}

func NewCommandExecutor(db *gorm.DB) *CommandExecutor {
	return &CommandExecutor{DB: db}
}

func (ce *CommandExecutor) Execute(logger *logrus.Logger) bool {
	args := os.Args[1:]
	run := true
	hasFlag := false

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			hasFlag = true
			switch arg {
			// case "--migrate":
			// 	ce.handleMigrate(logger)
			// case "--seed":
			// 	ce.handleSeed(logger)
			case "--run":
				run = true
			}
		}
	}

	if hasFlag && !contains(args, "--run") {
		run = false
	}

	return run
}

func contains(slice []string, target string) bool {
	return slices.Contains(slice, target)
}

// func (ce *CommandExecutor) handleMigrate(logger *logrus.Logger) {
// 	if err := migrations.Migrate(ce.DB); err != nil {
// 		logger.Fatalf("❌ Migration failed: %v", err)
// 	}
// 	logger.Println("✅ Migration completed")
// }

// func (ce *CommandExecutor) handleSeed(logger *logrus.Logger) {
// 	if err := migrations.Seeder(ce.DB, logger); err != nil {
// 		logger.Fatalf("❌ Seeder failed: %v", err)
// 	}
// 	logger.Println("✅ Seeder completed")
// }
