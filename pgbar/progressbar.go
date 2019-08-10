package pgbar

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// ProgressBar is a thread-safe, simple
// progress bar
type ProgressBar struct {
	max               int // max number of the counter
	size              int // size of the saucer
	currentNum        int
	currentPercent    int
	lastPercent       int
	currentSaucerSize int

	lastShown time.Time
	startTime time.Time
	w         io.Writer

	theme []string

	sync.RWMutex
	msg string
}

func (p *ProgressBar) SetTheme(theme []string) {
	p.Lock()
	defer p.Unlock()
	p.theme = theme
}

// New returns a new ProgressBar
// with the specified maximum
func New(max int) *ProgressBar {
	return &ProgressBar{
		max:       max,
		size:      100,
		theme:     []string{"█", " ", "|", "|"},
		w:         os.Stdout,
		lastShown: time.Now(),
		startTime: time.Now(),
	}
}

// Reset will reset the clock that is used
// to calculate current time and the time left.
func (p *ProgressBar) Reset() {
	p.Lock()
	defer p.Unlock()
	p.lastShown = time.Now()
	p.startTime = time.Now()
	p.currentNum = 0
}

// SetMax sets the total number of the progress bar
func (p *ProgressBar) SetMax(num int) {
	p.Lock()
	defer p.Unlock()
	p.max = num
}

// SetSize sets the size of the progress bar.
func (p *ProgressBar) SetSize(size int) {
	p.Lock()
	defer p.Unlock()
	p.size = size
}

// SetWriter will specify a different writer than os.Stdout
func (p *ProgressBar) SetWriter(w io.Writer) {
	p.Lock()
	defer p.Unlock()
	p.w = w
}

// Add with increase the current count on the progress bar
func (p *ProgressBar) Add(num int, msg string) error {
	p.Lock()
	defer p.Unlock()
	updateBar := false
	if msg != "" {
		if p.msg != msg {
			p.msg = msg
			updateBar = true
		}
	}
	//if p.max == 0 {
	//	return errors.New("max must be greater than 0")
	//}
	p.currentNum += num
	//percent := float64(p.currentNum) / float64(p.max)
	p.currentSaucerSize = num
	//p.currentSaucerSize = int(percent *100/p.size)
	//p.currentPercent = int(percent * 100)
	p.currentPercent = num
	if !updateBar {
		updateBar = p.currentPercent != p.lastPercent && p.currentPercent >= 0
	}

	p.lastPercent = p.currentPercent
	//if p.currentNum > p.max {
	//	return errors.New("current number exceeds max")
	//}
	theme1 := ""
	if p.size-p.currentSaucerSize > 0 {
		theme1 = strings.Repeat(p.theme[1], p.size-p.currentSaucerSize)
	}
	if updateBar {
		leftTime := time.Since(p.startTime).Seconds() / float64(p.currentNum) * (float64(p.max) - float64(p.currentNum))
		s := fmt.Sprintf("\r%4d%% %s%s%s%s [%s:%s]%s",
			p.currentPercent,
			p.theme[2],
			strings.Repeat(p.theme[0], p.currentSaucerSize),
			theme1,
			p.theme[3],
			(time.Duration(time.Since(p.startTime).Seconds()) * time.Second).String(),
			(time.Duration(leftTime) * time.Second).String(),
			p.msg,
		)
		_, err := io.WriteString(p.w, s)
		if err != nil {
			return err
		}

		if f, ok := p.w.(*os.File); ok {
			f.Sync()
		}
	}
	return nil
}
