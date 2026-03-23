//ff:func feature=history type=util control=sequence
//ff:what FileHistory가 존재하면 Hint를 추가

package history

func addHintIfExists(histories map[string]*FileHistory, relPath string, hint Hint) {
	h, ok := histories[relPath]
	if !ok {
		return
	}
	h.Hints = append(h.Hints, hint)
}
