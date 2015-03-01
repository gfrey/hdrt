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
			if !ok {
				continue
			}
			img, _ := disp.AllocImage(si.Rect, disp.ScreenImage.Pix, false, draw.Nofill)
			buf := make([]byte, 4*(si.Rect.Max.Y-si.Rect.Min.Y)*(si.Rect.Max.X-si.Rect.Min.X))
			for y := si.Rect.Min.Y; y < si.Rect.Max.Y; y++ {
				for x := si.Rect.Min.X; x < si.Rect.Max.X; x++ {
					o := si.PixOffset(x, y)
					
					buf[o+0] = si.Pix[o+2]
					buf[o+1] = si.Pix[o+1]
					buf[o+2] = si.Pix[o+0]
					buf[o+3] = si.Pix[o+3]
				}
			}
			img.Load(si.Rect, buf)
			disp.ScreenImage.Draw(si.Rect, img, nil, si.Rect.Min)
			disp.Flush()
			img.Free()
		case c := <-kbdIn.C:
			if c == 'q' {
				wrld.Abort()
				return nil
			}
		}
	}

	return err
}
