package repository

import (
	"arvanch/config"
	"arvanch/db"
	"arvanch/model"
	"testing"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type MessageRepoSuiteTest struct {
	suite.Suite
	repo MessageRepository
	db   *gorm.DB
}

func (suite *MessageRepoSuiteTest) SetupSuite() {
	database := db.WithRetry(db.Create, config.Init().Postgres)
	suite.db = database
	suite.repo = NewMessageRepo(database)
}

// nolint:funlen,gocognit
func (suite *MessageRepoSuiteTest) TestInsertMessage() {
	userID := uuid.New().String()
	accID := uuid.New().String()

	suite.NoError(suite.db.Create(&model.Account{ID: accID}).Error)

	suite.NoError(suite.db.Create(&model.User{AccountID: accID, Name: "user_test", ID: userID}).Error)

	tcs := []struct {
		name        string
		msg         *model.Message
		errExpected bool
	}{
		{
			name: "successful send en",
			msg: &model.Message{
				ID:       uuid.New().String(),
				UserID:   userID,
				Payload:  "payload 1",
				Language: "en",
			},
		},
		{
			name: "successful send fa",
			msg: &model.Message{
				ID:       uuid.New().String(),
				UserID:   userID,
				Payload:  "متن 1",
				Language: "en",
			},
		},
		{
			name: "failed due to nil userID",
			msg: &model.Message{
				ID:       uuid.New().String(),
				Payload:  "payload 2",
				Language: "en",
			},
			errExpected: true,
		},
		{
			name: "failed due to empty language",
			msg: &model.Message{
				ID:      uuid.New().String(),
				Payload: "payload 2",
				UserID:  userID,
			},
			errExpected: true,
		},
		{
			name: "failed due to empty payload",
			msg: &model.Message{
				ID:       uuid.New().String(),
				Language: "en",
				UserID:   userID,
			},
			errExpected: true,
		},
	}

	for i := range tcs {
		tc := tcs[i]

		suite.Run(tc.name, func() {
			err := suite.repo.InsertMessage(tc.msg)

			if tc.errExpected {
				suite.Error(err)
			} else {
				suite.NoError(err)

				var m model.Message
				suite.NoError(suite.db.Where("ID = ? ", tc.msg.ID).Find(&m).Error)

				suite.Equal(m.ID, tc.msg.ID)
				suite.Equal(m.UserID, tc.msg.UserID)
				suite.Equal(m.Language, tc.msg.Language)
				suite.Equal(m.Payload, tc.msg.Payload)
			}

		})
	}
}

func TestSMS(t *testing.T) {
	suite.Run(t, new(MessageRepoSuiteTest))
}
