package rtimer

type RTimer struct {
	nTickTerm int32 //触发时间
	nCurTick  int32 //当前时间片
	nCount    int32 //触发次数 0 代表无限触发
	nGoCount  int32 //累计流逝时间
	bOper     bool
}

func NewRTimer() *RTimer {
	ctimer := &RTimer{}
	ctimer.CleanUp()
	return ctimer
}
func (c *RTimer) CleanUp() {
	c.nTickTerm = 0
	c.bOper = false
	c.nCurTick = 0
	c.nCount = 0
	c.nGoCount = 0
}

//开始 nTerm触发时间
// nCount  触发次数，0 代表无限触发
func (c *RTimer) BeginTimer(nTerm int32, nCount int32) {
	c.bOper = true
	c.nTickTerm = nTerm
	c.nCurTick = nTerm
	c.nCount = nCount
	c.nGoCount = 0
}
func (c *RTimer) IsRun() bool {
	return c.bOper
}
func (c *RTimer) CountingTimer(nNow int32) bool {
	if !c.bOper {
		return false
	}
	c.nCurTick -= nNow
	c.nGoCount += nNow
	if c.nCurTick < 0 {
		c.nCurTick += c.nTickTerm
		if c.nCount != 0 {
			c.nCount--
			if c.nCount == 0 {
				c.CleanUp()
			}
		}
		return true
	}
	return false
}

//获取剩余时间
func (c *RTimer) GetRemainingTime() int32 {
	return c.nCurTick
}
