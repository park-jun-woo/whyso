//ff:func feature=history type=util control=iteration dimension=1
//ff:what parentUuid 체인을 역추적하여 원본 사용자 메시지를 찾음
package history

const maxChainDepth = 100

// FindUserRequest traces the parentUuid chain backwards to find the original user message.
func FindUserRequest(idx RecordIndex, recordUUID string) string {
	current, ok := idx[recordUUID]
	if !ok {
		return ""
	}

	for i := 0; i < maxChainDepth; i++ {
		if current.ParentUUID == nil && current.IsUserMessage() {
			return current.UserContent()
		}
		if current.ParentUUID == nil {
			return ""
		}

		parent, ok := idx[*current.ParentUUID]
		if !ok {
			return ""
		}

		if parent.IsUserMessage() {
			return parent.UserContent()
		}
		current = parent
	}

	return ""
}
