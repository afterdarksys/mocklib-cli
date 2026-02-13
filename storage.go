package main

import "fmt"

// Storage command implementations

func storageCreateBucket(bucketName string, provider ...string) error {
	prov := "s3"
	if len(provider) > 0 {
		prov = provider[0]
	}

	reqBody := map[string]interface{}{
		"Action":     "CreateBucket",
		"BucketName": bucketName,
		"Provider":   prov,
		"Region":     "us-east-1",
	}

	resp, err := makeRequest("POST", "/storage/bucket", reqBody)
	if err != nil {
		return err
	}

	// Print bucket name for easy capture
	fmt.Println(resp["BucketName"])
	return nil
}

func storageDeleteBucket(bucketName string) error {
	reqBody := map[string]interface{}{
		"Action":     "DeleteBucket",
		"BucketName": bucketName,
	}

	_, err := makeRequest("POST", "/storage/bucket", reqBody)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted bucket: %s\n", bucketName)
	return nil
}
