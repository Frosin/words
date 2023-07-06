package repository

import (
	"context"
	"os"
	"test/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testDB = "test.db"

func getTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(testDB), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entity.Phrase{}, &entity.UserSettings{}, &entity.Session{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func removeDB() {
	_ = os.Remove(testDB)
}

func Test_SavePhrase_SimpleCase_Success(t *testing.T) {
	db, err := getTestDB()
	assert.NoError(t, err)

	defer removeDB()

	rep := Repo{db}

	ctx := context.Background()
	phrase0 := entity.Phrase{
		Phrase: "test",
		LangID: 1,
		UserID: 123,
		Epoch:  1,
		Meta:   []byte(`{"sentences":["test"]}`),
	}

	id, err := rep.SavePhrase(ctx, phrase0)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), id)

	// check it
	phrase, err := rep.GetPhrase(ctx, 123, "test")
	assert.NoError(t, err)

	assert.Equal(t, "test", phrase.Phrase)
	assert.Equal(t, uint8(1), phrase.Epoch)
	assert.Equal(t, int64(123), phrase.UserID)
	assert.Equal(t, uint8(1), phrase.LangID)
}

func Test_UpdatePhrase_SimpleCase_Success(t *testing.T) {
	db, err := getTestDB()
	assert.NoError(t, err)

	defer removeDB()

	rep := Repo{db}

	ctx := context.Background()

	phrase0 := entity.Phrase{
		Phrase: "test",
		LangID: 1,
		UserID: 123,
		Epoch:  1,
		Meta:   []byte(`{"sentences":["test"]}`),
	}

	phrase := entity.Phrase{
		Phrase: "test",
		LangID: 2,
		UserID: 123,
		Epoch:  2,
		Meta:   []byte(`{"sentences":["test1","test2"]}`),
	}

	id, err := rep.SavePhrase(ctx, phrase0)
	assert.NoError(t, err)

	phrase.ID = id

	// update phrase
	_, err = rep.SavePhrase(ctx, phrase)
	assert.NoError(t, err)

	// check it
	gPhrase, err := rep.GetPhrase(ctx, 123, "test")
	assert.NoError(t, err)

	gMeta, err := entity.DeserializePhraseMeta(gPhrase.Meta)
	assert.NoError(t, err)

	assert.Equal(t, "test", gPhrase.Phrase)
	assert.Equal(t, uint8(2), gPhrase.Epoch)
	assert.Equal(t, int64(123), gPhrase.UserID)
	assert.Equal(t, uint8(2), gPhrase.LangID)
	assert.ElementsMatch(t, gMeta.Sentences, []string{"test1", "test2"})
}

func Test_SavePhrase_Not_Found_Fail(t *testing.T) {
	db, err := getTestDB()
	assert.NoError(t, err)

	defer removeDB()

	rep := Repo{db}

	phrase0 := entity.Phrase{
		Phrase: "test",
		LangID: 1,
		UserID: 123,
		Epoch:  1,
		Meta:   []byte(`{"sentences":["test"]}`),
	}

	phrase := entity.Phrase{
		Phrase: "test2",
		LangID: 1,
		UserID: 456,
		Epoch:  1,
		Meta:   []byte(`{"sentences":["test"]}`),
	}

	ctx := context.Background()
	_, err = rep.SavePhrase(ctx, phrase0)
	assert.NoError(t, err)

	_, err = rep.SavePhrase(ctx, phrase)
	assert.NoError(t, err)

	// check it
	gPhrase, err := rep.GetPhrase(ctx, 456, "test")
	assert.NoError(t, err)
	assert.Nil(t, gPhrase)
}

func Test_GetReminderPhrases_SimpleCase_Success(t *testing.T) {
	db, err := getTestDB()
	assert.NoError(t, err)

	defer removeDB()

	rep := Repo{db}

	// create test data
	err = db.Exec(`insert into phrases (phrase, user_id, epoch, updated_at) values
		('test1', 100, 0, datetime('now', '-2 hours')),
		('test2', 101, 0, datetime('now', '-2 hours')),
		('test3', 100, 1, datetime('now', '-1 day')),
		('test4', 101, 1, datetime('now', '-1 day')),
		('test5', 100, 2, datetime('now', '-14 days')),
		('test6', 101, 2, datetime('now', '-14 days')),
		('test7', 100, 3, datetime('now', '-2 months')),
		('test8', 101, 3, datetime('now', '-2 months')),
		('bad1', 101, 0, datetime('now', '-1 hour')),
		('bad2', 102, 1, datetime('now', '-20 hour')),
		('bad3', 101, 2, datetime('now', '-2 days')),
		('bad4', 102, 3, datetime('now', '-1 month')),
		('bad5', 103, 1, datetime('now', '-2 hours')),
		('bad6', 104, 2, datetime('now', '-1 day')),
		('bad7', 105, 3, datetime('now', '-14 days')),
		('bad8', 106, 4, datetime('now', '-2 months'));
	`).Error
	assert.NoError(t, err)

	ctx := context.Background()

	// check it
	phrases, err := rep.GetReminderPhrases(ctx)
	assert.NoError(t, err)
	assert.Len(t, phrases, 8)

	strs := []string{}
	for _, ph := range phrases {
		strs = append(strs, ph.Phrase)
	}

	assert.ElementsMatch(t, strs, []string{
		"test1", "test2", "test3", "test4", "test5", "test6", "test7", "test8",
	})
}

