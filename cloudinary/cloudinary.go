package cloudinary

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Manan-Rastogi/chezzTubeBackend-GO/configs"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/utils"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var CLOUDINARY *cloudinary.Cloudinary

func init() {
	CLOUDINARY = credentials()
}

func credentials() *cloudinary.Cloudinary {
	// Add your Cloudinary credentials, set configuration parameter
	// Secure=true to return "https" URLs, and create a context
	//===================
	cld, err := cloudinary.NewFromParams(configs.ENV.CloudinaryCloudName, configs.ENV.CloudinaryApiKey, configs.ENV.CloudinaryApiSecret)
	fmt.Printf("err: %v\n", err)

	cld.Config.URL.Secure = true

	return cld
}

// UploadImage uploads an image to Cloudinary.
// Parameters:
// - cld: a pointer to a Cloudinary instance
// - wg: a pointer to a WaitGroup for synchronization
// - timeout: the duration for the upload operation to complete
// - file: the image file to upload
// - publicId: the public ID for the uploaded image
// Returns the secure URL of the uploaded image.
func UploadImage(cld *cloudinary.Cloudinary, wg *sync.WaitGroup, timeout time.Duration, file interface{}, publicId string) string {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Upload the image.
	// Set the asset's public ID and allow overwriting the asset with new versions
	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       publicId,
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(true)})
	if err != nil {
		utils.Logger.Error(err.Error())
	}

	// Log the delivery URL
	return resp.SecureURL
}

// DestroyImage deletes an image using Cloudinary and sync.WaitGroup
func DestroyImage(cld *cloudinary.Cloudinary, wg *sync.WaitGroup, publicId string) {
	defer wg.Done()

	_, err := cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicId,
	})
	if err != nil {
		log.Printf("Error deleting image: %v\n", err)
	} else {
		log.Printf("Image deleted successfully: %s\n", publicId)
	}
}
