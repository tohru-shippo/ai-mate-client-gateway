package errors

import (
	"net/http"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	coreErrorReason = "AI_MATE_ERROR"
	coreCodeKey     = "code"
)

var clientMessages = map[string]map[int]string{
	"en-US": {
		10000: "Business state does not allow this action",
		10001: "Data already exists",
		10002: "Data not found",
		10003: "Deletion is not allowed",
		10004: "Update is not allowed",
		10005: "Parameter validation failed",
		20001: "Account already exists",
		20002: "Authentication failed",
		20003: "Linked user not found",
		20004: "User not found",
		20005: "Account not found",
		21001: "Conversation not found or access denied",
		21002: "Voice can only be generated for AI messages",
		21003: "Text content is empty",
		21004: "TTS settings are not configured",
		22001: "Character card not found",
		22002: "Character profile is missing",
		22003: "Character card is not published",
		22004: "Character card version not found",
		22005: "World card is not published",
		22006: "AI setting not found",
		22007: "AI setting is not published",
		23001: "Uploaded file is empty",
		23002: "Failed to read file",
		23003: "Asset not found",
		24001: "TTS provider not found",
		24002: "Image provider not found",
		24003: "Model mapping is not configured",
		24004: "TTS is not implemented",
		24005: "Image generation is not implemented",
		24006: "API key is not configured",
		24007: "TTS service failed",
		24008: "LLM generation failed",
		24009: "SSE connection closed",
		24010: "Memory is being organized. Please try again later.",
		25001: "Prompt template not found",
		26001: "Not enough points. Please recharge and try again.",
		26002: "Order not found",
		26003: "Order is already paid",
		26004: "Order has expired",
		27001: "Verification code is invalid",
		27002: "Verification code has expired",
		27003: "Too many requests. Please try again later.",
		28001: "OAuth authorization failed",
		31001: "Chat runtime session not found",
		40004: "Data not found",
		403:   "Access denied",
		404:   "API route not found",
		405:   "HTTP method not allowed",
		50000: "Service error",
		50001: "API is not implemented yet",
		-1:    "Service error",
	},
	"zh-CN": {
		10000: "现在的状态还不能这样操作，换一步试试吧",
		10001: "这条数据已经有啦，不用重复提交",
		10002: "没有找到对应数据，检查一下再试吧",
		10003: "这条内容暂时不能删除",
		10004: "这条内容暂时不能修改",
		10005: "请求参数有点不对，再检查一下吧",
		20001: "这个账号已经注册过啦，直接登录试试吧",
		20002: "登录状态好像失效了，请重新登录一下",
		20003: "没有找到关联用户，重新登录试试吧",
		20004: "没有找到这个用户，检查一下账号信息吧",
		20005: "没有找到这个账号，检查一下再试吧",
		21001: "没有找到这段对话，或者你暂时不能访问它",
		21002: "只有 AI 消息才能生成语音哦",
		21003: "消息内容还是空的，写一点再发送吧",
		21004: "语音设置还没配好，暂时不能生成语音",
		22001: "没有找到这个角色卡，换一个看看吧",
		22002: "角色卡的人设还没准备好，暂时不能使用",
		22003: "这个角色卡还没有发布，先看看别的吧",
		22004: "没有找到这个角色卡版本",
		22005: "这个世界卡还没有发布，先看看别的吧",
		22006: "AI 设定没有找到，暂时不能开始对话",
		22007: "这个 AI 设定还没有发布，先换一个试试吧",
		23001: "还没有选择文件，先选一个再上传吧",
		23002: "文件读取失败了，重新上传试试吧",
		23003: "没有找到这个文件资源",
		24001: "语音服务还没配置好，暂时不能生成语音",
		24002: "图片服务还没配置好，暂时不能生成图片",
		24003: "模型配置还没准备好，稍后再试吧",
		24004: "语音功能还在准备中，先稍等一下",
		24005: "图片生成功能还在准备中，先稍等一下",
		24006: "服务配置还没准备好，稍后再试吧",
		24007: "语音服务刚刚走神了，稍后再试吧",
		24008: "AI 生成失败了，稍等一下再试吧",
		24009: "连接断开了，刷新后再试试吧",
		24010: "我正在整理记忆，暂时不能继续回复，稍等一下再试吧",
		25001: "没有找到对应的 Prompt 模板",
		26001: "积分不足，无法完成本次操作，请充值后再试",
		26002: "没有找到这笔订单，检查一下再试吧",
		26003: "这笔订单已经支付过啦",
		26004: "这笔订单已经过期，请重新下单",
		27001: "验证码不对，再检查一下吧",
		27002: "验证码已经过期，请重新获取",
		27003: "请求过于频繁，请稍后再试",
		28001: "授权失败了，请重新授权试试",
		31001: "这次聊天状态没有找到，重新发送试试吧",
		40004: "没有找到对应数据，检查一下再试吧",
		403:   "暂时没有权限访问这里",
		404:   "这个接口还没有找到，检查一下路径吧",
		405:   "这个请求方式不支持，换个方式试试吧",
		50000: "服务刚刚走神了，稍后再试吧",
		50001: "这个接口还在准备中，先稍等一下",
		-1:    "服务刚刚走神了，稍后再试吧",
	},
	"zh-TW": {
		10000: "現在的狀態還不能這樣操作，換一步試試吧",
		10001: "這筆資料已經有啦，不用重複提交",
		10002: "沒有找到對應資料，檢查一下再試吧",
		10003: "這筆內容暫時不能刪除",
		10004: "這筆內容暫時不能修改",
		10005: "請求參數有點不對，再檢查一下吧",
		20001: "這個帳號已經註冊過啦，直接登入試試吧",
		20002: "登入狀態好像失效了，請重新登入一下",
		20003: "沒有找到關聯使用者，重新登入試試吧",
		20004: "沒有找到這個使用者，檢查一下帳號資訊吧",
		20005: "沒有找到這個帳號，檢查一下再試吧",
		21001: "沒有找到這段對話，或者你暫時不能存取它",
		21002: "只有 AI 訊息才能生成語音喔",
		21003: "訊息內容還是空的，寫一點再送出吧",
		21004: "語音設定還沒配好，暫時不能生成語音",
		22001: "沒有找到這個角色卡，換一個看看吧",
		22002: "角色卡的人設還沒準備好，暫時不能使用",
		22003: "這個角色卡還沒有發布，先看看別的吧",
		22004: "沒有找到這個角色卡版本",
		22005: "這個世界卡還沒有發布，先看看別的吧",
		22006: "AI 設定沒有找到，暫時不能開始對話",
		22007: "這個 AI 設定還沒有發布，先換一個試試吧",
		23001: "還沒有選擇檔案，先選一個再上傳吧",
		23002: "檔案讀取失敗了，重新上傳試試吧",
		23003: "沒有找到這個檔案資源",
		24001: "語音服務還沒設定好，暫時不能生成語音",
		24002: "圖片服務還沒設定好，暫時不能生成圖片",
		24003: "模型設定還沒準備好，稍後再試吧",
		24004: "語音功能還在準備中，先稍等一下",
		24005: "圖片生成功能還在準備中，先稍等一下",
		24006: "服務設定還沒準備好，稍後再試吧",
		24007: "語音服務剛剛走神了，稍後再試吧",
		24008: "AI 生成失敗了，稍等一下再試吧",
		24009: "連線斷開了，重新整理後再試試吧",
		24010: "我正在整理記憶，暫時不能繼續回覆，稍等一下再試吧",
		25001: "沒有找到對應的 Prompt 模板",
		26001: "積分不足，無法完成本次操作，請儲值後再試",
		26002: "沒有找到這筆訂單，檢查一下再試吧",
		26003: "這筆訂單已經支付過啦",
		26004: "這筆訂單已經過期，請重新下單",
		27001: "驗證碼不對，再檢查一下吧",
		27002: "驗證碼已經過期，請重新取得",
		27003: "請求太頻繁了，稍後再試吧",
		28001: "授權失敗了，請重新授權試試",
		31001: "這次聊天狀態沒有找到，重新送出試試吧",
		40004: "沒有找到對應資料，檢查一下再試吧",
		403:   "暫時沒有權限存取這裡",
		404:   "這個接口還沒有找到，檢查一下路徑吧",
		405:   "這個請求方式不支援，換個方式試試吧",
		50000: "服務剛剛走神了，稍後再試吧",
		50001: "這個接口還在準備中，先稍等一下",
		-1:    "服務剛剛走神了，稍後再試吧",
	},
	"ja-JP": {
		10000: "今の状態ではこの操作はできません。少し変えて試してみてください。",
		10001: "このデータはもうあります。重ねて送らなくて大丈夫です。",
		10002: "該当するデータが見つかりませんでした。確認してもう一度お試しください。",
		10003: "この内容は今は削除できません。",
		10004: "この内容は今は変更できません。",
		10005: "リクエスト内容が少し違うみたいです。もう一度確認してください。",
		20001: "このアカウントは登録済みです。ログインを試してみてください。",
		20002: "ログイン状態が切れたみたいです。もう一度ログインしてください。",
		20003: "紐づくユーザーが見つかりませんでした。ログインし直してみてください。",
		20004: "このユーザーが見つかりませんでした。アカウント情報を確認してください。",
		20005: "このアカウントが見つかりませんでした。確認してもう一度お試しください。",
		21001: "この会話が見つからないか、今はアクセスできません。",
		21002: "音声を作れるのは AI メッセージだけです。",
		21003: "メッセージが空です。少し書いてから送ってください。",
		21004: "音声設定がまだ整っていないため、今は生成できません。",
		22001: "このキャラクターカードが見つかりません。別の子も見てみましょう。",
		22002: "キャラクター設定がまだ準備中です。今は使えません。",
		22003: "このキャラクターカードはまだ公開されていません。別の子を見てみましょう。",
		22004: "このキャラクターカードのバージョンが見つかりません。",
		22005: "このワールドカードはまだ公開されていません。別のものを見てみましょう。",
		22006: "AI 設定が見つからないため、今は会話を始められません。",
		22007: "この AI 設定はまだ公開されていません。別の設定を試してみてください。",
		23001: "ファイルが選ばれていません。先にファイルを選んでください。",
		23002: "ファイルの読み込みに失敗しました。もう一度アップロードしてみてください。",
		23003: "このファイルが見つかりませんでした。",
		24001: "音声サービスの準備がまだです。今は音声を生成できません。",
		24002: "画像サービスの準備がまだです。今は画像を生成できません。",
		24003: "モデル設定がまだ準備中です。少し待ってからお試しください。",
		24004: "音声機能はまだ準備中です。少しだけ待ってください。",
		24005: "画像生成機能はまだ準備中です。少しだけ待ってください。",
		24006: "サービス設定がまだ整っていません。少し待ってからお試しください。",
		24007: "音声サービスが少し迷子になりました。後でもう一度お試しください。",
		24008: "AI の生成に失敗しました。少し待ってからもう一度お試しください。",
		24009: "接続が切れました。更新してもう一度お試しください。",
		24010: "記憶を整理しています。今は続けて返信できないので、少し待ってください。",
		25001: "該当する Prompt テンプレートが見つかりませんでした。",
		26001: "ポイントが足りません。チャージしてからもう一度お試しください。",
		26002: "この注文が見つかりませんでした。確認してもう一度お試しください。",
		26003: "この注文はすでに支払い済みです。",
		26004: "この注文は期限切れです。もう一度注文してください。",
		27001: "認証コードが違うみたいです。もう一度確認してください。",
		27002: "認証コードの期限が切れました。もう一度取得してください。",
		27003: "リクエストが多すぎます。少し待ってからお試しください。",
		28001: "認可に失敗しました。もう一度認可してみてください。",
		31001: "今回のチャット状態が見つかりません。もう一度送ってみてください。",
		40004: "該当するデータが見つかりませんでした。確認してもう一度お試しください。",
		403:   "ここには今アクセスできません。",
		404:   "この API が見つかりません。パスを確認してください。",
		405:   "このリクエスト方法には対応していません。別の方法で試してください。",
		50000: "サービスが少し迷子になりました。後でもう一度お試しください。",
		50001: "この API はまだ準備中です。少し待ってください。",
		-1:    "サービスが少し迷子になりました。後でもう一度お試しください。",
	},
	"ko-KR": {
		10000: "지금 상태에서는 이 작업을 할 수 없어요. 다른 방법으로 시도해 주세요.",
		10001: "이미 있는 데이터예요. 다시 보내지 않아도 괜찮아요.",
		10002: "해당 데이터를 찾지 못했어요. 확인하고 다시 시도해 주세요.",
		10003: "이 내용은 지금 삭제할 수 없어요.",
		10004: "이 내용은 지금 수정할 수 없어요.",
		10005: "요청 값이 조금 맞지 않아요. 다시 확인해 주세요.",
		20001: "이미 가입된 계정이에요. 바로 로그인해 보세요.",
		20002: "로그인 상태가 만료된 것 같아요. 다시 로그인해 주세요.",
		20003: "연결된 사용자를 찾지 못했어요. 다시 로그인해 보세요.",
		20004: "이 사용자를 찾지 못했어요. 계정 정보를 확인해 주세요.",
		20005: "이 계정을 찾지 못했어요. 확인하고 다시 시도해 주세요.",
		21001: "이 대화를 찾지 못했거나 지금은 접근할 수 없어요.",
		21002: "음성은 AI 메시지에서만 만들 수 있어요.",
		21003: "메시지가 비어 있어요. 조금 작성한 뒤 보내 주세요.",
		21004: "음성 설정이 아직 준비되지 않아 지금은 생성할 수 없어요.",
		22001: "이 캐릭터 카드를 찾지 못했어요. 다른 캐릭터도 둘러보세요.",
		22002: "캐릭터 설정이 아직 준비되지 않아 지금은 사용할 수 없어요.",
		22003: "이 캐릭터 카드는 아직 공개되지 않았어요. 다른 캐릭터를 먼저 봐 주세요.",
		22004: "이 캐릭터 카드 버전을 찾지 못했어요.",
		22005: "이 월드 카드는 아직 공개되지 않았어요. 다른 월드를 먼저 봐 주세요.",
		22006: "AI 설정을 찾지 못해 지금은 대화를 시작할 수 없어요.",
		22007: "이 AI 설정은 아직 공개되지 않았어요. 다른 설정을 시도해 주세요.",
		23001: "파일을 아직 선택하지 않았어요. 먼저 파일을 골라 주세요.",
		23002: "파일을 읽지 못했어요. 다시 업로드해 보세요.",
		23003: "이 파일 리소스를 찾지 못했어요.",
		24001: "음성 서비스가 아직 준비되지 않아 지금은 생성할 수 없어요.",
		24002: "이미지 서비스가 아직 준비되지 않아 지금은 생성할 수 없어요.",
		24003: "모델 설정이 아직 준비 중이에요. 잠시 후 다시 시도해 주세요.",
		24004: "음성 기능이 아직 준비 중이에요. 잠시만 기다려 주세요.",
		24005: "이미지 생성 기능이 아직 준비 중이에요. 잠시만 기다려 주세요.",
		24006: "서비스 설정이 아직 준비되지 않았어요. 잠시 후 다시 시도해 주세요.",
		24007: "음성 서비스가 잠깐 길을 잃었어요. 잠시 후 다시 시도해 주세요.",
		24008: "AI 생성에 실패했어요. 잠시 후 다시 시도해 주세요.",
		24009: "연결이 끊어졌어요. 새로고침 후 다시 시도해 주세요.",
		24010: "기억을 정리하는 중이에요. 잠시 후 다시 시도해 주세요.",
		25001: "해당 Prompt 템플릿을 찾지 못했어요.",
		26001: "포인트가 부족해요. 충전한 뒤 다시 시도해 주세요.",
		26002: "이 주문을 찾지 못했어요. 확인하고 다시 시도해 주세요.",
		26003: "이 주문은 이미 결제되었어요.",
		26004: "이 주문은 만료되었어요. 다시 주문해 주세요.",
		27001: "인증 코드가 맞지 않아요. 다시 확인해 주세요.",
		27002: "인증 코드가 만료되었어요. 다시 받아 주세요.",
		27003: "요청이 너무 잦아요. 잠시 후 다시 시도해 주세요.",
		28001: "인증에 실패했어요. 다시 승인해 보세요.",
		31001: "이번 채팅 상태를 찾지 못했어요. 다시 보내 보세요.",
		40004: "해당 데이터를 찾지 못했어요. 확인하고 다시 시도해 주세요.",
		403:   "지금은 여기에 접근할 수 없어요.",
		404:   "이 API를 찾지 못했어요. 경로를 확인해 주세요.",
		405:   "이 요청 방식은 지원하지 않아요. 다른 방식으로 시도해 주세요.",
		50000: "서비스가 잠깐 길을 잃었어요. 잠시 후 다시 시도해 주세요.",
		50001: "이 API는 아직 준비 중이에요. 잠시만 기다려 주세요.",
		-1:    "서비스가 잠깐 길을 잃었어요. 잠시 후 다시 시도해 주세요.",
	},
	"ru-RU": {
		10000: "Сейчас это действие недоступно. Попробуйте немного иначе.",
		10001: "Такие данные уже есть. Повторно отправлять не нужно.",
		10002: "Данные не найдены. Проверьте и попробуйте снова.",
		10003: "Сейчас это нельзя удалить.",
		10004: "Сейчас это нельзя изменить.",
		10005: "В запросе что-то не так. Проверьте поля и попробуйте снова.",
		20001: "Этот аккаунт уже зарегистрирован. Попробуйте войти.",
		20002: "Похоже, вход истек. Пожалуйста, войдите снова.",
		20003: "Связанный пользователь не найден. Попробуйте войти снова.",
		20004: "Пользователь не найден. Проверьте данные аккаунта.",
		20005: "Аккаунт не найден. Проверьте и попробуйте снова.",
		21001: "Диалог не найден или сейчас недоступен.",
		21002: "Голос можно создать только для сообщения AI.",
		21003: "Сообщение пустое. Напишите немного и отправьте снова.",
		21004: "Голосовые настройки еще не готовы, поэтому создать голос пока нельзя.",
		22001: "Карточка персонажа не найдена. Попробуйте выбрать другую.",
		22002: "Профиль персонажа еще не готов, пока использовать его нельзя.",
		22003: "Карточка персонажа еще не опубликована. Посмотрите другую.",
		22004: "Версия карточки персонажа не найдена.",
		22005: "Карточка мира еще не опубликована. Посмотрите другую.",
		22006: "AI-настройка не найдена, поэтому начать диалог сейчас нельзя.",
		22007: "Эта AI-настройка еще не опубликована. Попробуйте другую.",
		23001: "Файл не выбран. Сначала выберите файл для загрузки.",
		23002: "Не удалось прочитать файл. Попробуйте загрузить его еще раз.",
		23003: "Файл не найден.",
		24001: "Голосовой сервис еще не готов, голос пока создать нельзя.",
		24002: "Сервис изображений еще не готов, изображение пока создать нельзя.",
		24003: "Настройка модели еще не готова. Попробуйте чуть позже.",
		24004: "Голосовая функция еще готовится. Подождите немного.",
		24005: "Генерация изображений еще готовится. Подождите немного.",
		24006: "Настройка сервиса еще не готова. Попробуйте чуть позже.",
		24007: "Голосовой сервис немного задумался. Попробуйте позже.",
		24008: "AI не смог сгенерировать ответ. Попробуйте чуть позже.",
		24009: "Соединение прервалось. Обновите страницу и попробуйте снова.",
		24010: "Я привожу память в порядок. Сейчас нельзя продолжить ответ, попробуйте чуть позже.",
		25001: "Нужный Prompt-шаблон не найден.",
		26001: "Недостаточно баллов. Пополните баланс и попробуйте снова.",
		26002: "Заказ не найден. Проверьте и попробуйте снова.",
		26003: "Этот заказ уже оплачен.",
		26004: "Срок действия заказа истек. Создайте новый заказ.",
		27001: "Код подтверждения неверный. Проверьте его еще раз.",
		27002: "Код подтверждения истек. Получите новый.",
		27003: "Слишком много запросов. Попробуйте чуть позже.",
		28001: "Авторизация не удалась. Попробуйте разрешить доступ снова.",
		31001: "Состояние этого чата не найдено. Попробуйте отправить еще раз.",
		40004: "Данные не найдены. Проверьте и попробуйте снова.",
		403:   "Сейчас у вас нет доступа сюда.",
		404:   "Этот API не найден. Проверьте путь.",
		405:   "Этот метод запроса не поддерживается. Попробуйте другой.",
		50000: "Сервис немного задумался. Попробуйте позже.",
		50001: "Этот API еще готовится. Подождите немного.",
		-1:    "Сервис немного задумался. Попробуйте позже.",
	},
}

