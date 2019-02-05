package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/minio/minio-go"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func SaveFile(file multipart.File, path string) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

func UploadFile(filename string) {
	endpoint := "sfo2.digitaloceanspaces.com"
	accessKeyID := os.Getenv("SPACES_KEY")
	secretAccessKey := os.Getenv("SPACES_SECRET")
	useSSL := true

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
		return
	}

	bucketName := "feather"
	objectName := filename
	filePath := "/tmp/" + filename
	contentType := "image/png"

	userMetadata := map[string]string{
		"x-amz-acl": "public-read",
	}

	// Upload the zip file with FPutObject
	n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType, UserMetadata: userMetadata})
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
}

func RoundTime(input float64) int {
	var result float64

	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}

	// only interested in integer, ignore fractional
	i, _ := math.Modf(result)

	return int(i)
}

func TimeElapsed(inputDate *time.Time, full bool) string {
	var precise [8]string // this an array, not slice
	var text string
	var future bool // our crystal ball

	// get years, months and days
	// and get hours, minutes, seconds
	now := time.Now()
	year2, month2, day2 := now.Date()
	hour2, minute2, second2 := now.Clock()

	year1, month1, day1 := inputDate.Date()
	hour1, minute1, second1 := inputDate.Clock()

	// are we forecasting the future?
	if (year2 - year1) < 0 {
		future = true
	}

	if (month2 - month1) < 0 {
		future = true
	}
	if (day2 - day1) < 0 {
		future = true
	}
	if (hour2 - hour1) < 0 {
		future = true
	}
	if (minute2 - minute1) < 0 {
		future = true
	}
	if (second2 - second1) < 0 {
		future = true
	}

	// convert negative to positive numbers
	year := math.Abs(float64(int(year2 - year1)))
	month := math.Abs(float64(int(month2 - month1)))
	day := math.Abs(float64(int(day2 - day1)))
	hour := math.Abs(float64(int(hour2 - hour1)))
	minute := math.Abs(float64(int(minute2 - minute1)))
	second := math.Abs(float64(int(second2 - second1)))

	week := math.Floor(day / 7)

	// Ouch!, no if-else short hand - see https://golang.org/doc/faq#Does_Go_have_a_ternary_form

	if year > 0 {
		if int(year) == 1 {
			precise[0] = strconv.Itoa(int(year)) + " year"
		} else {
			precise[0] = strconv.Itoa(int(year)) + " years"
		}
	}

	if month > 0 {
		if int(month) == 1 {
			precise[1] = strconv.Itoa(int(month)) + " month"
		} else {
			precise[1] = strconv.Itoa(int(month)) + " months"
		}
	}

	if week > 0 {
		if int(week) == 1 {
			precise[2] = strconv.Itoa(int(week)) + " week"
		} else {
			precise[2] = strconv.Itoa(int(week)) + " weeks"
		}
	}

	if day > 0 {
		if int(day) == 1 {
			precise[3] = strconv.Itoa(int(day)) + " day"
		} else {
			precise[3] = strconv.Itoa(int(day)) + " days"
		}
	}

	if hour > 0 {
		if int(hour) == 1 {
			precise[4] = strconv.Itoa(int(hour)) + " hour"
		} else {
			precise[4] = strconv.Itoa(int(hour)) + " hours"
		}
	}

	if minute > 0 {
		if int(minute) == 1 {
			precise[5] = strconv.Itoa(int(minute)) + " minute"
		} else {
			precise[5] = strconv.Itoa(int(minute)) + " minutes"
		}
	}

	if second > 0 {
		if int(second) == 1 {
			precise[6] = strconv.Itoa(int(second)) + " second"
		} else {
			precise[6] = strconv.Itoa(int(second)) + " seconds"
		}
	}

	for _, v := range precise {
		if v != "" {
			// no comma after second
			if v[len(v)-5:len(v)-1] != "cond" {
				precise[7] += v + ", "
			} else {
				precise[7] += v
			}
		}
	}

	if !future {
		text = " ago."
	} else {
		text = " in future."
	}

	if full {
		return precise[7] + text
	} else {
		// return the first non-empty position
		for k, v := range precise {
			if v != "" {
				return precise[k] + text
			}
		}
	}
	return "invalid date"
}
