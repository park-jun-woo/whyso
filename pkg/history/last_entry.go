//ff:func feature=history type=util control=sequence
//ff:what FileHistory의 마지막 ChangeEntry 포인터를 반환
package history

func lastEntry(h *FileHistory) *ChangeEntry {
	if len(h.History) == 0 {
		return nil
	}
	return &h.History[len(h.History)-1]
}
