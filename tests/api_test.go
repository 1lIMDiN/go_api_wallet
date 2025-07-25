package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"wallet-service/internal/api"
	"wallet-service/internal/models"
	"wallet-service/internal/repository"
	"wallet-service/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	db       *sqlx.DB
	repo     repository.WalletRepository
	service  service.WalletService
	router   *mux.Router
	walletID uuid.UUID
}

func (suite *APITestSuite) SetupSuite() {
	connStr := "host=localhost port=5432 user=wallet_user password=wallet_password dbname=wallet_test sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	require.NoError(suite.T(), err)

	suite.db = db
	suite.repo = repository.NewWallet(db)
	suite.service = service.NewWalletService(suite.repo)
	handler := api.NewHandler(suite.service)

	suite.router = mux.NewRouter()
	suite.router.HandleFunc("/api/v1/wallets/{walletId}", handler.GetWallet).Methods("GET")
	suite.router.HandleFunc("/api/v1/wallet", handler.PostWallet).Methods("POST")

	suite.walletID = uuid.New()
}

func (suite *APITestSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *APITestSuite) SetupTest() {
	// Очистка таблицы перед каждым тестом
	suite.db.Exec("TRUNCATE TABLE wallets RESTART IDENTITY CASCADE")
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (suite *APITestSuite) TestCreateAndGetWallet() {
	// Тест депозита
	op := models.WalletOperation{
		WalletID:      suite.walletID,
		OperationType: models.Deposit,
		Amount:        1000,
	}

	body, err := json.Marshal(op)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(body))
	require.NoError(suite.T(), err)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var wallet models.Wallet
	err = json.NewDecoder(rr.Body).Decode(&wallet)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.walletID, wallet.ID)
	assert.Equal(suite.T(), int64(1000), wallet.Balance)

	// Тест получения кошелька
	req, err = http.NewRequest("GET", "/api/v1/wallets/"+suite.walletID.String(), nil)
	require.NoError(suite.T(), err)

	rr = httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var retrievedWallet models.Wallet
	err = json.NewDecoder(rr.Body).Decode(&retrievedWallet)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), wallet, retrievedWallet)
}

func (suite *APITestSuite) TestWithdraw() {
	// Сначала создаем кошелек с балансом
	err := suite.repo.CreateWallet(context.Background(), suite.walletID)
	require.NoError(suite.T(), err)
	err = suite.repo.UpdateBalance(context.Background(), suite.walletID, 1000)
	require.NoError(suite.T(), err)

	// Тест снятия средств
	op := models.WalletOperation{
		WalletID:      suite.walletID,
		OperationType: models.Withdraw,
		Amount:        500,
	}

	body, err := json.Marshal(op)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(body))
	require.NoError(suite.T(), err)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var wallet models.Wallet
	err = json.NewDecoder(rr.Body).Decode(&wallet)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.walletID, wallet.ID)
	assert.Equal(suite.T(), int64(500), wallet.Balance)
}

func (suite *APITestSuite) TestInsufficientFunds() {
	err := suite.repo.CreateWallet(context.Background(), suite.walletID)
	require.NoError(suite.T(), err)

	op := models.WalletOperation{
		WalletID:      suite.walletID,
		OperationType: models.Withdraw,
		Amount:        500,
	}

	body, err := json.Marshal(op)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(body))
	require.NoError(suite.T(), err)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "insufficient funds")
}
