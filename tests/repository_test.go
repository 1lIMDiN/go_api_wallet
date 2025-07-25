package tests

import (
	"context"
	"testing"

	"wallet-service/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	db   *sqlx.DB
	repo repository.WalletRepository
}

func (suite *RepositoryTestSuite) SetupSuite() {
	connStr := "host=localhost port=5432 user=wallet_user password=wallet_password dbname=wallet_test sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	require.NoError(suite.T(), err)

	suite.db = db
	suite.repo = repository.NewWallet(db)
}

func (suite *RepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *RepositoryTestSuite) SetupTest() {
	// Очистка таблицы перед каждым тестом
	suite.db.Exec("TRUNCATE TABLE wallets RESTART IDENTITY CASCADE")
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestCreateAndGetWallet() {
	walletID := uuid.New()

	// Создаем кошелек
	err := suite.repo.CreateWallet(context.Background(), walletID)
	assert.NoError(suite.T(), err)

	// Получаем кошелек
	wallet, err := suite.repo.GetWallet(context.Background(), walletID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), walletID, wallet.ID)
	assert.Equal(suite.T(), int64(0), wallet.Balance)
}

func (suite *RepositoryTestSuite) TestUpdateBalance() {
	walletID := uuid.New()
	err := suite.repo.CreateWallet(context.Background(), walletID)
	require.NoError(suite.T(), err)

	// Депозит
	err = suite.repo.UpdateBalance(context.Background(), walletID, 1000)
	assert.NoError(suite.T(), err)

	wallet, err := suite.repo.GetWallet(context.Background(), walletID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1000), wallet.Balance)

	// Снятие
	err = suite.repo.UpdateBalance(context.Background(), walletID, -500)
	assert.NoError(suite.T(), err)

	wallet, err = suite.repo.GetWallet(context.Background(), walletID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(500), wallet.Balance)

	// Недостаточно средств
	err = suite.repo.UpdateBalance(context.Background(), walletID, -600)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "insufficient funds")
}

func (suite *RepositoryTestSuite) TestConcurrentUpdates() {
	walletID := uuid.New()
	err := suite.repo.CreateWallet(context.Background(), walletID)
	require.NoError(suite.T(), err)

	// Запускаем 10 горутин для конкурентного обновления баланса
	const goroutines = 10
	const iterations = 100
	errs := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			var err error
			for j := 0; j < iterations; j++ {
				if err = suite.repo.UpdateBalance(context.Background(), walletID, 1); err != nil {
					break
				}
			}
			errs <- err
		}()
	}

	// Ждем завершения всех горутин
	for i := 0; i < goroutines; i++ {
		err := <-errs
		assert.NoError(suite.T(), err)
	}

	// Проверяем итоговый баланс
	wallet, err := suite.repo.GetWallet(context.Background(), walletID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(goroutines*iterations), wallet.Balance)
}
