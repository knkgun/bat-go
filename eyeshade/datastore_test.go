package eyeshade

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brave-intl/bat-go/datastore/grantserver"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DatastoreMockTestSuite struct {
	suite.Suite
	ctx  context.Context
	db   Datastore
	mock sqlmock.Sqlmock
}

// type DatastoreTestSuite struct {
// 	suite.Suite
// 	ctx context.Context
// 	db  Datastore
// }

func TestDatastoreMockTestSuite(t *testing.T) {
	suite.Run(t, new(DatastoreMockTestSuite))
	// if os.Getenv("EYESHADE_DB_URL") != "" {
	// 	suite.Run(t, new(DatastoreTestSuite))
	// }
}

func (suite *DatastoreMockTestSuite) SetupSuite() {
	ctx := context.Background()
	// setup mock DB we will inject into our pg
	mockDB, mock, err := sqlmock.New()
	suite.Require().NoError(err, "failed to create a sql mock")

	name := "sqlmock"
	suite.db = NewFromConnection(&grantserver.Postgres{
		DB: sqlx.NewDb(mockDB, name),
	}, name)
	suite.ctx = ctx
	suite.mock = mock
}

func SetupMockGetAccountEarnings(
	mock sqlmock.Sqlmock,
	options AccountEarningsOptions,
) []AccountEarnings {
	getRows := sqlmock.NewRows(
		[]string{"channel", "earnings", "account_id"},
	)
	rows := []AccountEarnings{}
	for i := 0; i < options.Limit; i++ {
		accountID := fmt.Sprintf("publishers#uuid:%s", uuid.NewV4().String())
		earnings := decimal.NewFromFloat(
			float64(rand.Intn(100)),
		).Div(
			decimal.NewFromFloat(10),
		)
		channel := uuid.NewV4().String()
		rows = append(rows, AccountEarnings{channel, earnings, accountID})
		// append sql result rows
		getRows = getRows.AddRow(
			channel,
			earnings,
			accountID,
		)
	}
	mock.ExpectQuery(`
select
	channel,
	(.+) as earnings,
	account_id
from account_transactions
where account_type = 'owner'
	and transaction_type = (.+)
group by (.+)
order by earnings (.+)
limit (.+)`).
		WithArgs(options.Type[:len(options.Type)-1], options.Limit).
		WillReturnRows(getRows)
	return rows
}

func SetupMockGetAccountSettlementEarnings(
	mock sqlmock.Sqlmock,
	options AccountSettlementEarningsOptions,
) []AccountSettlementEarnings {
	getRows := sqlmock.NewRows(
		[]string{"channel", "paid", "account_id"},
	)
	rows := []AccountSettlementEarnings{}
	args := []driver.Value{
		fmt.Sprintf("%s_settlement", options.Type[:len(options.Type)-1]),
		options.Limit,
	}
	targetTime := options.StartDate
	if targetTime == nil {
		now := time.Now()
		targetTime = &now
	} else {
		args = append(args, targetTime)
	}
	untilDate := options.UntilDate
	if untilDate == nil {
		untilDatePrep := targetTime.Add(time.Hour * 24 * time.Duration(options.Limit))
		untilDate = &untilDatePrep
	} else {
		args = append(args, untilDate)
	}
	for i := 0; i < options.Limit; i++ {
		accountID := fmt.Sprintf("publishers#uuid:%s", uuid.NewV4().String())
		paid := decimal.NewFromFloat(
			float64(rand.Intn(100)),
		).Div(
			decimal.NewFromFloat(10),
		)
		channel := uuid.NewV4().String()

		if untilDate.Before(targetTime.Add(time.Duration(i) * time.Hour * 24)) {
			break
		}
		rows = append(rows, AccountSettlementEarnings{channel, paid, accountID})
		// append sql result rows
		getRows = getRows.AddRow(
			channel,
			paid,
			accountID,
		)
	}
	mock.ExpectQuery(`
select
	channel,
	(.+) as paid,
	account_id
from account_transactions
where (.+)
group by (.+)
order by paid (.+)
limit (.+)`).
		WithArgs(args...).
		WillReturnRows(getRows)
	return rows
}

func SetupMockGetPending(
	mock sqlmock.Sqlmock,
	accountIDs []string,
) []Votes {
	getRows := sqlmock.NewRows(
		[]string{"channel", "balance"},
	)
	rows := []Votes{}
	for _, channel := range accountIDs {
		balance := RandomDecimal()
		rows = append(rows, Votes{channel, balance})
		// append sql result rows
		getRows = getRows.AddRow(
			channel,
			balance,
		)
	}
	mock.ExpectQuery(`
SELECT
	V.channel,
	(.+) as balance
FROM votes V
INNER JOIN surveyor_groups S
ON V.surveyor_id = S.id
WHERE
	V.channel = (.+)
	AND NOT V.transacted
	AND NOT V.excluded
GROUP BY channel`).
		WithArgs(fmt.Sprintf("{%s}", strings.Join(accountIDs, ","))).
		WillReturnRows(getRows)
	return rows
}

