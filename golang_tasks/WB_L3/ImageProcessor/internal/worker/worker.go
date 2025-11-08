package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	kafka "github.com/wb-go/wbf/kafka"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	_ "embed"
)

type ImageTask struct {
	ID  string `json:"id"`
	Ext string `json:"ext"`
}

type Worker struct {
	Consumer *kafka.Consumer
}

func NewWorker(consumer *kafka.Consumer) *Worker {
	return &Worker{Consumer: consumer}
}

func (w *Worker) Start(ctx context.Context) {
	out := make(chan []byte)
	go w.consumeKafka(ctx, out)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-out:
			if !ok {
				return
			}
			if err := w.processMessage(msg); err != nil {
				log.Println("Worker error:", err)
			}
		}
	}
}

func (w *Worker) consumeKafka(ctx context.Context, out chan<- []byte) {
	defer close(out)
	for {
		msg, err := w.Consumer.Fetch(ctx)
		if err != nil {
			log.Println("Kafka fetch error:", err)
			return
		}
		out <- msg.Value
		w.Consumer.Commit(ctx, msg)
	}
}

func (w *Worker) processMessage(msg []byte) error {
	var task ImageTask
	if err := json.Unmarshal(msg, &task); err != nil {
		return fmt.Errorf("unmarshal Kafka message: %w", err)
	}

	log.Printf("Got task: id=%s, ext=%s\n", task.ID, task.Ext)
	return w.processImage(task.ID, task.Ext)
}

func (w *Worker) processImage(id, ext string) error {
	origPath := filepath.Join("storage", "original", id+ext)
	procPath := filepath.Join("storage", "processed", id+ext)
	thumbPath := filepath.Join("storage", "thumbnails", id+ext)

	log.Println("📂 Processing:", origPath)

	os.MkdirAll(filepath.Dir(procPath), os.ModePerm)
	os.MkdirAll(filepath.Dir(thumbPath), os.ModePerm)

	src, err := imaging.Open(origPath)
	if err != nil {
		return fmt.Errorf("open image: %w", err)
	}

	//Resize
	resized := imaging.Resize(src, src.Bounds().Dx()/2, 0, imaging.Lanczos)

	//Watermark
	watermarked := addWatermark(resized, "WB")

	//Save processed
	if err := imaging.Save(watermarked, procPath); err != nil {
		return fmt.Errorf("save processed: %w", err)
	}

	// миниатюра
	thumbnail := imaging.Thumbnail(src, 200, 200, imaging.Lanczos)
	if err := imaging.Save(thumbnail, thumbPath); err != nil {
		return fmt.Errorf("save thumbnail: %w", err)
	}

	log.Println("Successfully processed:", id)
	return nil
}

func addWatermark(img image.Image, text string) image.Image {
	rgba := imaging.Clone(img)
	bounds := rgba.Bounds()

	// позиция watermark
	x := bounds.Dx() - len(text)*10 - 10
	y := bounds.Dy() - 10

	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(color.NRGBA{255, 255, 255, 180}),
		Face: basicfont.Face7x13,
		Dot: fixed.Point26_6{
			X: fixed.I(x),
			Y: fixed.I(y),
		},
	}

	d.DrawString(text)
	return rgba
}
