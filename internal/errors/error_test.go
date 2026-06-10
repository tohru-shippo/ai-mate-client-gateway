package errors

import (
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLocalizeUsesClientTone(t *testing.T) {
	err := Localize(BadRequest(""), "zh-CN")
	if err == nil {
		t.Fatal("expected localized error")
	}
	if err.Message != "请求参数有点不对，再检查一下吧" {
		t.Fatalf("unexpected message: %s", err.Message)
	}
}

func TestSupportedLocaleMessages(t *testing.T) {
	locales := []string{"en-US", "zh-CN", "zh-TW", "ja-JP", "ko-KR", "ru-RU"}
	for _, locale := range locales {
		t.Run(locale, func(t *testing.T) {
			message := MessageForCode(10005, locale)
			if message == "" {
				t.Fatal("expected message")
			}
			if locale != "en-US" && message == MessageForCode(10005, "en-US") {
				t.Fatalf("expected locale-specific message, got fallback: %s", message)
			}
		})
	}
}

func TestLocalizeKeepsFieldValidationMessage(t *testing.T) {
	st := status.New(codes.InvalidArgument, "invalid argument")
	st, detailErr := st.WithDetails(
		&errdetails.ErrorInfo{Reason: coreErrorReason, Metadata: map[string]string{coreCodeKey: "10005"}},
		&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{Field: "userId", Description: "required"}}},
	)
	if detailErr != nil {
		t.Fatalf("build status details: %v", detailErr)
	}
	err := Localize(FromCoreError(st.Err(), "zh-CN"), "zh-CN")
	if err == nil {
		t.Fatal("expected localized error")
	}
	if err.Message != "用户 ID 不能为空" {
		t.Fatalf("unexpected message: %s", err.Message)
	}
}
