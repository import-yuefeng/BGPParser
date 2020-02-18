package analysis

type AspathList []string
type BGPInfo struct {
	Aspath     AspathList
	Prefix     []string
	Aspath2str string
	Hashcode   string
	isSorted   bool
	content    string
}

func NewBGPInfo(content string) *BGPInfo {
	return &BGPInfo{
		content: content,
	}
}

func (b *BGPInfo) AnalysisBGPData() {
	b.FindPrefix()
	b.FindAsPath()
	b.CleanContent()
	b.SortASpathBySize()
	b.ConvertHashcode()
}
