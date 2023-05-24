package internal

type Message struct {
  Id       string      `json:"id"`
  Role     string      `json:"role"`
  Name     *string     `json:"name"`
  Metadata interface{} `json:"metadata"`
}

type Content struct {
  ContentType string   `json:"content_type"`
  Parts       []string `json:"parts"`
}

type NodeMessage struct {
  Id         string      `json:"id"`
  Author     Message     `json:"author"`
  CreateTime float64     `json:"create_time"`
  UpdateTime *float64    `json:"update_time"`
  Content    Content     `json:"content"`
  EndTurn    *bool       `json:"end_turn"`
  Weight     float64     `json:"weight"`
  Metadata   interface{} `json:"metadata"`
  Recipient  string      `json:"recipient"`
}

type Node struct {
  Id       string       `json:"id"`
  Message  *NodeMessage `json:"message"`
  Parent   *string      `json:"parent"`
  Children []string     `json:"children"`
}

type Conversation struct {
  Id                string          `json:"id"`
  Title             string          `json:"title"`
  CreateTime        float64         `json:"create_time"`
  UpdateTime        float64         `json:"update_time"`
  Mapping           map[string]Node `json:"mapping"`
  ModerationResults []interface{}   `json:"moderation_results"`
  CurrentNode       string          `json:"current_node"`
  PluginIds         *string         `json:"plugin_ids"`
}
