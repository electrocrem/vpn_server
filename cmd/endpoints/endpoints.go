package endpoints

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/electrocrem/vpn_server/cmd/oss"
)

func Page(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}
func GetProfile(w http.ResponseWriter, r *http.Request) {
	randomName := strconv.Itoa(rand.Intn(1000))
	fmt.Printf("%v", randomName)
	oss.LaunchScript("/bin/sh", "./generate_profile.sh", randomName)
	filePath := "profiles/" + randomName + ".ovpn"
	fmt.Printf("\n%v\n", filePath)
	DonwnloadProfile(w, r, filePath)
	os.Remove(filePath)

}

func DonwnloadProfile(w http.ResponseWriter, r *http.Request, filePath string) (err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close() // Close the file after function return
	// Reading header info from the opened file, this will be used for response header "Content-Type"
	fileHeader := make([]byte, 512)
	_, err = file.Read(fileHeader) // File offset is now len(fileHeader)
	if err != nil {
		return err
	}
	// Get file info which we will use for the response headers "Content-Disposition" and "Content-Length"
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	// Set default headers
	// attachment is required to tell some (older) browsers who follow an href to download the file
	// instead of showing/printing the content to the screen.
	// For example, if you click a link to an image, the browser will pop up the download dialog
	// box compared to drawing the image in the browser tab.
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s""`, fileInfo.Name()))
	// A must for every request that has a body (see RFC 2616 section 7.2.1)
	w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	// Tell the client we accept ranges, this gives clients the option to pause the transfer
	// and pick up later where they left up. Or download managers to establish multiple connections
	w.Header().Set("Accept-Ranges", "bytes")
	// Check if the client requests a range from the file (see RFC 7233 section 4.2)
	requestRange := r.Header.Get("range")
	if requestRange == "" {
		// No range is defined, tell the client the incoming length of data, the size of the open file
		w.Header().Set("Content-Length", strconv.Itoa(int(fileInfo.Size())))
		// Since we read 512 bytes for 'fileHeader' earlier, we set the reader offset back
		// to 0 starting from the beginning of the the file (the 0 in the second argument)
		file.Seek(0, 0)
		// Stream the file to the client
		io.Copy(w, file)
		return nil
	}
	// Client requests a part of the file
	// Decode the request header to integers we can use for offset
	requestRange = requestRange[6:] // Strip the "bytes=", left over is now "begin-end"
	splitRange := strings.Split(requestRange, "-")
	if len(splitRange) != 2 {
		return fmt.Errorf("invalid values for header 'Range'")
	}
	begin, err := strconv.ParseInt(splitRange[0], 10, 64)
	if err != nil {
		return err
	}
	end, err := strconv.ParseInt(splitRange[1], 10, 64)
	if err != nil {
		return err
	}
	if begin > fileInfo.Size() || end > fileInfo.Size() {
		return fmt.Errorf("range out of bounds for file")
	}
	if begin >= end {
		return fmt.Errorf("range begin cannot be bigger than range end")
	}
	// Tell the amount bytes the client will receive
	w.Header().Set("Content-Length", strconv.FormatInt(end-begin+1, 10))
	// Confirm the range values to the client, and the total size of the file
	// 'Content-Range' : 'bytes begin-end/totalFileSize'
	w.Header().Set("Content-Range",
		fmt.Sprintf("bytes %d-%d/%d", begin, end, fileInfo.Size()))
	// Response http status code 206
	w.WriteHeader(http.StatusPartialContent)
	// Set the file offset to the requested beginning
	file.Seek(begin, 0)
	// Send the (end-begin) amount of bytes to the client
	io.CopyN(w, file, end-begin)
	return nil
}
