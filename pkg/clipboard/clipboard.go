package clipboard

import (
	"log"
	"time"

	"golang.design/x/clipboard"
)

func Copy(msg string) error {
	log.Println("Initializing clipboard")
	err := clipboard.Init()
	if err != nil {
		log.Println("Clipboard checking failed:", err)
		return err
	}
	log.Println("Writing to clipboard")
	clipboard.Write(clipboard.FmtText, []byte(msg))
	// TODO: Technical Debt (If dont sleep, it can not save to clipboard)
	time.Sleep(time.Second)
	log.Println("Copied to clipboard")
	return nil
}
