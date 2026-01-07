package prompt

// TranslateSystemPrompt 翻译文本的 system prompt
//
// 用于将非英文文本（中文、俄文等）翻译为英文。
// 从嵌入的模板文件中加载。
var TranslateSystemPrompt = MustLoadTemplate("translate.md")
