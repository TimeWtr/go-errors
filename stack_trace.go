package errors

import (
	"runtime"
	"strconv"
	"strings"
)

// captureStackTrace 捕获堆栈跟踪信息
func captureStackTrace(skip int) string {
	const depth = 32
	var pcs [depth]uintptr

	// 获取调用堆栈
	n := runtime.Callers(skip+2, pcs[:])
	if n == 0 {
		return ""
	}

	// 获取堆栈帧
	frames := runtime.CallersFrames(pcs[:n])

	var stack strings.Builder
	stack.WriteString("Stack Trace:\n")

	// 遍历所有帧
	for {
		frame, more := frames.Next()

		// 格式化堆栈信息
		stack.WriteString("  ")
		stack.WriteString(frame.Function)
		stack.WriteString("\n    ")
		stack.WriteString(frame.File)
		stack.WriteString(":")
		stack.WriteString(strconv.Itoa(frame.Line))
		stack.WriteString("\n")

		if !more {
			break
		}
	}

	return stack.String()
}

// getCallerInfo 获取调用者信息（文件名和行号）
func getCallerInfo(skip int) (string, int) {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown", 0
	}

	// 获取函数名
	funcName := runtime.FuncForPC(pc).Name()

	// 简化文件路径，只显示最后两部分
	parts := strings.Split(file, "/")
	if len(parts) > 2 {
		file = strings.Join(parts[len(parts)-2:], "/")
	}

	// 简化函数名
	funcParts := strings.Split(funcName, "/")
	funcName = funcParts[len(funcParts)-1]

	return file + ":" + funcName, line
}

// getSimplifiedStackTrace 获取简化的堆栈跟踪（生产环境友好）
func getSimplifiedStackTrace(skip int, maxDepth int) string {
	if maxDepth <= 0 {
		maxDepth = 8
	}

	var pcs [32]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	if n == 0 {
		return ""
	}

	frames := runtime.CallersFrames(pcs[:n])
	var stack strings.Builder
	stack.WriteString("Simplified Stack:\n")

	count := 0
	for {
		if count >= maxDepth {
			break
		}

		frame, more := frames.Next()

		// 过滤掉系统库和错误库本身的调用
		if shouldIncludeFrame(frame) {
			stack.WriteString("  ")
			stack.WriteString(extractSimpleFunctionName(frame.Function))
			stack.WriteString(" (")
			stack.WriteString(extractFileName(frame.File))
			stack.WriteString(":")
			stack.WriteString(strconv.Itoa(frame.Line))
			stack.WriteString(")\n")
			count++
		}

		if !more {
			break
		}
	}

	return stack.String()
}

// shouldIncludeFrame 判断是否应该包含该堆栈帧
func shouldIncludeFrame(frame runtime.Frame) bool {
	// 排除系统库
	if strings.Contains(frame.File, "/usr/lib/") ||
		strings.Contains(frame.File, "runtime/") ||
		strings.Contains(frame.File, "reflect/") {
		return false
	}

	// 排除错误库本身的调用
	if strings.Contains(frame.File, "pkg/errors/") {
		return false
	}

	// 排除一些常见的第三方库（根据实际项目调整）
	if strings.Contains(frame.Function, "vendor/") {
		return false
	}

	return true
}

// extractSimpleFunctionName 提取简化的函数名
func extractSimpleFunctionName(funcName string) string {
	// 移除包路径，只保留最后的函数名
	parts := strings.Split(funcName, "/")
	lastPart := parts[len(parts)-1]

	// 进一步简化
	dotParts := strings.Split(lastPart, ".")
	if len(dotParts) > 1 {
		return dotParts[len(dotParts)-1]
	}

	return lastPart
}

// extractFileName 提取文件名
func extractFileName(filePath string) string {
	parts := strings.Split(filePath, "/")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return filePath
}

// GetCurrentFunctionName 获取当前函数名（用于调试）
func GetCurrentFunctionName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	funcObj := runtime.FuncForPC(pc)
	if funcObj == nil {
		return "unknown"
	}

	return extractSimpleFunctionName(funcObj.Name())
}