var clientFieldNames = map[string]map[string]string{
	"en-US": {
		"userId": "user ID",
		"update": "update payload",
	},
	"zh-CN": {
		"userId": "用户 ID",
		"update": "更新内容",
	},
	"zh-TW": {
		"userId": "使用者 ID",
		"update": "更新內容",
	},
	"ja-JP": {
		"userId": "ユーザー ID",
		"update": "更新内容",
	},
	"ko-KR": {
		"userId": "사용자 ID",
		"update": "수정 내용",
	},
	"ru-RU": {
		"userId": "ID пользователя",
		"update": "данные для обновления",
	},
}

// Core status details 对应的用户端错误响应。
func FromCoreError(err error, language string) *Error {
	st, ok := status.FromError(err)
	if !ok {
		return Internal()
	}
	language = normalizeLanguage(language)
	code, hasCode := businessCode(st)
	if !hasCode {
		return fromGRPCCode(st.Code(), language)
	}
	message := messageForCode(code, language)
	if validationMessage, ok := firstValidationMessage(st, language); ok {
		message = validationMessage
	}
	return &Error{Code: code, Message: message, HTTPStatus: httpStatusForGRPCCode(st.Code())}
}

// Core ErrorInfo details 中的业务码。
func businessCode(st *status.Status) (int, bool) {
	for _, detail := range st.Details() {
		info, ok := detail.(*errdetails.ErrorInfo)
		if !ok || info.Reason != coreErrorReason {
			continue
		}
		code, err := strconv.Atoi(info.Metadata[coreCodeKey])
		return code, err == nil
	}
	return 0, false
}