func Test_SaveSettings_Not_Found_Fail(t *testing.T) {
	db, err := getTestDB()
	assert.NoError(t, err)

	defer removeDB()

	rep := Repo{db}

	ctx := context.Background()

	sets := entity.Settings{
		Langs: []entity.LangSettings{
			{
				LangID:    0,
				PhraseNum: 3,
			},
			{
				LangID:    1,
				PhraseNum: 4,
			},
			{
				LangID:    2,
				PhraseNum: 5,
			},
		},
	}
	serialisedSets, err := sets.Serialize()
	assert.NoError(t, err)

	err = rep.SaveSettings(ctx, 123, serialisedSets)
	assert.NoError(t, err)

	// check it
	gSerSets, err := rep.GetSettings(ctx, 456)
	assert.NoError(t, err)
	assert.Nil(t, gSerSets)
}

func Test_SaveSettings_SimpleCase_Success(t *testing.T) {
	db, err := getTestDB()
	assert.NoError(t, err)

	defer removeDB()

	rep := Repo{db}

	ctx := context.Background()

	sets := entity.Settings{
		Langs: []entity.LangSettings{
			{
				LangID:    0,
				PhraseNum: 3,
			},
			{
				LangID:    1,
				PhraseNum: 4,
			},
			{
				LangID:    2,
				PhraseNum: 5,
			},
		},
	}
	serialisedSets, err := sets.Serialize()
	assert.NoError(t, err)

	err = rep.SaveSettings(ctx, 123, serialisedSets)
	assert.NoError(t, err)

	// check it
	gSerSets, err := rep.GetSettings(ctx, 123)
	assert.NoError(t, err)

	gSets, err := entity.DeserializeSettings(gSerSets)
	assert.NoError(t, err)

	assert.Len(t, sets.Langs, 3)

	for i := range sets.Langs {
		assert.Equal(t, sets.Langs[i].LangID, gSets.Langs[i].LangID)
		assert.Equal(t, sets.Langs[i].PhraseNum, gSets.Langs[i].PhraseNum)
	}
}

func Test_SaveSession_SimpleCase_Success(t *testing.T) {
	db, err := getTestDB()
	assert.NoError(t, err)

	defer removeDB()

	rep := Repo{db}

	ctx := context.Background()
	session := entity.Session{
		UserID:      123,
		CurrentPage: "settings",
		LastMsgID:   10000,
		ChatID:      12,
	}

	err = rep.SaveSession(ctx, session)
	assert.NoError(t, err)

	// check it
	gSess, err := rep.GetSession(ctx, 123)
	assert.NoError(t, err)

	assert.Equal(t, session.UserID, gSess.UserID)
	assert.Equal(t, session.CurrentPage, gSess.CurrentPage)
	assert.Equal(t, session.LastMsgID, gSess.LastMsgID)
	assert.Equal(t, session.ChatID, gSess.ChatID)
}

func Test_SaveSession_Not_Found_Fail(t *testing.T) {
	db, err := getTestDB()
	assert.NoError(t, err)

	defer removeDB()

	rep := Repo{db}

	ctx := context.Background()
	session := entity.Session{
		UserID:      123,
		CurrentPage: "settings",
		LastMsgID:   10000,
		ChatID:      12,
		Data:        "",
	}

	err = rep.SaveSession(ctx, session)
	assert.NoError(t, err)

	// check it
	gSess, err := rep.GetSession(ctx, 456)
	assert.NoError(t, err)
	assert.Nil(t, gSess)
}

func Test_GetSessions_SimpleCase_Success(t *testing.T) {
	db, err := getTestDB()
	assert.NoError(t, err)

	defer removeDB()

	rep := Repo{db}

	ctx := context.Background()
	sessions := []entity.Session{
		{
			UserID:      123,
			CurrentPage: "settings",
			LastMsgID:   10000,
			ChatID:      12,
			Data:        "test",
		},
		{
			UserID:      456,
			CurrentPage: "settings",
			LastMsgID:   10001,
			ChatID:      12,
			Data:        "test",
		},
		{
			UserID:      789,
			CurrentPage: "settings",
			LastMsgID:   10002,
			ChatID:      12,
			Data:        "test",
		},
	}

	for _, sess := range sessions {
		err = rep.SaveSession(ctx, sess)
		assert.NoError(t, err)
	}

	// check it
	gSessions, err := rep.GetSessions(ctx, []int64{789, 123})
	assert.NoError(t, err)

	assert.Len(t, gSessions, 2)

	assert.ElementsMatch(t, []int{789, 123}, []int{gSessions[0].UserID, gSessions[1].UserID})
}
