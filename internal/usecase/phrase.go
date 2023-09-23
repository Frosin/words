package usecase

import (
	"context"
	"fmt"
	"test/internal/entity"
	"test/internal/metrics"
)

func (u *Uc) CreatePhrase(ctx context.Context, userID int64, phrase string) error {
	obj, err := u.repo.GetPhrase(ctx, userID, phrase)
	if err != nil {
		return err
	}

	// if phrase already exist
	if obj != nil {
		return nil
	}

	metaSerialized, err := entity.PhraseMeta{}.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize phrase meta")
	}

	newObj := entity.Phrase{
		Phrase: phrase,
		LangID: 0,
		Epoch:  0,
		UserID: userID,
		Meta:   metaSerialized,
	}

	_, err = u.repo.SavePhrase(ctx, newObj)

	return err
}

func (u *Uc) UpdatePhrase(ctx context.Context, userID int64, phrase string, sentence string) error {
	obj, err := u.repo.GetPhrase(ctx, userID, phrase)
	if err != nil {
		return err
	}

	if obj == nil {
		return fmt.Errorf("phrase for update not found: %s", phrase)
	}

	// send epoch mitrics
	switch obj.Epoch {
	case 1:
		metrics.WordsPhraseEpoch1.Inc()
	case 2:
		metrics.WordsPhraseEpoch2.Inc()
	case 3:
		metrics.WordsPhraseEpoch3.Inc()
	}

	obj.Epoch++

	meta, err := entity.DeserializePhraseMeta(obj.Meta)
	if err != nil {
		return fmt.Errorf("failed to deserialize phrase meta")
	}

	meta.Sentences = append(meta.Sentences, sentence)

	metaSerialized, err := meta.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize phrase meta")
	}

	obj.Meta = metaSerialized

	_, err = u.repo.SavePhrase(ctx, *obj)

	return err
}

func (u *Uc) DeletePhraseSentence(ctx context.Context, userID int64, phrase string, sentence string) error {
	obj, err := u.repo.GetPhrase(ctx, userID, phrase)
	if err != nil {
		return err
	}

	if obj == nil {
		return fmt.Errorf("phrase for update not found: %s", phrase)
	}

	// send epoch mitrics
	switch obj.Epoch {
	case 1:
		metrics.WordsPhraseEpoch1.Inc()
	case 2:
		metrics.WordsPhraseEpoch2.Inc()
	case 3:
		metrics.WordsPhraseEpoch3.Inc()
	}

	meta, err := entity.DeserializePhraseMeta(obj.Meta)
	if err != nil {
		return fmt.Errorf("failed to deserialize phrase meta")
	}

	deleted := false
	for i, ms := range meta.Sentences {
		if ms == sentence {
			// delete this sentence
			if i < len(meta.Sentences)-1 {
				copy(meta.Sentences[i:], meta.Sentences[i+1:])
			}
			meta.Sentences[len(meta.Sentences)-1] = ""
			meta.Sentences = meta.Sentences[:len(meta.Sentences)-1]

			deleted = true
			obj.Epoch--
		}
	}

	if !deleted {
		return fmt.Errorf("sentence %s:%s not found ", phrase, sentence)
	}

	metaSerialized, err := meta.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize phrase meta")
	}

	obj.Meta = metaSerialized

	_, err = u.repo.SavePhrase(ctx, *obj)

	return err
}

func (u *Uc) GetReminderPhrases(ctx context.Context) ([]*entity.Phrase, error) {
	return u.repo.GetReminderPhrases(ctx)
}

func (u *Uc) GetPhraseInfo(ctx context.Context, userID int64, phrase string) (*entity.Phrase, error) {
	obj, err := u.repo.GetPhrase(ctx, userID, phrase)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (u *Uc) DeletePhrase(ctx context.Context, userID int64, phrase string) error {
	obj, err := u.repo.GetPhrase(ctx, userID, phrase)
	if err != nil {
		return err
	}

	if obj == nil {
		return fmt.Errorf("phrase '%s' for user '%d' not found", phrase, userID)
	}

	return u.repo.DeletePhrase(ctx, obj)
}