// Core 字段校验明细对应的用户端提示。
func firstValidationMessage(st *status.Status, language string) (string, bool) {
	for _, detail := range st.Details() {
		badRequest, ok := detail.(*errdetails.BadRequest)
		if !ok || len(badRequest.FieldViolations) == 0 {
			continue
		}
		item := badRequest.FieldViolations[0]
		fieldName := clientFieldName(item.GetField(), language)
		return validationMessage(fieldName, item.GetDescription(), language), true
	}
	return "", false
}

// 对外 JSON 字段名对应的用户端可读名称。
func clientFieldName(field string, language string) string {
	if names, ok := clientFieldNames[language]; ok {
		if name, ok := names[field]; ok {
			return name
		}
	}
	if name, ok := clientFieldNames["en-US"][field]; ok {
		return name
	}
	return field
}

// 用户端标准文案，未知业务码按服务错误处理。
func messageForCode(code int, language string) string {
	if messages, ok := clientMessages[language]; ok {
		if message, ok := messages[code]; ok {
			return message
		}
	}
	if message, ok := clientMessages["en-US"][code]; ok {
		return message
	}
	return InternalError.Message
}

// 指定语言下的用户端错误文案。
func MessageForCode(code int, language string) string {
	return messageForCode(code, normalizeLanguage(language))
}

