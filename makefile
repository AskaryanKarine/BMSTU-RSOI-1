MIGRATION_DIR = "migration"

create-migration:
ifeq ($(name),)
	@echo "Введите название миграции с помощью параметра name"
else
	goose --dir=$(MIGRATION_DIR) create $(name) sql
endif

lint:
	golangci-lint run --timeout=5m

.PHONY : create-migration lint

