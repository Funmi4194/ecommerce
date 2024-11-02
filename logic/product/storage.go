package product

import (
	"database/sql"
	"errors"
	"fmt"

	"mime/multipart"

	"github.com/funmi4194/bifrost"
	"github.com/funmi4194/ecommerce/enum"
	"github.com/funmi4194/ecommerce/helper"
	"github.com/funmi4194/ecommerce/primer"
	userRepository "github.com/funmi4194/ecommerce/repository/user"
	"github.com/funmi4194/ecommerce/storage"
	"github.com/funmi4194/ecommerce/types"
	"github.com/opensaucerer/barf"
)

// Store allows the upload of image file
func Store(fs []*multipart.FileHeader, userId string) ([]types.Object, error) {

	user := userRepository.User{
		ID: userId,
	}

	// find user by Id
	err := user.FByKeyVal("id", user.ID, true)
	if err != nil {
		barf.Logger().Errorf(`[product.Store] [user.FByKeyVal("id", user.ID, true)] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, errors.New("looks like your account no longer exists. please contact support")
		}
		return nil, errors.New("we're having issues uploading product. please try again later")
	}

	if user.Role != enum.Admin {
		return nil, errors.New("you do not have the permission to this feature")
	}

	if len(fs) == 0 {
		return nil, errors.New("no files to upload")
	}

	acl := bifrost.ACLPublicRead

	files := bifrost.MultiFile{}

	for _, file := range fs {

		if helper.DetermineFileFormat(file.Filename) == "" {
			barf.Logger().Errorf(`[product.Store] [DetermineFileFormat(file.Filename)] %s`, "file format not supported")
			return nil, fmt.Errorf("unsupported file format detected on file %s", file.Filename)
		}

		f, err := file.Open()
		if err != nil {
			barf.Logger().Errorf(`[product.Store] [f, err := file.Open()] %s`, err.Error())
			return nil, err
		}

		files.Files = append(files.Files, bifrost.File{
			Handle:   f,
			Filename: helper.GenerateFilename(file.Filename),
			Options: map[string]interface{}{
				bifrost.OptACL: acl,
				bifrost.OptMetadata: map[string]string{
					"originalName": file.Filename,
				},
			},
		})
	}

	bridge, err := storage.NewGCSRainbowBridge(primer.ENV.OriginalBucket)
	if err != nil {
		barf.Logger().Fatalf(`[product.Store] [storage.NewGCSRainbowBridge(primer.ENV.OriginalBucket)] %s`, err.Error())
		return nil, errors.New("we're having some trouble uploading your files, please try again later")
	}

	objects, err := bridge.UploadMultiFile(files)
	if err != nil {
		barf.Logger().Errorf(`[product.Store] [bridge.UploadMultiFile(files)] %s`, err.Error())
		return nil, err
	}

	var objs = []types.Object{}

	for i, object := range objects {

		objs = append(objs, types.Object{
			Name:          object.Name,
			OriginalName:  fs[i].Filename,
			RemoteAddress: object.Preview,
			Size:          object.Size,
			FileFormat:    helper.DetermineFileFormat(object.Name),
			Error:         object.Error,
		})
	}

	return objs, nil
}
