package dialog

type Dialog struct {
	Title       string
	Description string
	Item        string
}

func (dialog *Dialog) SerializeDialog() []byte {
	body := make([]any, 0, 2)

	if len(dialog.Item) != 0 {
		body = append(body, struct {
			Type string
			Item any
		}{
			Type: "minecraft:item",
			Item: struct {
				Id string
			}{
				Id: "minecraft:clock",
			},
		})
	}

	if len(dialog.Description) != 0 {
		body = append(body, struct {
			Type     string
			Contents string
		}{
			Type:     "minecraft:plain_message",
			Contents: dialog.Description,
		})
	}

	return serializeNbt(&struct {
		Type      string
		Title     string
		Body      []any
		Closeable uint8 `nbt:"can_close_with_escape"`
		Columns   uint32
	}{
		Type:      "minecraft:notice",
		Title:     dialog.Title,
		Body:      body,
		Closeable: 0,
		Columns:   0,
	})
}