// Core 未携带业务码时，按传输层分类兜底。
func fromGRPCCode(code codes.Code, language string) *Error {
	switch code {
	case codes.InvalidArgument:
		return WithMessage(ParameterValidationError, messageForCode(ParameterValidationError.Code, language))
	case codes.NotFound:
		return WithMessage(DataNotFound, messageForCode(DataNotFound.Code, language))
	case codes.Unauthenticated:
		return WithMessage(InvalidCredentials, messageForCode(InvalidCredentials.Code, language))
	case codes.PermissionDenied:
		return WithMessage(AccessDenied, messageForCode(AccessDenied.Code, language))
	case codes.FailedPrecondition:
		return WithMessage(commonBizResult(), messageForCode(10000, language))
	default:
		return Internal()
	}
}

// gRPC 状态对应的用户端 HTTP 状态。
func httpStatusForGRPCCode(code codes.Code) int {
	switch code {
	case codes.InvalidArgument, codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

func commonBizResult() ResultCode {
	return ResultCode{Code: 10000, Message: "Business state does not allow this action", HTTPStatus: http.StatusBadRequest}
}

func normalizeLanguage(language string) string {
	switch language {
	case "en-US", "zh-CN", "zh-TW", "ja-JP", "ko-KR", "ru-RU":
		return language
	default:
		return "en-US"
	}
}

func validationMessage(fieldName string, reason string, language string) string {
	switch language {
	case "zh-CN":
		switch reason {
		case "required":
			return fieldName + " 不能为空"
		case "invalid_format":
			return fieldName + " 格式不正确"
		default:
			return fieldName + " 不正确"
		}
	case "zh-TW":
		switch reason {
		case "required":
			return fieldName + " 不能為空"
		case "invalid_format":
			return fieldName + " 格式不正確"
		default:
			return fieldName + " 不正確"
		}
	case "ja-JP":
		switch reason {
		case "required":
			return fieldName + " は必須です"
		case "invalid_format":
			return fieldName + " の形式が正しくありません"
		default:
			return fieldName + " が正しくありません"
		}
	case "ko-KR":
		switch reason {
		case "required":
			return fieldName + "은(는) 필수예요"
		case "invalid_format":
			return fieldName + " 형식이 올바르지 않아요"
		default:
			return fieldName + "이(가) 올바르지 않아요"
		}
	case "ru-RU":
		switch reason {
		case "required":
			return fieldName + " обязателен"
		case "invalid_format":
			return fieldName + " имеет неверный формат"
		default:
			return fieldName + " указан неверно"
		}
	}
	switch reason {
	case "required":
		return fieldName + " is required"
	case "invalid_format":
		return fieldName + " has an invalid format"
	default:
		return fieldName + " is invalid"
	}
}
