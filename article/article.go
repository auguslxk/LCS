package article

import (
	"sort"
	"sync"

	"github.com/auguslxk/LCS/lib"
)

var endPunct = map[rune]bool{
	'。': true,
	'.': true,
	'!': true,
	'！': true,
	'?': true,
	'？': true,
}

type Article struct {
	CharSize      int
	DuplicateSize int
	Sentences     []Sentence
}

type Sentence struct {
	Id             int // 标记是第几句
	Content        []rune
	DuplicateId    int
	DuplicateRatio float64
	DupIndex       []int
}

type Result struct {
	SentenceAId int
	SentenceBId int
	SentenceA   []rune
	SentenceB   []rune
	LCS         []rune
	LCSSize     int
	IndexA      []int
	IndexB      []int
}

type ResultList []Result

func (rl ResultList) Len() int           { return len(rl) }
func (rl ResultList) Swap(i, j int)      { rl[i], rl[j] = rl[j], rl[i] }
func (rl ResultList) Less(i, j int) bool { return rl[i].SentenceAId < rl[j].SentenceAId }

func ArticleInit(article string) *Article {
	sentence, articleSize := splitArticle(article)
	return &Article{articleSize, 0, sentence}
}

func splitArticle(article string) ([]Sentence, int) {
	sentence := []Sentence{}
	articleSize := 0
	id := 0
	art := []rune(article)
	cur := &Sentence{id, []rune{}, -1, 0, []int{}}
	for index, char := range art {
		cur.Content = append(cur.Content, art[index])
		if _, ok := endPunct[char]; ok {
			sentence = append(sentence, *cur)
			id += 1
			cur = &Sentence{id, []rune{}, -1, 0, []int{}}
		}
	}
	sentence = append(sentence, *cur)
	for _, val := range sentence {
		articleSize += len(val.Content)
	}

	return sentence, articleSize
}

func DuplicateChecking(arta, artb *Article) []Result {
	lock := &sync.Mutex{}
	wg := sync.WaitGroup{}
	resultList := ResultList{}
	// 开100个
	coroutines := 100
	groupSize := len(arta.Sentences) / coroutines
	if groupSize == 0 {
		//arta句子数量小于100个
		coroutines = len(arta.Sentences)
		groupSize = 1
	}

	for g := 0; g < coroutines-1; g++ {
		wg.Add(1)
		go checkUtil(g*groupSize, (g+1)*groupSize, arta, artb, &resultList, lock, &wg)
	}

	wg.Add(1)
	go checkUtil((coroutines-1)*groupSize, len(arta.Sentences), arta, artb, &resultList, lock, &wg)
	wg.Wait()
	sort.Sort(resultList)
	statistics(arta, artb, resultList)
	return resultList
}

func checkUtil(start, end int, arta, artb *Article, resultList *ResultList, lock *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	matchedMax := 0.0
	resList := ResultList{}
	for start < end {
		result := Result{}
		matched := false
		for index, _ := range artb.Sentences {
			lcs, indexA, indexB := lib.GetLCS(arta.Sentences[start].Content, artb.Sentences[index].Content)
			if len(lcs) == 0 {
				continue
			}
			match := float64(len(lcs)) / float64(lib.Max(len(arta.Sentences[start].Content), len(artb.Sentences[index].Content)).(int))
			if match > matchedMax {
				matched = true
				matchedMax = match
				result.LCS = lcs
				result.LCSSize = len(lcs)
				result.SentenceAId = arta.Sentences[start].Id
				result.SentenceBId = artb.Sentences[index].Id
				result.SentenceA = arta.Sentences[start].Content
				result.SentenceB = artb.Sentences[index].Content
				result.IndexA = indexA
				result.IndexB = indexB
			}
		}
		start += 1
		if matched {
			resList = append(resList, result)
		}
	}
	lock.Lock()
	*resultList = append(*resultList, resList...)
	lock.Unlock()
}

func statistics(arta, artb *Article, list ResultList) {
	dupId := 1
	for _, result := range list {
		// 为日后debug写下一个注释
		// Id和index的值是一样的 可能以后手抽动了哪里就对不上了😂
		sA := &arta.Sentences[result.SentenceAId]
		sB := &artb.Sentences[result.SentenceBId]
		sA.DuplicateRatio = float64(result.LCSSize) / float64(len(sA.Content))
		sB.DuplicateRatio = float64(result.LCSSize) / float64(len(sB.Content))
		if sA.DuplicateRatio > 0.3 {
			sA.DupIndex = result.IndexA
			sB.DupIndex = result.IndexB
			sA.DuplicateId = dupId
			sB.DuplicateId = dupId
			dupId += 1
			arta.DuplicateSize += result.LCSSize
			artb.DuplicateSize += result.LCSSize
		}
	}
}
