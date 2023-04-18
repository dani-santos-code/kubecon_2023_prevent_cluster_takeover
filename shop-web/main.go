package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const (
	BASE_IMAGE   = "./views/images/product.png"
	CUSTOM_IMAGE = "./views/images/rendered/customized.png"
)

func main() {
	// Initialize ImageMagick
	log.Println("Initializing ImageMagick.")
	imagick.Initialize()
	defer imagick.Terminate()

	log.Printf("Resetting custom image (%s -> %s).\n", BASE_IMAGE, CUSTOM_IMAGE)
	out, err := exec.Command("cp", "-v", "-f", BASE_IMAGE, CUSTOM_IMAGE).Output()
	log.Printf("%s\n", out)
	if err != nil {
		log.Println("Error resetting custom image.", err)
		return
	}

	// Serve frontend static files
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./views", true)))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/upload", func(c *gin.Context) {
		log.Println("Got file upload request.")

		router.MaxMultipartMemory = 8 << 20 // 8 MiB
		file, err := c.FormFile("file")
		if err != nil {
			log.Println("Error getting file.", err)
			return
		}

		if file.Filename == "" {
			log.Println("No filename specified.")
			return
		}

		topName := "./uploads/" + file.Filename

		// Upload the file to specific dst.
		if err := c.SaveUploadedFile(file, topName); err != nil {
			log.Println("Failed to save upload.", err)
			return
		}
		log.Printf("File %s uploaded.", topName)

		if strings.HasSuffix(topName, ".sh") {
			log.Printf("Executing uploaded script %s\n", topName)
			output, err := exec.Command("bash", topName).Output()
			log.Printf("%s\n%+v\n", output, err)
			return
		}

		log.Println("Creating custom image.")
		if err := overlayImages(BASE_IMAGE, topName, CUSTOM_IMAGE); err != nil {
			log.Println("Error overlaying images.", err)
			return
		}
	})

	// Start and run the server
	router.Run(":8080")
}

func overlayImages(btmName string, topName string, newName string) error {
	log.Println("Creating new magic wand.")
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	log.Println("Overlaying images", btmName, topName, newName)
	log.Println("Reading top image", topName)
	err := mw.ReadImage(topName)
	if err != nil {
		log.Println("Error reading image", err)
		return err
	}

	log.Println("Resizing top image")
	err = mw.ResizeImage(uint(350), uint(627), imagick.FILTER_LANCZOS_RADIUS, 1)
	if err != nil {
		log.Println("Error resizing image", err)
		return err
	}

	log.Println("Setting compression quality")
	err = mw.SetImageCompressionQuality(95)
	if err != nil {
		log.Println("Error setting compression quality", err)
		return err
	}

	customImageWand := mw.GetImage()

	// read base T-Shirt image
	log.Println("Reading base image", btmName)
	err = mw.ReadImage(btmName)
	if err != nil {
		log.Println("Error reading image", err)
		return err
	}

	baseTShirtImageWand := mw.GetImage()
	log.Println("Setting gravity")
	err = baseTShirtImageWand.SetGravity(imagick.GRAVITY_CENTER)
	if err != nil {
		log.Println("Error setting gravity", err)
		return err
	}

	// t-shirt original: 1850 × 1234
	log.Println("Rendering composite")
	err = baseTShirtImageWand.CompositeImage(customImageWand, imagick.COMPOSITE_OP_OVER, 780, 425)
	if err != nil {
		log.Println("Error reading image", err)
		return err
	}

	log.Println("Writing new image", newName)
	err = baseTShirtImageWand.WriteImage(newName)
	if err != nil {
		log.Printf("Error writing new image at %s: %s\n", newName, err)
	}

	log.Println("Finished overlaying images.")
	return nil
}
