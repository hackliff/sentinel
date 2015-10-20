package radio

import "testing"

func TestPushbulletAdapter_implements(t *testing.T) {
	var _ Adapter = &PushbulletAdapter{}
}
