mock:
	mockgen -source=./service/manager.go -destination=./service/mock/mock.go *
	mockgen -source=./repository/manager.go -destination=./repository/mock/mock.go *

.PHONY: mock