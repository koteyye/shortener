package deleter

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/koteyye/shortener/internal/app/models"
	mock_storage "github.com/koteyye/shortener/internal/app/storage/mocks"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestDeleter_StartWorker(t *testing.T) {
	t.Run("send_batch_on_os_signal", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		repo := mock_storage.NewMockURLStorage(c)

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

		repo.EXPECT().GetURLByUser(gomock.Any(), gomock.Any()).Return(testURLList, error(nil))
		repo.EXPECT().DeleteURLByUser(gomock.Any(), url)

		delCh := make(chan DeleteURL)
		d := &Deleter{
			storage:  repo,
			ticker:   time.NewTicker(10 * time.Second),
			delURLch: delCh,
		}

		userID, err := uuid.NewRandom()
		assert.NoError(t, err)
		testDel := DeleteURL{URL: url, UserID: userID.String()}

		go func() {
			d.delURLch <- testDel
		}()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		d.StartWorker(ctx)
	})
	// t.Run("send_batch_on_timeout", func(t *testing.T) {
	// 	c := gomock.NewController(t)
	// 	defer c.Finish()

	// 	repo := mock_storage.NewMockURLStorage(c)

	// 	url := make([]string, 20)
	// 	for i := range url {
	// 		url[i] = randSeq(10)
	// 	}

	// 	testURLList := make([]*models.URLList, 20)
	// 	for i := range testURLList {
	// 		testURLList[i] = &models.URLList{
	// 			Number:   i,
	// 			URL:      "http://localhost:8080/" + randSeq(10),
	// 			ShortURL: url[i],
	// 		}
	// 	}

	// 	repo.EXPECT().GetURLByUser(gomock.Any(), gomock.Any()).Return(testURLList, error(nil))
	// 	repo.EXPECT().DeleteURLByUser(gomock.Any(), url)

	// 	delCh := make(chan DeleteURL)
	// 	d := &Deleter{
	// 		storage:  repo,
	// 		ticker:   time.NewTicker(2 * time.Second),
	// 		delURLch: delCh,
	// 	}

	// 	userID, err := uuid.NewRandom()
	// 	assert.NoError(t, err)
	// 	testDel := DeleteURL{URL: url, UserID: userID.String()}

	// 	go func() {
	// 		d.delURLch <- testDel
	// 	}()

	// 	d.StartWorker(context.Background())
	// })
	// t.Run("send_batch_on_timeout", func(t *testing.T) {
	// 	c := gomock.NewController(t)
	// 	defer c.Finish()

	// 	repo := mock_storage.NewMockURLStorage(c)

	// 	url := make([]string, 20)
	// 	for i := range url {
	// 		url[i] = randSeq(10)
	// 	}

	// 	testURLList := make([]*models.URLList, 20)
	// 	for i := range testURLList {
	// 		testURLList[i] = &models.URLList{
	// 			Number:   i,
	// 			URL:      "http://localhost:8080/" + randSeq(10),
	// 			ShortURL: url[i],
	// 		}
	// 	}

	// 	repo.EXPECT().GetURLByUser(gomock.Any(), gomock.Any()).Return(testURLList, error(nil))
	// 	repo.EXPECT().DeleteURLByUser(gomock.Any(), url)

	// 	delCh := make(chan DeleteURL)
	// 	d := &Deleter{
	// 		storage:  repo,
	// 		ticker:   time.NewTicker(2 * time.Second),
	// 		delURLch: delCh,
	// 	}

	// 	userID, err := uuid.NewRandom()
	// 	assert.NoError(t, err)
	// 	testDel := DeleteURL{URL: url, UserID: userID.String()}

	// 	go func() {
	// 		d.delURLch <- testDel
	// 	}()

	// 	d.StartWorker(context.Background())
	// })
}

// func TestDeleter_InitDeleter(t *testing.T) {
// 	t.Run("success", func(t *testing.T) {
// 		testRepo := storage.URLHandler{}
// 		deleter := InitDeleter(testRepo, &zap.SugaredLogger{})

// 		wantDeleter := &Deleter{storage: testRepo, logger: &zap.SugaredLogger{}, mutex: sync.Mutex{}, test: &unitTest{isTest: false}}

// 		assert.Equal(t, wantDeleter.storage, deleter.storage)
// 		assert.Equal(t, wantDeleter.logger, deleter.logger)
// 		assert.Equal(t, wantDeleter.test, deleter.test)
// 	})
// }

// func testInitDeleter(t *testing.T) (*Deleter, *mock_storage.MockURLStorage) {
// 	c := gomock.NewController(t)
// 	defer c.Finish()

// 	logger, err := zap.NewDevelopment()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer logger.Sync()
// 	log := *logger.Sugar()

// 	repo := mock_storage.NewMockURLStorage(c)
// 	deleter := &Deleter{
// 		URL:     make(chan string, 50),
// 		ticker:  time.NewTicker(2 * time.Second),
// 		storage: repo,
// 		logger:  &log,
// 		mutex:   sync.Mutex{},
// 		test:    &unitTest{isTest: true, msg: make(chan string), mutex: sync.Mutex{}},
// 	}

// 	return deleter, repo
// }

// func TestStartWorker(t *testing.T) {

// 	t.Run("testWorker", func(t *testing.T) {
// 		d, s := testInitDeleter(t)

// 		g, gCtx := errgroup.WithContext(context.Background())
// 		g.Go(func() error {
// 			d.StartWorker(gCtx)
// 			return nil
// 		})
// 		// t.Run("max content", func(t *testing.T) {

// 		// 	url := make([]string, 50)
// 		// 	for i := range url {
// 		// 		url[i] = randSeq(10)
// 		// 	}

// 		// 	testURLList := make([]*models.URLList, 50)
// 		// 	for i := range testURLList {
// 		// 		testURLList[i] = &models.URLList{
// 		// 			Number:   i,
// 		// 			URL:      "http://localhost:8080/" + randSeq(10),
// 		// 			ShortURL: url[i],
// 		// 		}
// 		// 	}

// 		// 	userID, err := uuid.NewRandom()
// 		// 	assert.NoError(t, err)

// 		// 	s.EXPECT().GetURLByUser(gomock.Any(), gomock.Any()).Return(testURLList, error(nil))
// 		// 	s.EXPECT().DeleteURLByUser(gomock.Any(), gomock.Any()).Return(error(nil))

// 		// 	d.Receive(context.Background(), url, userID.String())

// 		// 	d.test.mutex.Lock()
// 		// 	assert.Equal(t, maxItemMsg, <-d.test.msg)
// 		// 	d.test.mutex.Unlock()
// 		// })
// 		t.Run("timer", func(t *testing.T) {
// 			url := make([]string, 40)
// 			for i := range url {
// 				url[i] = randSeq(10)
// 			}

// 			testURLList := make([]*models.URLList, 40)
// 			for i := range testURLList {
// 				testURLList[i] = &models.URLList{
// 					Number:   i,
// 					URL:      "http://localhost:8080/" + randSeq(10),
// 					ShortURL: url[i],
// 				}
// 			}

// 			userID, err := uuid.NewRandom()
// 			assert.NoError(t, err)

// 			s.EXPECT().GetURLByUser(gomock.Any(), gomock.Any()).Return(testURLList, error(nil))
// 			s.EXPECT().DeleteURLByUser(gomock.Any(), gomock.Any()).Return(error(nil))

// 			d.Receive(url, userID.String())

// 			d.test.mutex.Lock()
// 			assert.Equal(t, timeoutMsg, <-d.test.msg)
// 			d.test.mutex.Unlock()
// 		},
// 		)
// 	})
// }

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
