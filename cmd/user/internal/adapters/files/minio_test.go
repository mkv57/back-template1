//go:build integration

package files_test

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ZergsLaw/back-template/cmd/user/internal/app"
)

const (
	avatarFilePath = `testdata/test.jpg`
)

func TestClient_Smoke(t *testing.T) {
	t.Parallel()

	ctx, fileStore, assert := start(t)

	fImg, err := os.Open(avatarFilePath)
	assert.NoError(err)
	imgBuf, err := io.ReadAll(fImg)
	assert.NoError(err)
	_, err = fImg.Seek(0, io.SeekStart)
	assert.NoError(err)
	imgID := uuid.Must(uuid.NewV4())
	st, err := fImg.Stat()
	assert.NoError(err)

	id, err := fileStore.UploadFile(ctx, app.Avatar{
		ID:             imgID,
		Name:           fImg.Name(),
		ContentType:    getContentType(t, assert, fImg),
		Size:           st.Size(),
		ModTime:        st.ModTime(),
		ReadSeekCloser: fImg,
	})
	assert.NoError(err)
	assert.NotEmpty(id)

	fImgFromStore, err := fileStore.DownloadFile(ctx, id)
	assert.NoError(err)
	t.Cleanup(func() {
		assert.NoError(fImgFromStore.Close())
	})

	imgFromStoreBuf, err := io.ReadAll(fImgFromStore)
	assert.NoError(err)
	assert.Equal(imgBuf, imgFromStoreBuf)
}

func getContentType(t *testing.T, assert *require.Assertions, r io.ReadSeeker) string {
	t.Helper()

	var buf [512]byte
	_, err := r.Read(buf[:])
	assert.NoError(err)
	_, err = r.Seek(0, io.SeekStart)
	assert.NoError(err)

	return http.DetectContentType(buf[:])
}
