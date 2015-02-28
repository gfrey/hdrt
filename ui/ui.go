package ui

import (
	"fmt"

	"github.com/gfrey/hdrt"

	"9fans.net/go/draw"
)

func Run(wrld *hdrt.World) error {
	res := fmt.Sprintf("%dx%d", wrld.Viewplane.ResX, wrld.Viewplane.ResY)
	disp, err := draw.Init(nil, "", "hdrt", res)
	if err != nil {
		return err
	}
	defer disp.Close()

	if err := disp.Attach(draw.Refmesg); err != nil {
		return err
	}

	kbdIn := disp.InitKeyboard()
	sImgOutChan := wrld.Render()
	for {
		select {
		case si, ok := <-sImgOutChan:
			if !ok { continue }
			img, _ := disp.AllocImage(si.Rect, disp.ScreenImage.Pix, false, 0)
			_, err := img.Load(si.Rect, si.Buf)
			if err != nil {
				return err
			}
			disp.ScreenImage.Draw(si.Rect, img, nil, si.Rect.Min)
			disp.Flush()
		case c := <-kbdIn.C:
			if c == 'q' {
				wrld.Abort()
				return nil
			}
		}
	}

	return err
}