func SetupMockGetBalances(
	mock sqlmock.Sqlmock,
	accountIDs []string,
) []Balance {
	getRows := sqlmock.NewRows(
		[]string{"account_id", "account_type", "balance"},
	)
	rows := []Balance{}
	for _, accountID := range accountIDs {
		balance := RandomDecimal()
		accountType := uuid.NewV4().String()
		rows = append(rows, Balance{accountID, accountType, balance})
		// append sql result rows
		getRows = getRows.AddRow(
			accountID,
			accountType,
			balance,
		)
	}
	mock.ExpectQuery(`
	SELECT
		account_transactions.account_type as account_type,
		account_transactions.account_id as account_id,
		(.+) as balance
	FROM account_transactions
	WHERE account_id = (.+)
	GROUP BY (.+)`).
		WithArgs(fmt.Sprintf("{%s}", strings.Join(accountIDs, ","))).
		WillReturnRows(getRows)
	return rows
}

func SetupMockGetTransactions(
	mock sqlmock.Sqlmock,
	accountID string,
	txTypes ...string,
) []Transaction {
	getRows := sqlmock.NewRows(
		[]string{
			"created_at",
			"description",
			"channel",
			"amount",
			"from_account",
			"to_account",
			"to_account_type",
			"settlement_currency",
			"settlement_amount",
			"transaction_type",
		},
	)
	rows := []Transaction{}
	args := []driver.Value{accountID}
	var txTypesHash *map[string]bool
	if len(txTypes) > 0 {
		args = append(args, joinStringList(txTypes))
		txTypesHash := &map[string]bool{}
		for _, txType := range txTypes {
			(*txTypesHash)[txType] = true
		}
	}
	channels := CreateIDs(3)
	providerID := uuid.NewV4().String()
	for _, channel := range channels {
		rows = append(rows, ContributeTransaction(channel))
		rows = append(rows, ReferralTransaction(accountID, channel))
	}
	for i := range channels {
		targetIndex := decimal.NewFromFloat(
			float64(rand.Intn(2)),
		).Mul(decimal.NewFromFloat(float64(i)))
		target := rows[targetIndex.IntPart()]
		rows = append(rows, SettlementTransaction(
			*target.ToAccount, // fromAccount
			target.Channel,    // channel
			providerID,        // toAccount
			target.TransactionType,
		))
	}
	for _, tx := range rows {
		if txTypesHash != nil && !(*txTypesHash)[tx.TransactionType] {
			continue
		}
		getRows = getRows.AddRow(
			tx.CreatedAt,
			tx.Description,
			tx.Channel,
			tx.Amount,
			tx.FromAccount,
			tx.ToAccount,
			tx.ToAccountType,
			tx.SettlementCurrency,
			tx.SettlementAmount,
			tx.TransactionType,
		)
	}
	mock.ExpectQuery(`
SELECT
	created_at,
	description,
	channel,
	amount,
	from_account,
	to_account,
	to_account_type,
	settlement_currency,
	settlement_amount,
	transaction_type
FROM transactions
WHERE (.+)
ORDER BY created_at`).
		WithArgs(args...).
		WillReturnRows(getRows)
	return rows
}

func SettlementTransaction(fromAccount, channel, toAccountID, transactionType string) Transaction {
	transactionType = transactionType + "_settlement"
	provider := "uphold"
	return Transaction{
		Channel:         channel,
		CreatedAt:       time.Now(),
		Description:     uuid.NewV4().String(),
		FromAccount:     fromAccount,
		ToAccount:       &toAccountID,
		ToAccountType:   &provider,
		Amount:          RandomDecimal(),
		TransactionType: transactionType,
	}
}
func ReferralTransaction(accountID, channel string) Transaction {
	toAccountType := "type"
	return Transaction{
		Channel:         channel,
		CreatedAt:       time.Now(),
		Description:     uuid.NewV4().String(),
		FromAccount:     uuid.NewV4().String(),
		ToAccount:       &accountID,
		ToAccountType:   &toAccountType,
		Amount:          RandomDecimal(),
		TransactionType: "referral",
	}
}
func ContributeTransaction(account string) Transaction {
	toAccountType := "type"
	return Transaction{
		Channel:         uuid.NewV4().String(),
		CreatedAt:       time.Now(),
		Description:     uuid.NewV4().String(),
		FromAccount:     uuid.NewV4().String(),
		ToAccount:       &account,
		ToAccountType:   &toAccountType,
		Amount:          RandomDecimal(),
		TransactionType: "contribution",
	}
}

