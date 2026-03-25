package repository

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRoomRepositoryReturnsMessageCounts(t *testing.T) {
	db := openRoomTestDB(t)
	repo := NewRoomRepository(db)

	ownerID := uuid.New()
	room := domain.Room{
		UniqueString: "room-1",
		Name:         "Room 1",
		OwnerID:      ownerID,
	}
	require.NoError(t, db.Create(&room).Error)

	require.NoError(t, db.Create(&domain.Message{RoomId: room.ID, FromUnique: "sender-1", Text: "first"}).Error)
	require.NoError(t, db.Create(&domain.Message{RoomId: room.ID, FromUnique: "sender-2", Text: "second"}).Error)

	single, err := repo.GetByUniqueString(room.UniqueString)
	require.NoError(t, err)
	require.Equal(t, 2, single.MessagesCnt)

	listed, err := repo.GetByOwnerId(ownerID.String(), options.BaseFetchOptions{})
	require.NoError(t, err)
	require.Len(t, listed.Data, 1)
	require.Equal(t, 2, listed.Data[0].MessagesCnt)
}

func openRoomTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "rooms.sqlite")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, sqlDB.Close())
	})

	require.NoError(t, db.AutoMigrate(&domain.User{}, &domain.Room{}, &domain.Message{}))

	return db
}
