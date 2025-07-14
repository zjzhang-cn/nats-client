# Git提交信息参考

## 本次修改的建议提交信息

### 标准格式提交信息
```
feat(ui): add NATS server address selection component

- Add dynamic server address selection dropdown with preset options
- Add custom WebSocket URL input with validation  
- Add editable authentication token configuration
- Improve UI layout with better visual hierarchy
- Enhance connection logic with dynamic configuration
- Add TypeScript helper functions for configuration management

Breaking Changes: None
Backward Compatible: Yes
```

### 简洁版提交信息
```
feat: add NATS server address selection component

Add dynamic server selection, custom URL input, and token configuration
for flexible NATS connection management without code changes.
```

### 详细版提交信息
```
feat(ui): implement dynamic NATS server configuration

Features:
- Server address dropdown with preset options (default, local)
- Custom WebSocket URL input with ws:// and wss:// validation
- Editable authentication token field with default fallback
- Responsive UI with conditional custom input display

Technical:
- Add getSelectedServerUrl() helper function
- Add getAuthToken() helper function  
- Add handleServerSelectChange() event handler
- Enhance connectToNats() with dynamic configuration
- Update setupUIHandlers() for new UI components

Files modified:
- html/index.html: Add server selection and token UI components
- html/src/index.ts: Implement dynamic configuration logic

Testing:
- ✅ TypeScript compilation successful
- ✅ Build process completes without errors
- ✅ All existing functionality preserved

Co-authored-by: GitHub Copilot <copilot@github.com>
```

## 提交建议

推荐使用以下Git命令序列：

```bash
# 查看修改状态
git status

# 添加修改的文件
git add html/index.html html/src/index.ts

# 添加文档文件
git add CHANGELOG.md docs/development-log-2025-07-11.md html/CHANGES.md

# 提交修改
git commit -m "feat(ui): add NATS server address selection component

- Add dynamic server address selection dropdown with preset options
- Add custom WebSocket URL input with validation  
- Add editable authentication token configuration
- Improve UI layout with better visual hierarchy
- Enhance connection logic with dynamic configuration

Backward Compatible: Yes"

# 可选: 创建标签
git tag -a v1.1.0 -m "Release v1.1.0: Add NATS server address selection"
```
