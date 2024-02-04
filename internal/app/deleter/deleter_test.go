package deleter

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/storage"
	mock_storage "github.com/koteyye/shortener/internal/app/storage/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestDeleter_InitDeleter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testRepo := storage.URLHandler{}
		deleter := InitDeleter(testRepo, &zap.SugaredLogger{})

		wantDeleter := &Deleter{URL: make(chan string, 50), storage: testRepo, logger: &zap.SugaredLogger{}, mutex: sync.Mutex{}, ticker: time.NewTicker(10 * time.Second), test: &unitTest{isTest: false}}

		assert.Equal(t, wantDeleter.storage, deleter.storage)
		assert.Equal(t, wantDeleter.logger, deleter.logger)
		assert.Equal(t, wantDeleter.test, deleter.test)
	})
}

func testInitDeleter(t *testing.T) (*Deleter, *mock_storage.MockURLStorage) {
	g, gCtx := errgroup.WithContext(context.Background())

	c := gomock.NewController(t)
	defer c.Finish()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	log := *logger.Sugar()

	repo := mock_storage.NewMockURLStorage(c)
	deleter := &Deleter{
		URL:     make(chan string, 50),
		ticker:  time.NewTicker(2 * time.Second),
		storage: repo,
		logger:  &log,
		mutex:   sync.Mutex{},
		test:    &unitTest{isTest: true, msg: make(chan string)},
	}

	g.Go(func() error {
		deleter.StartWorker(gCtx)
		return nil
	})

	return deleter, repo
}

func TestStartWorker(t *testing.T) {

	t.Run("testWorker", func(t *testing.T) {
		
		t.Run("timer", func(t *testing.T) {
			d, s := testInitDeleter(t)


			url := make([]string, 40)
			for i := range url {
				url[i] = randSeq(10)
			}

			testURLList := make([]*models.URLList, 40)
			for i := range testURLList {
				testURLList[i] = &models.URLList{
					Number:   i,
					URL:      "http://localhost:8080/" + randSeq(10),
					ShortURL: url[i],
				}
			}

			userID, err := uuid.NewRandom()
			assert.NoError(t, err)

			s.EXPECT().GetURLByUser(gomock.Any(), gomock.Any()).Return(testURLList, error(nil))
			s.EXPECT().DeleteURLByUser(gomock.Any(), gomock.Any()).Return(error(nil))

			d.Receive(context.Background(), url, userID.String())

			assert.Equal(t, timeoutMsg, <-d.test.msg)
		},
		)
	})
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

