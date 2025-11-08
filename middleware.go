package errors

//
//import (

//)
//
//type ErrorResponse struct {
//	// 是否处理成功
//	Success bool `json:"success"`
//	// 错误码
//	ErrCode *ErrCode `json:"errCode"`
//	// 错误详情
//	Message string `json:"message"`
//	// 请求ID
//	RequestID string `json:"requestId,omitempty"`
//	// 纬度ID
//	SpanID string `json:"spanId,omitempty"`
//	// 跟踪ID
//	TraceID string `json:"traceID,omitempty"`
//	// 详情信息
//	Details map[string]any `json:"details"`
//	// 时间戳
//	Timestamp string `json:"timestamp"`
//}
//
//type Handler struct {
//	// 日志适配器
//	l logger.Logger
//	// 是否显示错误详情
//	showDetails bool
//	// 是有隐藏内部
//	hideInternal bool
//	// 是否启用监控
//	enableMonitor bool
//	// 环境
//	environment string
//}
//
//func NewHandler(l logger.Logger, opts ...HandlerOption) *Handler {
//	h := &Handler{
//		l: l,
//	}
//
//	for _, opt := range opts {
//		opt(h)
//	}
//
//	return h
//}
//
//// RecoveryMiddleware 适配Gin框架的错误恢复中间件
//func (h *Handler) RecoveryMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		defer func() {
//			r := recover()
//			if r == nil {
//				return
//			}
//
//			panicError := h.createPanicError(r, debug.Stack())
//			h.handleError(c, panicError)
//			h.l.Error("panic recovered",
//				logger.AnyField("panic", r),
//				logger.StringField("stack", string(debug.Stack())),
//				logger.StringField("path", c.Request.URL.Path))
//
//			c.Abort()
//		}()
//
//		c.Next()
//	}
//}
//
//// HandleError 处理错误，包括日志的记录、监控的记录和错误响应的构建
//func (h *Handler) handleError(c *gin.Context, err Error) {
//	h.recordError(c, err)
//	h.buildErrorResponse(c, err)
//}
//
//// recordError 记录错误到日志和监控系统
//func (h *Handler) recordError(c *gin.Context, err Error) {
//	if h.enableMonitor {
//		// TODO 记录错误到监控系统
//	}
//
//	// 记录机构化日志
//	fields := []logger.Field{
//		logger.StringField("code", err.Code()),
//		logger.StringField("type", err.Type().String()),
//		logger.StringField("path", c.Request.URL.Path),
//		logger.StringField("method", c.Request.Method),
//		logger.IntField("http_status", err.HttpStatus()),
//		logger.TimeField("timestamp", err.Timestamp()),
//	}
//
//	// 添加其它的信息
//	// 记录SpanID信息，部分会叫做RequestID
//	requestID := c.GetString("request_id")
//	if requestID != "" {
//		fields = append(fields, logger.StringField("request_id", requestID))
//	}
//
//	spanID := c.GetString("span_id")
//	if spanID != "" {
//		fields = append(fields, logger.StringField("span_id", spanID))
//	}
//
//	// 记录客户端IP
//	if client := c.ClientIP(); client != "" {
//		fields = append(fields, logger.StringField("client_ip", client))
//	}
//
//	// 记录其他的元数据信息
//	if metadata := err.Metadata(); len(metadata) > 0 {
//		fields = append(fields, logger.AnyField("metadata", metadata))
//	}
//
//	// 根据错误类型来决定日志记录的级别
//	switch err.Type() {
//	case ErrTypeInternal, ErrTypeTimeout, ErrTypeExternal:
//		h.l.Error("internal error", fields...)
//	case ErrTypeBusiness, ErrTypeValidation:
//		h.l.Warn("business error", fields...)
//	default:
//		h.l.Info("unknown error", fields...)
//	}
//}
//
//// buildErrorResponse 构建错误响应
//func (h *Handler) buildErrorResponse(c *gin.Context, err Error) {
//	errResponse := ErrorResponse{
//		Success:   false,
//		Message:   err.Message(),
//		Timestamp: err.Timestamp().Format(time.RFC3339),
//	}
//
//	requestID := c.GetString("request_id")
//	if requestID != "" {
//		errResponse.RequestID = requestID
//	}
//
//	spanID := c.GetString("span_id")
//	if spanID != "" {
//		errResponse.SpanID = spanID
//	}
//
//	traceID := c.GetString("trace_id")
//	if traceID != "" {
//		errResponse.TraceID = traceID
//	}
//
//	status := err.HttpStatus()
//	if status == 0 {
//		status = http.StatusInternalServerError
//	}
//
//	c.JSON(status, errResponse)
//}
//
//// createPanicError 创建panic错误
//func (h *Handler) createPanicError(panicValue any, stack []byte) Error {
//	var panicMessage string
//	switch v := panicValue.(type) {
//	case error:
//		panicMessage = v.Error()
//	case string:
//		panicMessage = v
//	default:
//		panicMessage = fmt.Sprintf("%v", v)
//	}
//
//	return NewBuilder().
//		WithCode(&ErrCode{
//			Code:       "PANIC_RECOVERED",
//			Message:    "Service encountered a panic and recovered",
//			HttpStatus: http.StatusInternalServerError,
//			Type:       ErrTypeInternal,
//		}).
//		WithMessage(fmt.Sprintf("panic recovered: %s", panicMessage)).
//		WithMetadata("panic_value", panicValue).
//		WithMetadata("panic_type", fmt.Sprintf("%T", panicValue)).
//		WithMetadata("stack_trace", string(stack)).
//		Build()
//}
//
//// createPanicErrorSimple 创建简化版的panic错误
//func (h *Handler) createPanicErrorSimple(panicValue any, stack []byte) Error {
//	panicMsg := fmt.Sprintf("%v", panicValue)
//	return Newf(ErrInternal, "panic recovered: %s", panicMsg).
//		WithMetadata("panic_value", panicValue).
//		WithMetadata("panic_type", fmt.Sprintf("%T", panicValue)).
//		WithMetadata("stack_trace", string(stack))
//}
//
////// LoggingMiddleware 适配Gin框架中记录日志的中间件
////func (h *Handler) LoggingMiddleware() gin.HandlerFunc {
////	return func(c *gin.Context) {
////		start := time.Now()
////		h.l.Debug("start request",
////			logger.StringField("method", c.Request.Method),
////			logger.StringField("path", c.Request.URL.Path),
////			logger.StringField("client_ip", c.ClientIP()),
////			logger.StringField("user_agent", c.Request.UserAgent()),
////			logger.StringField("request_id", c.GetString("request_id")),
////			logger.TimeField("timestamp", start))
////
////		c.Next()
////		latency := time.Since(start)
////		status := c.Writer.Status()
////		fields := []logger.Field{
////			logger.StringField("method", c.Request.Method),
////			logger.StringField("path", c.Request.URL.Path),
////			logger.StringField("client_ip", c.ClientIP()),
////			logger.StringField("user_agent", c.Request.UserAgent()),
////			logger.StringField("request_id", c.GetString("request_id")),
////			logger.DurationField("latency", latency),
////			logger.TimeField("timestamp", start),
////			logger.IntField("status", status),
////		}
////
////		// 响应的大小
////		if size := c.Writer.Size(); size > 0 {
////			fields = append(fields, logger.IntField("size", size))
////		}
////
////		// 根据状态码来决定使用什么级别
////		switch {
////		case status >= 500:
////			// 服务器内部错误
////			h.l.Error("server error", fields...)
////		case status >= 400:
////			// 客户端错误
////			h.l.Warn("client error", fields...)
////		default:
////			// 正常请求
////			h.l.Info("request completed", fields...)
////		}
////	}
////}
//
//func (h *Handler) ErrorMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 先处理请求
//		c.Next()
//
//		// 处理错误
//		if c.Errors != nil && len(c.Errors) > 0 {
//			for _, err := range c.Errors {
//				h.handleError(c, err.Err.(Error))
//			}
//		}
//	}
//}
//
//// TimeoutMiddleware 适配Gin框架的请求超时处理中间件
//func (h *Handler) TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		ctx, cancel := context.WithTimeout(context.Background(), timeout)
//		defer cancel()
//
//		c.Request.WithContext(ctx)
//		done := make(chan struct{})
//		go func() {
//			defer close(done)
//			c.Next()
//		}()
//
//		select {
//		case <-done:
//		// 正常处理完整，未超时
//		case <-ctx.Done():
//			// 请求处理超时了
//			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
//				timeoutErr := Newf(ErrTimeout, "request timed out").
//					WithMetadata("method", c.Request.Method).
//					WithMetadata("path", c.Request.URL.Path)
//				h.handleError(c, timeoutErr)
//				c.Abort()
//			}
//		}
//	}
//}
