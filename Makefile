MOCK_OUT=pkg/mock/pkg_gen.go
MOCK_SRC=./pkg/domain/interfaces
MOCK_INTERFACES=Slack Policy UseCases

all: mock

mock: $(MOCK_OUT)

$(MOCK_OUT): $(MOCK_SRC)/*
	go run github.com/matryer/moq@v0.5.1 -pkg mock -out $(MOCK_OUT) $(MOCK_SRC) $(MOCK_INTERFACES)

clean:
	rm -f $(MOCK_OUT)