func RandomDecimal() decimal.Decimal {
	return decimal.NewFromFloat(
		float64(rand.Intn(100)),
	).Div(
		decimal.NewFromFloat(10),
	)
}

func CreateIDs(count int) []string {
	list := []string{}
	for i := 0; i < count; i++ {
		list = append(list, uuid.NewV4().String())
	}
	return list
}

func MustMarshal(
	assertions *require.Assertions,
	structure interface{},
) string {
	marshalled, err := json.Marshal(structure)
	assertions.NoError(err)
	return string(marshalled)
}

func (suite *DatastoreMockTestSuite) TestGetAccountEarnings() {
	options := AccountEarningsOptions{
		Limit:     5,
		Ascending: true,
		Type:      "contributions",
	}
	expected := SetupMockGetAccountEarnings(suite.mock, options)
	actual := suite.GetAccountEarnings(
		options,
	)

	suite.Require().JSONEq(
		MustMarshal(suite.Require(), expected),
		MustMarshal(suite.Require(), actual),
	)
}

func (suite *DatastoreMockTestSuite) GetAccountEarnings(
	options AccountEarningsOptions,
) *[]AccountEarnings {
	earnings, err := suite.db.GetAccountEarnings(
		suite.ctx,
		options,
	)
	suite.Require().NoError(err)
	suite.Require().Len(*earnings, options.Limit)
	return earnings
}
func (suite *DatastoreMockTestSuite) TestGetAccountSettlementEarnings() {
	options := AccountSettlementEarningsOptions{
		Limit:     5,
		Ascending: true,
		Type:      "contributions",
	}
	expectSettlementEarnings := SetupMockGetAccountSettlementEarnings(suite.mock, options)
	actualSettlementEarnings := suite.GetAccountSettlementEarnings(options)
	suite.Require().JSONEq(
		MustMarshal(suite.Require(), expectSettlementEarnings),
		MustMarshal(suite.Require(), actualSettlementEarnings),
	)
}

func (suite *DatastoreMockTestSuite) GetAccountSettlementEarnings(
	options AccountSettlementEarningsOptions,
) *[]AccountSettlementEarnings {
	earnings, err := suite.db.GetAccountSettlementEarnings(
		suite.ctx,
		options,
	)
	suite.Require().NoError(err)
	suite.Require().Len(*earnings, options.Limit)
	return earnings
}

func (suite *DatastoreMockTestSuite) TestGetBalances() {
	accountIDs := CreateIDs(3)

	expectedBalances := SetupMockGetBalances(
		suite.mock,
		accountIDs,
	)
	actualBalances := suite.GetBalances(accountIDs)
	suite.Require().JSONEq(
		MustMarshal(suite.Require(), expectedBalances),
		MustMarshal(suite.Require(), actualBalances),
	)
}

func (suite *DatastoreMockTestSuite) GetBalances(accountIDs []string) *[]Balance {
	balances, err := suite.db.GetBalances(
		suite.ctx,
		accountIDs,
	)
	suite.Require().NoError(err)
	suite.Require().Len(*balances, len(accountIDs))
	return balances
}

func (suite *DatastoreMockTestSuite) TestGetPending() {
	accountIDs := CreateIDs(3)

	expectedVotes := SetupMockGetPending(
		suite.mock,
		accountIDs,
	)
	actualVotes := suite.GetPending(accountIDs)
	suite.Require().JSONEq(
		MustMarshal(suite.Require(), expectedVotes),
		MustMarshal(suite.Require(), actualVotes),
	)
}

func (suite *DatastoreMockTestSuite) GetPending(accountIDs []string) *[]Votes {
	votes, err := suite.db.GetPending(
		suite.ctx,
		accountIDs,
	)
	suite.Require().NoError(err)
	suite.Require().Len(*votes, len(accountIDs))
	return votes
}
func (suite *DatastoreMockTestSuite) TestGetTransactions() {
	accountID := CreateIDs(1)[0]

	expectedTransactions := SetupMockGetTransactions(
		suite.mock,
		accountID,
	)
	actualTransaction := suite.GetTransactions(
		len(expectedTransactions),
		accountID,
		nil,
	)
	suite.Require().JSONEq(
		MustMarshal(suite.Require(), expectedTransactions),
		MustMarshal(suite.Require(), actualTransaction),
	)
}

func (suite *DatastoreMockTestSuite) GetTransactions(
	count int,
	accountID string,
	txTypes []string,
) *[]Transaction {
	transactions, err := suite.db.GetTransactions(
		suite.ctx,
		accountID,
		txTypes,
	)
	suite.Require().NoError(err)
	suite.Require().Len(*transactions, count)
	return transactions
}
